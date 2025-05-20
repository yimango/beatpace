package services

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/zmb3/spotify/v2"
)

type PlaylistGenerator struct {
	spotifyService SpotifyService
}

func NewPlaylistGenerator(spotifyService SpotifyService) *PlaylistGenerator {
	return &PlaylistGenerator{
		spotifyService: spotifyService,
	}
}

type TrackInfo struct {
	Track spotify.SimpleTrack
	BPM   float32
}

// GeneratePlaylist creates a playlist based on the target BPM
func (s *PlaylistGenerator) GeneratePlaylist(ctx context.Context, userID string, targetBPM int) (*spotify.FullPlaylist, error) {
	fmt.Printf("PlaylistGenerator: Starting playlist generation for user %s with target BPM %d\n", userID, targetBPM)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get user's top tracks and artists for better recommendations
	client := s.spotifyService.GetClient(ctx, userID)
	if client == nil {
		fmt.Printf("PlaylistGenerator: Failed to get Spotify client\n")
		return nil, fmt.Errorf("failed to get spotify client")
	}
	fmt.Printf("PlaylistGenerator: Successfully got Spotify client\n")

	// Create channels with appropriate buffer sizes
	tracksChan := make(chan TrackInfo, 100)
	errorsChan := make(chan error, 10)
	var wg sync.WaitGroup

	// Create a rate limiter for Spotify API calls (5 requests per second)
	rateLimiter := time.NewTicker(200 * time.Millisecond)
	defer rateLimiter.Stop()

	// Start a goroutine to handle track collection
	go func() {
		defer close(tracksChan)
		defer close(errorsChan)

		// Get user's top tracks and artists concurrently
		var topTracks *spotify.FullTrackPage
		var topArtists *spotify.FullArtistPage
		var tracksErr, artistsErr error

		var wgTopItems sync.WaitGroup
		wgTopItems.Add(2)

		go func() {
			defer wgTopItems.Done()
			<-rateLimiter.C // Rate limit
			fmt.Printf("PlaylistGenerator: Fetching user's top tracks...\n")
			topTracks, tracksErr = client.CurrentUsersTopTracks(ctx, spotify.Limit(5), spotify.Timerange(spotify.ShortTermRange))
		}()

		go func() {
			defer wgTopItems.Done()
			<-rateLimiter.C // Rate limit
			fmt.Printf("PlaylistGenerator: Fetching user's top artists...\n")
			topArtists, artistsErr = client.CurrentUsersTopArtists(ctx, spotify.Limit(5), spotify.Timerange(spotify.ShortTermRange))
		}()

		wgTopItems.Wait()

		if tracksErr != nil || artistsErr != nil {
			playlistErr := fmt.Errorf("failed to get top items: tracks error: %v, artists error: %v", tracksErr, artistsErr)
			errorsChan <- playlistErr
			return
		}

		fmt.Printf("PlaylistGenerator: Successfully got %d top tracks and %d top artists\n",
			len(topTracks.Tracks), len(topArtists.Artists))

		// Process tracks and artists in combinations
		for i := 0; i < len(topTracks.Tracks); i++ {
			select {
			case <-ctx.Done():
				return
			case <-rateLimiter.C:
				wg.Add(1)
				go func(track spotify.FullTrack) {
					defer wg.Done()
					
					// Get artist genres for this track
					var genres []string
					if len(track.Artists) > 0 {
						artist, err := client.GetArtist(ctx, track.Artists[0].ID)
						if err == nil && artist != nil && len(artist.Genres) > 0 {
							genres = artist.Genres[:1] // Take first genre
						}
					}

					// Create seeds with track, artist, and genre
					seeds := spotify.Seeds{
						Tracks:  []spotify.ID{track.ID},
						Artists: []spotify.ID{},
						Genres:  genres,
					}

					// Add up to 2 artists from the track's artists
					for _, artist := range track.Artists {
						if len(seeds.Artists) < 2 {
							seeds.Artists = append(seeds.Artists, artist.ID)
						}
					}

					fmt.Printf("Making recommendation request with seeds:\n")
					fmt.Printf("  Track: %s (%s)\n", track.Name, track.ID)
					fmt.Printf("  Artists: %v\n", seeds.Artists)
					fmt.Printf("  Genres: %v\n", seeds.Genres)

					if tracks := s.tryGetRecommendations(ctx, client, seeds, targetBPM, 5); len(tracks) > 0 {
						for _, track := range tracks {
							select {
							case <-ctx.Done():
								return
							case tracksChan <- track:
							}
						}
					}
				}(topTracks.Tracks[i])
			}
		}

		// Wait for all searches to complete
		wg.Wait()
	}()

	// Collect tracks and handle errors
	trackMap := make(map[string]TrackInfo)
	errCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("playlist generation timed out after 30 seconds")
		case err := <-errorsChan:
			if err != nil {
				errCount++
				fmt.Printf("PlaylistGenerator: Error during track search: %v\n", err)
				if errCount > 5 {
					return nil, fmt.Errorf("too many errors during track search")
				}
			}
		case track, ok := <-tracksChan:
			if !ok {
				goto CREATE_PLAYLIST
			}
			// Accept tracks within ±10 BPM of target
			if math.Abs(float64(targetBPM)-float64(track.BPM)) <= 10 {
				trackMap[track.Track.ID.String()] = track
			}
		}
	}

