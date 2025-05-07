// fetch userID from Spotify
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchUserService is a service to fetch user data from Spotify
type FetchUserService struct {
	AccessToken string
}
type SpotifyUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

// FetchUser fetches the user data from Spotify
func (s *FetchUserService) FetchUser() (*SpotifyUser, error) {
	url := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user data: %s", resp.Status)
	}

	var user SpotifyUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &user, nil
}
func (s *FetchUserService) FetchUserID() (string, error) {
	user, err := s.FetchUser()
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %v", err)
	}
	return user.ID, nil
}

