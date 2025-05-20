// services/spotify_service.go
package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/yimango/beatpace-backend/model"
	"github.com/yimango/beatpace-backend/repository"
)

// PlaylistResponse is the shape returned by GeneratePlaylistForPace.
type PlaylistResponse struct {
	URL    string   `json:"url"`
	Tracks []string `json:"tracks"`
}

type SpotifyServiceImpl struct {
	userRepo      repository.UserRepository
	tokenRepo     repository.TokenRepository
	clientID      string
	clientSecret  string
	redirectURI   string
	auth          *spotifyauth.Authenticator
}

// SpotifyTokenResponse represents the response from Spotify's token endpoint
type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// SpotifyUserProfile represents the response from Spotify's user profile endpoint
type SpotifyUserProfile struct {
	ID    string  `json:"id"`
	Email *string `json:"email"`
}


// NewSpotifyService stays the same...
func NewSpotifyService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
) *SpotifyServiceImpl {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserTopRead,
		),
	)

	return &SpotifyServiceImpl{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   redirectURI,
		auth:         auth,
	}
}

// GetAuthURL returns the Spotify authorization URL
func (s *SpotifyServiceImpl) GetAuthURL() string {
	scopes := []string{
		"user-read-private",
		"user-read-email",
		"playlist-modify-public",
		"playlist-modify-private",
		"user-top-read",
	}
	
	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", s.redirectURI)
	params.Set("scope", strings.Join(scopes, " "))
	params.Set("show_dialog", "true")

	return "https://accounts.spotify.com/authorize?" + params.Encode()
}

// HandleCallback exchanges code & persists tokens
func (s *SpotifyServiceImpl) HandleCallback(ctx context.Context, code string) (string, error) {
	fmt.Printf("HandleCallback: Starting callback process with code length: %d\n", len(code))
	fmt.Printf("HandleCallback: Using redirect URI: %s\n", s.redirectURI)
	
	// Exchange code for token using the auth library
	fmt.Printf("HandleCallback: Attempting to exchange code for token...\n")
	token, err := s.auth.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("HandleCallback: Failed to exchange code for token: %v\n", err)
		return "", fmt.Errorf("failed to exchange code: %v", err)
	}
	fmt.Printf("HandleCallback: Successfully exchanged code for token\n")

	// Create a client using the token
	fmt.Printf("HandleCallback: Creating Spotify client...\n")
	client := spotify.New(s.auth.Client(ctx, token))

	// Get user profile
	fmt.Printf("HandleCallback: Fetching user profile...\n")
	user, err := client.CurrentUser(ctx)
	if err != nil {
		fmt.Printf("HandleCallback: Failed to get user profile: %v\n", err)
		return "", fmt.Errorf("failed to get user profile: %v", err)
	}
	fmt.Printf("HandleCallback: Successfully got user profile - ID: %s, Email: %s\n", user.ID, user.Email)

	// Create or update user
	fmt.Printf("HandleCallback: Checking if user exists in database...\n")
	dbUser, err := s.userRepo.GetUserBySpotifyID(ctx, user.ID)
	if err != nil {
		fmt.Printf("HandleCallback: User not found in database (error: %v), creating new user...\n", err)
		// User doesn't exist, create new one
		dbUser = &model.User{
			ID:           uuid.New(),
			SpotifyUserID: user.ID,
			Email:        &user.Email,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.userRepo.CreateUser(ctx, dbUser); err != nil {
			fmt.Printf("HandleCallback: Failed to create user in database: %v\n", err)
			return "", fmt.Errorf("failed to create user: %v", err)
		}
		fmt.Printf("HandleCallback: Successfully created new user with ID: %s\n", dbUser.ID)
	} else {
		fmt.Printf("HandleCallback: Found existing user in database with ID: %s\n", dbUser.ID)
		// Update existing user's email if it has changed
		if user.Email != "" && (dbUser.Email == nil || *dbUser.Email != user.Email) {
			fmt.Printf("HandleCallback: Updating user email from %v to %s\n", dbUser.Email, user.Email)
			dbUser.Email = &user.Email
			dbUser.UpdatedAt = time.Now()
			if err := s.userRepo.UpdateUser(ctx, dbUser); err != nil {
				fmt.Printf("HandleCallback: Failed to update user in database: %v\n", err)
				return "", fmt.Errorf("failed to update user: %v", err)
			}
			fmt.Printf("HandleCallback: Successfully updated user email\n")
		}
	}

	// Save Spotify tokens
	fmt.Printf("HandleCallback: Saving Spotify tokens for user %s...\n", dbUser.ID)
	spotifyToken := &model.SpotifyToken{
		InternalUserID: dbUser.ID.String(),
		SpotifyUserID:  user.ID,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		GeneratedAt:    time.Now(),
		ExpiresAt:      token.Expiry,
	}

	if err := s.tokenRepo.SaveSpotifyToken(ctx, spotifyToken); err != nil {
		fmt.Printf("HandleCallback: Failed to save Spotify token: %v\n", err)
		return "", fmt.Errorf("failed to save spotify token: %v", err)
	}
	fmt.Printf("HandleCallback: Successfully saved Spotify tokens\n")

	fmt.Printf("HandleCallback: Successfully completed Spotify authentication for user %s\n", dbUser.ID)
	return dbUser.ID.String(), nil
}