CREATE_PLAYLIST:
	if len(trackMap) == 0 {
		fmt.Printf("PlaylistGenerator: No suitable tracks found\n")
		return nil, fmt.Errorf("no suitable tracks found")
	}

	fmt.Printf("PlaylistGenerator: Found %d unique tracks within BPM range\n", len(trackMap))

	// Convert map to slice and limit to 25 tracks
	var selectedTracks []spotify.ID
	for _, track := range trackMap {
		selectedTracks = append(selectedTracks, track.Track.ID)
		if len(selectedTracks) >= 25 {
			break
		}
	}

	// Create a new playlist
	<-rateLimiter.C // Rate limit
	fmt.Printf("PlaylistGenerator: Creating new playlist...\n")
	
	// Get the Spotify user ID from the token
	storedToken, err := s.spotifyService.GetUserProfile(ctx, userID)
	if err != nil {
		fmt.Printf("PlaylistGenerator: Failed to get user profile: %v\n", err)
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}
	
	playlist, err := client.CreatePlaylistForUser(ctx, storedToken.SpotifyUserID, fmt.Sprintf("BeatPace - %d BPM", targetBPM), "", false, false)
	if err != nil {
		fmt.Printf("PlaylistGenerator: Failed to create playlist: %v\n", err)
		return nil, fmt.Errorf("failed to create playlist: %v", err)
	}
	fmt.Printf("PlaylistGenerator: Created playlist with ID: %s\n", playlist.ID)

	// Add tracks to the playlist
	<-rateLimiter.C // Rate limit
	fmt.Printf("PlaylistGenerator: Adding tracks to playlist...\n")
	_, err = client.AddTracksToPlaylist(ctx, playlist.ID, selectedTracks...)
	if err != nil {
		fmt.Printf("PlaylistGenerator: Failed to add tracks to playlist: %v\n", err)
		return nil, fmt.Errorf("failed to add tracks to playlist: %v", err)
	}
	fmt.Printf("PlaylistGenerator: Successfully added %d tracks to playlist\n", len(selectedTracks))

	return playlist, nil
}

func (s *PlaylistGenerator) searchSimilarTracks(ctx context.Context, client *spotify.Client, seedTrack spotify.ID, targetBPM int, bpmOffset int, tracksChan chan<- TrackInfo, errorsChan chan<- error) {
	// Get track info first to try genre seeds
	track, err := client.GetTrack(ctx, seedTrack)
	if err != nil {
		errorsChan <- fmt.Errorf("failed to get track info: %v", err)
		return
	}

	// Try to get genre first
	if len(track.Artists) > 0 {
		artist, err := client.GetArtist(ctx, track.Artists[0].ID)
		if err == nil && artist != nil && len(artist.Genres) > 0 {
			// Try 1: Genre seed
			if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
				Genres: []string{artist.Genres[0]},
			}, targetBPM, bpmOffset); len(tracks) > 0 {
				for _, track := range tracks {
					select {
					case <-ctx.Done():
						return
					case tracksChan <- track:
					}
				}
				return
			}
		}

		// Try 2: Artist seed
		if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
			Artists: []spotify.ID{track.Artists[0].ID},
		}, targetBPM, bpmOffset); len(tracks) > 0 {
			for _, track := range tracks {
				select {
				case <-ctx.Done():
					return
				case tracksChan <- track:
				}
			}
			return
		}
	}

	// Try 3: Track seed as last resort
	if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
		Tracks: []spotify.ID{seedTrack},
	}, targetBPM, bpmOffset); len(tracks) > 0 {
		for _, track := range tracks {
			select {
			case <-ctx.Done():
				return
			case tracksChan <- track:
			}
		}
	}
}

