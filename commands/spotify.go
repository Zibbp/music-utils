package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/zibbp/music-utils/config"
	"github.com/zibbp/music-utils/spotify"
	"github.com/zibbp/music-utils/tidal"
	"github.com/zibbp/music-utils/utils"
)

func SaveSpotifyPlaylists(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService) error {
	// Initialize Spotify client
	spotifyClient, err := spotify.Initialize(envConfig.SpotifyClientId, envConfig.SpotifyClientSecret, envConfig.SpotifyRedirectUri, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing spotify client: %v", err)
	}

	log.Info().Msg("fetching spotify playlists")

	// Get all user's Spotify playlists
	spotifyPlaylists, err := spotifyClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	for _, spotifyPlaylist := range spotifyPlaylists {
		log.Info().Str("playlist", spotifyPlaylist.Name).Msg("processing playlist")

		// Get playlist tracks
		spotifyPlaylistTracks, err := spotifyClient.GetPlaylistTracks(spotifyPlaylist.ID)
		if err != nil {
			log.Error().Err(err).Msg("error getting Spotify playlist tracks")
			continue
		}

		// Convert struct to `map[string]any`
		var playlistMap map[string]any
		playlistJson, _ := json.Marshal(spotifyPlaylist) // Convert to JSON
		json.Unmarshal(playlistJson, &playlistMap)       // Convert JSON to map

		// Override "tracks" field
		playlistMap["tracks"] = spotifyPlaylistTracks

		// Write updated playlist to file
		if err := utils.WriteJsonToFile("/data/spotify", string(spotifyPlaylist.ID), playlistMap); err != nil {
			log.Error().Err(err).Msg("error writing playlist to file")
			continue
		}
	}

	return nil
}

func PrintSpotifyPlaylists(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService) error {
	// Initialize Spotify client
	spotifyClient, err := spotify.Initialize(envConfig.SpotifyClientId, envConfig.SpotifyClientSecret, envConfig.SpotifyRedirectUri, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing spotify client: %v", err)
	}

	// Get all user's Spotify playlists
	spotifyPlaylists, err := spotifyClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	data := [][]string{
		{"ID", "Name", "Description"},
	}

	for _, spotifyPlaylist := range spotifyPlaylists {
		data = append(data, []string{spotifyPlaylist.ID.String(), spotifyPlaylist.Name, spotifyPlaylist.Description})
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, row := range data {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()

	return nil
}

func CreateSpotifyPlaylistsOnTidal(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService) error {
	// Initialize Spotify client
	spotifyClient, err := spotify.Initialize(envConfig.SpotifyClientId, envConfig.SpotifyClientSecret, envConfig.SpotifyRedirectUri, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing spotify client: %v", err)
	}

	// Initialize Tidal client
	tidalClient, err := tidal.Initialize(envConfig.TidalClientId, envConfig.TidalClientSecret, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing tidal client: %v", err)
	}

	// Get all user's Spotify playlists
	spotifyPlaylists, err := spotifyClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	// Get all user's Tidal playlists
	tidalPlaylists, err := tidalClient.GetUserPlaylists()
	if err != nil {
		return fmt.Errorf("error getting tidal playlists: %v", err)
	}

	for _, spotifyPlaylist := range spotifyPlaylists {

		// check if spotify playlist id exists in tidal description
		// if it does not exist, create a new playlist
		found := false
		tidalPlaylist := tidal.PlaylistItemV2{}

		checkString := spotifyPlaylist.ID.String()
		playlistDescription := fmt.Sprintf("%s\n%s", spotifyPlaylist.Description, string(spotifyPlaylist.ID))

		for _, tidtidalPlaylist := range tidalPlaylists {
			if strings.Contains(tidtidalPlaylist.Data.Description, checkString) {
				found = true
				tidalPlaylist = tidtidalPlaylist
				break
			}
		}

		if !found {
			// create new playlist
			var playlistName string
			if spotifyPlaylist.Name == "" {
				playlistName = "Untitled"
			} else {
				playlistName = spotifyPlaylist.Name
			}
			log.Info().Msgf("creating tidal playlist: %s", spotifyPlaylist.Name)
			_, err := tidalClient.CreatePlaylist(playlistName, playlistDescription)
			if err != nil {
				return err
			}
			continue
		}

		// check if playlist needs to be updated
		if tidalPlaylist.Data.UUID != "" && (tidalPlaylist.Data.Title != spotifyPlaylist.Name && spotifyPlaylist.Name != "") || tidalPlaylist.Data.Description != playlistDescription {
			log.Info().Msgf("updating tidal playlist: %s", spotifyPlaylist.Name)
			err := tidalClient.UpdatePlaylist(tidalPlaylist.Data.UUID, spotifyPlaylist.Name, playlistDescription)
			if err != nil {
				return err
			}
		}
	}

	// Print mapping
	fmt.Println("Spotify playlist to Tidal playlist mapping")
	for _, spotifyPlaylist := range spotifyPlaylists {
		checkString := spotifyPlaylist.ID.String()
		for _, tidtidalPlaylist := range tidalPlaylists {
			if strings.Contains(tidtidalPlaylist.Data.Description, checkString) {
				fmt.Printf("%s:%s\n", spotifyPlaylist.ID.String(), tidtidalPlaylist.Data.UUID)
				break
			}
		}
	}

	return nil
}
