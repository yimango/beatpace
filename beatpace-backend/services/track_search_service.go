package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SearchTracks searches Spotify for tracks based on given criteria
func SearchTracks(accessToken, mood, genre string, bpm int) ([]string, error) {
	searchURL := "https://api.spotify.com/v1/search"

	// Construct the query based on mood, genre, and BPM
	query := fmt.Sprintf("genre:%s mood:%s bpm:%d", genre, mood, bpm)
	url := fmt.Sprintf("%s?q=%s&type=track&limit=10", searchURL, query)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("failed to search tracks: %s", resp.Status)
	}

	var searchResponse struct {
		Tracks struct {
			Items []struct {
				ExternalURL struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	var trackLinks []string
	for _, track := range searchResponse.Tracks.Items {
		trackLinks = append(trackLinks, track.ExternalURL.Spotify)
	}

	return trackLinks, nil
}
