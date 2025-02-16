package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/zibbp/music-utils/config"
	"github.com/zibbp/music-utils/tidal"
	"github.com/zibbp/music-utils/utils"
)

func SaveTidalPlaylists(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService, saveNavidromeFormat bool) error {
	// Initialize Tidal client
	tidalClient, err := tidal.Initialize(envConfig.TidalClientId, envConfig.TidalClientSecret, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing tidal client: %v", err)
	}

	log.Info().Msg("fetching tidal playlists")

	// Get all user's Tidal playlists
	tidalPlaylists, err := tidalClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	for _, tmpTidalPlaylist := range tidalPlaylists {
		log.Info().Str("playlist", tmpTidalPlaylist.Data.Title).Msg("processing playlist")
		// Get full playlist
		tidalPlaylist, err := tidalClient.GetPlaylist(string(tmpTidalPlaylist.Data.UUID))
		if err != nil {
			log.Error().Err(err).Msg("error getting tidal playlist")
		}

		// Get playlist tracks
		tidalPlaylistTracks, err := tidalClient.GetPlaylistTracks(string(tidalPlaylist.UUID))
		if err != nil {
			log.Error().Err(err).Msg("error getting tidal playlist tracks")
		}
		tidalPlaylist.Tracks = append(tidalPlaylist.Tracks, tidalPlaylistTracks.Items...)

		if err := utils.WriteJsonToFile("/data/tidal", tidalPlaylist.UUID, tidalPlaylist); err != nil {
			log.Error().Err(err).Msg("error writing playlist to file")
			continue
		}

		if saveNavidromeFormat {
			navidromePlaylist, err := tidalClient.ToNavidromePlaylist(tidalPlaylist)
			if err != nil {
				log.Error().Err(err).Msg("error converting tidal playlist to navidrome playlist")
				continue
			}
			if err := utils.WriteJsonToFile("/data/navidrome", tidalPlaylist.UUID+"_navidrome", navidromePlaylist); err != nil {
				log.Error().Err(err).Msg("error writing navidrome playlist to file")
				continue
			}
		}
	}

	return nil
}

func PrintTidalPlaylists(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService) error {
	// Initialize Tidal client
	tidalClient, err := tidal.Initialize(envConfig.TidalClientId, envConfig.TidalClientSecret, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing tidal client: %v", err)
	}

	// Get all user's Spotify playlists
	tidalPlaylists, err := tidalClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	data := [][]string{
		{"ID", "Name", "Description"},
	}

	for _, tidalPlaylist := range tidalPlaylists {
		data = append(data, []string{tidalPlaylist.Data.UUID, tidalPlaylist.Data.Title, tidalPlaylist.Data.Description})
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, row := range data {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()

	return nil
}

func PrintTidalPlaylistsLinks(ctx context.Context, envConfig *config.Config, jsonConfig *config.JsonConfigService) error {
	// Initialize Tidal client
	tidalClient, err := tidal.Initialize(envConfig.TidalClientId, envConfig.TidalClientSecret, jsonConfig)
	if err != nil {
		return fmt.Errorf("error initializing tidal client: %v", err)
	}

	// Get all user's Spotify playlists
	tidalPlaylists, err := tidalClient.GetUserPlaylists()
	if err != nil {
		return err
	}

	for _, tidalPlaylist := range tidalPlaylists {
		fmt.Printf("https://listen.tidal.com/playlist/%s\n", tidalPlaylist.Data.UUID)
	}

	return nil
}
