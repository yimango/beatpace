package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// getSpotifyAccessToken exchanges an authorization code for an access token
func GetSpotifyAccessToken(authCode string) (*SpotifyTokenResponse, error) {
	// Spotify token endpoint
	url := "https://accounts.spotify.com/api/token"

	// Prepare the data
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s",
		authCode, os.Getenv("REDIRECT_URI"))
	body := bytes.NewBufferString(data)

	// Create a new request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add Basic Auth using Client ID and Client Secret
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read and decode the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d, body: %s", resp.StatusCode, bodyBytes)
	}

	// Parse the JSON response
	var tokenResponse SpotifyTokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &tokenResponse, nil
}