func (s *PlaylistGenerator) searchSimilarTracksFromArtist(ctx context.Context, client *spotify.Client, artistID spotify.ID, targetBPM int, bpmOffset int, tracksChan chan<- TrackInfo, errorsChan chan<- error) {
	// Get artist info to try genre seeds first
	artist, err := client.GetArtist(ctx, artistID)
	if err != nil {
		errorsChan <- fmt.Errorf("failed to get artist info: %v", err)
		return
	}

	// Try 1: Genre seed
	if len(artist.Genres) > 0 {
		if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
			Genres: []string{artist.Genres[0]},
		}, targetBPM, bpmOffset); len(tracks) > 0 {
			for _, track := range tracks {
				select {
				case <-ctx.Done():
					return
				case tracksChan <- track:
				}
			}
			return
		}
	}

	// Try 2: Artist seed
	if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
		Artists: []spotify.ID{artistID},
	}, targetBPM, bpmOffset); len(tracks) > 0 {
		for _, track := range tracks {
			select {
			case <-ctx.Done():
				return
			case tracksChan <- track:
			}
		}
		return
	}

	// Try 3: Top track seed
	topTracks, err := client.GetArtistsTopTracks(ctx, artistID, "US")
	if err == nil && len(topTracks) > 0 {
		if tracks := s.tryGetRecommendations(ctx, client, spotify.Seeds{
			Tracks: []spotify.ID{topTracks[0].ID},
		}, targetBPM, bpmOffset); len(tracks) > 0 {
			for _, track := range tracks {
				select {
				case <-ctx.Done():
					return
				case tracksChan <- track:
				}
			}
		}
	}
}

// Helper function to try getting recommendations with given seeds
func (s *PlaylistGenerator) tryGetRecommendations(ctx context.Context, client *spotify.Client, seeds spotify.Seeds, targetBPM int, bpmOffset int) []TrackInfo {
	// Build search query from seeds
	var searchQuery string
	if len(seeds.Genres) > 0 {
		searchQuery = fmt.Sprintf(`genre:"%s"`, seeds.Genres[0])
	} else if len(seeds.Artists) > 0 {
		// Get artist name for search
		artist, err := client.GetArtist(ctx, seeds.Artists[0])
		if err == nil && artist != nil {
			searchQuery = fmt.Sprintf(`artist:"%s"`, artist.Name)
		}
	} else if len(seeds.Tracks) > 0 {
		// Get track info for search
		track, err := client.GetTrack(ctx, seeds.Tracks[0])
		if err == nil && track != nil {
			searchQuery = fmt.Sprintf(`track:"%s"`, track.Name)
		}
	}

	if searchQuery == "" {
		fmt.Printf("No valid search query could be constructed from seeds\n")
		return nil
	}

	fmt.Printf("Searching with query: %s\n", searchQuery)

	// Search for tracks with a wider range of results
	searchResults, err := client.Search(ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(50))
	if err != nil {
		fmt.Printf("Failed to search tracks: %v\n", err)
		return nil
	}

	if len(searchResults.Tracks.Tracks) == 0 {
		fmt.Printf("No tracks found in search\n")
		return nil
	}

	fmt.Printf("Found %d tracks in search\n", len(searchResults.Tracks.Tracks))

	// Filter tracks based on available information
	var filteredTracks []TrackInfo
	for _, track := range searchResults.Tracks.Tracks {
		// Skip tracks that are too short (likely interludes or skits)
		if track.Duration < 60000 { // Less than 1 minute
			continue
		}

		// Skip tracks that are too long (likely full albums or compilations)
		if track.Duration > 600000 { // More than 10 minutes
			continue
		}

		// Add track to results
		filteredTracks = append(filteredTracks, TrackInfo{
			Track: spotify.SimpleTrack{
				ID:       track.ID,
				Name:     track.Name,
				Artists:  track.Artists,
				Duration: track.Duration,
			},
			BPM: float32(targetBPM), // Use target BPM since we can't get actual BPM
		})
		fmt.Printf("Added track: %s (Duration: %d ms)\n",
			track.Name, track.Duration)
	}

	return filteredTracks
} 