func (s *SpotifyServiceImpl) exchangeCodeForToken(code string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectURI)
	data.Set("scope", "user-read-private user-read-email playlist-modify-public playlist-modify-private user-top-read")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify returned status code %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &tokenResp, nil
}

func (s *SpotifyServiceImpl) getUserProfile(accessToken string) (*SpotifyUserProfile, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify returned status code %d", resp.StatusCode)
	}

	var profile SpotifyUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// GetValidAccessToken loads (and refreshes) the access token for a user
func (s *SpotifyServiceImpl) GetValidAccessToken(ctx context.Context, internalUserID string) (string, error) {
	token, err := s.tokenRepo.GetSpotifyToken(ctx, internalUserID)
	if err != nil {
		return "", fmt.Errorf("failed to get spotify token: %v", err)
	}

	// Check if token needs refresh
	if time.Now().After(token.ExpiresAt) {
		newToken, err := s.RefreshToken(ctx, internalUserID)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %v", err)
		}
		return newToken.AccessToken, nil
	}

	return token.AccessToken, nil
}

// GeneratePlaylistForPace creates a playlist based on pace/gender/height
func (s *SpotifyServiceImpl) GeneratePlaylistForPace(
	ctx context.Context,
	internalUserID string,
	paceSeconds int,
	gender string,
	height int,
) (*PlaylistResponse, error) {
	fmt.Printf("Generating playlist for user %s with pace %d seconds\n", internalUserID, paceSeconds)

	// Calculate target BPM based on pace
	// Running cadence is typically between 150-180 BPM
	// Faster pace = higher BPM
	// For a 5:00 min/km pace (300 seconds), we want ~180 BPM
	// For a 6:00 min/km pace (360 seconds), we want ~150 BPM
	// Linear interpolation: BPM = 270 - (paceSeconds/3)
	targetBPM := int(270 - (float64(paceSeconds) / 3.0))
	
	// Clamp BPM to reasonable range
	if targetBPM < 140 {
		targetBPM = 140
	}
	if targetBPM > 180 {
		targetBPM = 180
	}
	
	fmt.Printf("Calculated target BPM: %d\n", targetBPM)

	// Create playlist generator
	generator := NewPlaylistGenerator(s)

	// Generate playlist
	playlist, err := generator.GeneratePlaylist(ctx, internalUserID, targetBPM)
	if err != nil {
		fmt.Printf("Failed to generate playlist: %v\n", err)
		return nil, fmt.Errorf("failed to generate playlist: %v", err)
	}
	fmt.Printf("Generated playlist with ID: %s\n", playlist.ID)

	// Get the tracks in the playlist
	client := s.GetClient(ctx, internalUserID)
	if client == nil {
		fmt.Printf("Failed to get Spotify client\n")
		return nil, fmt.Errorf("failed to get spotify client")
	}
	fmt.Printf("Got Spotify client\n")

	tracks, err := client.GetPlaylistTracks(ctx, playlist.ID)
	if err != nil {
		fmt.Printf("Failed to get playlist tracks: %v\n", err)
		return nil, fmt.Errorf("failed to get playlist tracks: %v", err)
	}
	fmt.Printf("Got %d tracks from playlist\n", len(tracks.Tracks))

	// Extract track URLs
	var trackURLs []string
	for _, item := range tracks.Tracks {
		trackURLs = append(trackURLs, fmt.Sprintf("https://open.spotify.com/track/%s", item.Track.ID))
	}
	fmt.Printf("Extracted %d track URLs\n", len(trackURLs))

	response := &PlaylistResponse{
		URL:    fmt.Sprintf("https://open.spotify.com/playlist/%s", playlist.ID),
		Tracks: trackURLs,
	}
	fmt.Printf("Returning response with URL: %s and %d tracks\n", response.URL, len(response.Tracks))

	return response, nil
}

