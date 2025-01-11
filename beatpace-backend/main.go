package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"encoding/base64"
	"io/ioutil"

	"github.com/rs/cors"
)

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// getSpotifyAccessToken exchanges an authorization code for an access token
func getSpotifyAccessToken(authCode string) (*SpotifyTokenResponse, error) {
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

func generatePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	authorizationCode, _ := requestData["authorization_code"].(string)
	/* paceUnit, _ := requestData["paceUnit"].(string)
	paceInSeconds, _ := requestData["paceInSeconds"].(float64)
	gender, _ := requestData["gender"].(string)
	height, _ := requestData["height"].(float64)
	heightUnit, _ := requestData["heightUnit"].(string) */

	// Exchange auth code for access token
	tokenResponse, err := getSpotifyAccessToken(authorizationCode)
	if err != nil {
		log.Println("Error getting access token:", err)
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Log and return the access token
	log.Println("Access token:", tokenResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenResponse.AccessToken,
		"message":      "Access token successfully retrieved",
	})
}

func main() {
	// Create a new CORS handler allowing requests from localhost:3000
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Change this to your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Set up your routes
	http.HandleFunc("/api/generate-playlist", generatePlaylistHandler)

	// Wrap your handlers with the CORS handler
	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	fmt.Println("Server is running on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", handler))
}
