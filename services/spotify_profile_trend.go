package services

import (
	"fmt"
	"math/rand"
	"time"
)

// TrendServiceParams contains parameters for fetching user trends
type TrendServiceParams struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

// TrendTrackData represents a track in the trends response
type TrendTrackData struct {
	Rank   int            `json:"rank"`
	Title  string         `json:"title"`
	Artist string         `json:"artist"`
	Album  TrendAlbumData `json:"album"`
}

// TrendAlbumData represents album information in the trend
type TrendAlbumData struct {
	Name     string `json:"name"`
	CoverURL string `json:"coverUrl"`
}

// BehindLyricsData represents the "behind the lyrics" section
type BehindLyricsData struct {
	SummaryEn string          `json:"summaryEn"`
	SummaryZh string          `json:"summaryZh"`
	Title     string          `json:"title"`
	Artist    string          `json:"artist"`
	Album     BehindAlbumData `json:"album"`
}

// BehindAlbumData represents album information for behind the lyrics
type BehindAlbumData struct {
	Name     string `json:"name"`
	CoverURL string `json:"coverUrl"`
}

// TrendsData represents the complete trends data
type TrendsData struct {
	Tracks          []TrendTrackData `json:"tracks"`
	BehindTheLyrics BehindLyricsData `json:"behindTheLyrics"`
}

// ValidateTrendRequestParams validates that required parameters are present
func ValidateTrendRequestParams(params TrendServiceParams) error {
	if params.ClientID == "" {
		return fmt.Errorf("clientID is required")
	}
	if params.ClientSecret == "" {
		return fmt.Errorf("clientSecret is required")
	}
	if params.RefreshToken == "" {
		return fmt.Errorf("refreshToken is required")
	}
	return nil
}

// FetchAndRefreshToken refreshes the Spotify access token
func FetchAndRefreshToken(params TrendServiceParams) (*TokenResponse, error) {
	return RefreshAccessToken(RefreshTokenParams{
		RefreshToken: params.RefreshToken,
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
	})
}

// TransformTracksToTrendData converts Spotify tracks to trend track data
func TransformTracksToTrendData(tracks []SpotifyTrack) []TrendTrackData {
	trendTracks := make([]TrendTrackData, 0, len(tracks))

	for i, track := range tracks {
		// Get the first artist name (tracks can have multiple artists)
		artistName := ""
		if len(track.Artists) > 0 {
			artistName = track.Artists[0].Name
		}

		// Get album cover URL (prefer the largest image)
		coverURL := ""
		if len(track.Album.Images) > 0 {
			coverURL = track.Album.Images[0].URL
		}

		trendTracks = append(trendTracks, TrendTrackData{
			Rank:   i + 1, // Rank starts from 1
			Title:  track.Name,
			Artist: artistName,
			Album: TrendAlbumData{
				Name:     track.Album.Name,
				CoverURL: coverURL,
			},
		})
	}

	return trendTracks
}

// SelectRandomTrack randomly selects a track from the provided tracks
func SelectRandomTrack(tracks []TrendTrackData) *TrendTrackData {
	if len(tracks) == 0 {
		return nil
	}

	// Initialize random seed if not already done
	rand.Seed(time.Now().UnixNano())

	// Select random index
	randomIndex := rand.Intn(len(tracks))
	return &tracks[randomIndex]
}

// BuildBehindLyricsData builds the behind lyrics data with LLM summaries
func BuildBehindLyricsData(track *TrendTrackData, agentClient *AgentClient) (BehindLyricsData, error) {
	var behindLyrics BehindLyricsData

	if track == nil {
		return behindLyrics, fmt.Errorf("no track provided")
	}

	// Call agent to get summaries
	agentResp, err := agentClient.GetBehindTheLyrics(
		track.Title,
		track.Artist,
		track.Album.Name,
		track.Album.CoverURL,
	)
	if err != nil {
		return behindLyrics, fmt.Errorf("failed to get summaries from agent: %w", err)
	}

	// Build the behind lyrics data
	behindLyrics = BehindLyricsData{
		SummaryEn: agentResp.SummaryEn,
		SummaryZh: agentResp.SummaryZh,
		Title:     track.Title,
		Artist:    track.Artist,
		Album: BehindAlbumData{
			Name:     track.Album.Name,
			CoverURL: track.Album.CoverURL,
		},
	}

	return behindLyrics, nil
}

// GetUserTrends orchestrates the entire process of fetching user trends
func GetUserTrends(params TrendServiceParams) (*TrendsData, error) {
	// Validate parameters
	if err := ValidateTrendRequestParams(params); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Refresh access token
	tokenResp, err := FetchAndRefreshToken(params)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	// Fetch user's top tracks (20 tracks, medium term)
	topTracksResp, err := GetUserTopTracks(tokenResp.AccessToken, 20, "medium_term")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top tracks: %w", err)
	}

	// Transform tracks to trend data
	trendTracks := TransformTracksToTrendData(topTracksResp.Items)

	// Select a random track for behind the lyrics
	randomTrack := SelectRandomTrack(trendTracks)

	// Initialize agent client
	agentClient := NewAgentClient()

	// Build behind lyrics data with agent summaries
	behindLyrics, err := BuildBehindLyricsData(randomTrack, agentClient)
	if err != nil {
		// Log the error but don't fail the entire request
		// Return empty summaries instead
		if randomTrack != nil {
			behindLyrics = BehindLyricsData{
				SummaryEn: "",
				SummaryZh: "",
				Title:     randomTrack.Title,
				Artist:    randomTrack.Artist,
				Album: BehindAlbumData{
					Name:     randomTrack.Album.Name,
					CoverURL: randomTrack.Album.CoverURL,
				},
			}
		}
	}

	// Build and return the final trends data
	trendsData := &TrendsData{
		Tracks:          trendTracks,
		BehindTheLyrics: behindLyrics,
	}

	return trendsData, nil
}