func (s *SpotifyServiceImpl) GetUserProfile(ctx context.Context, userID string) (*model.SpotifyToken, error) {
	token, err := s.tokenRepo.GetSpotifyToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spotify token: %v", err)
	}

	// Check if token needs refresh
	if time.Now().After(token.ExpiresAt) {
		return s.RefreshToken(ctx, userID)
	}

	return token, nil
}

// GetClient returns a Spotify client for the given user
func (s *SpotifyServiceImpl) GetClient(ctx context.Context, userID string) *spotify.Client {
	fmt.Printf("GetClient: Getting client for user %s\n", userID)

	// Get the stored token from the database
	storedToken, err := s.tokenRepo.GetSpotifyToken(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to get stored token: %v\n", err)
		return nil
	}
	fmt.Printf("GetClient: Got stored token for Spotify user %s\n", storedToken.SpotifyUserID)

	// Create token
	token := &oauth2.Token{
		AccessToken:  storedToken.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: storedToken.RefreshToken,
		Expiry:      storedToken.ExpiresAt,
	}

	// Create HTTP client with token source
	httpClient := s.auth.Client(ctx, token)

	// Create Spotify client
	client := spotify.New(httpClient)

	// Test the client by getting the current user's profile
	user, err := client.CurrentUser(ctx)
	if err != nil {
		fmt.Printf("GetClient: Failed to get current user profile: %v\n", err)
		return nil
	}
	fmt.Printf("GetClient: Successfully verified client for Spotify user %s\n", user.ID)

	return client
}

// RefreshToken refreshes the Spotify access token for a user
func (s *SpotifyServiceImpl) RefreshToken(ctx context.Context, userID string) (*model.SpotifyToken, error) {
	// Get current token
	storedToken, err := s.tokenRepo.GetSpotifyToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spotify token: %v", err)
	}

	// Create token for refresh
	token := &oauth2.Token{
		AccessToken:  storedToken.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: storedToken.RefreshToken,
		Expiry:      storedToken.ExpiresAt,
	}

	// Use the authenticator to refresh the token
	client := s.auth.Client(ctx, token)
	newToken, err := client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}

	// Update token in database
	updatedToken := &model.SpotifyToken{
		InternalUserID: userID,
		SpotifyUserID:  storedToken.SpotifyUserID,
		AccessToken:    newToken.AccessToken,
		RefreshToken:   newToken.RefreshToken,
		GeneratedAt:    time.Now(),
		ExpiresAt:      newToken.Expiry,
	}

	if err := s.tokenRepo.SaveSpotifyToken(ctx, updatedToken); err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %v", err)
	}

	return updatedToken, nil
}