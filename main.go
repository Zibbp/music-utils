package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/zibbp/music-utils/commands"
	"github.com/zibbp/music-utils/config"
)

func initialize() (*config.Config, *config.JsonConfigService) {
	if os.Getenv("JSON_LOG") != "true" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// initialize config
	c, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// load json config which has credentials
	jsonConfig := config.NewJsonConfigService("/data/config.json")
	err = jsonConfig.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load Spotify config")
	}

	return c, jsonConfig
}

func main() {
	app := &cli.App{
		Name:  "music-utils",
		Usage: "A collection of small utlities for my music setup.",
		Commands: []*cli.Command{
			{
				Name:  "tidal",
				Usage: "Tidal related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "save",
						Usage: "Save all user tidal playlists to a JSON file",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "save-navidrome-format",
								Usage: "Save a version of the tidal playlist in a format for importing into Navidrome",
							},
						},
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							saveNavidromeFormat := cCtx.Bool("save-navidrome-format")

							err := commands.SaveTidalPlaylists(cCtx.Context, c, jsonConfig, saveNavidromeFormat)
							if err != nil {
								log.Fatal().Err(err).Msg("error saving tidal playlists")
							}

							return nil
						},
					},
					{
						Name:  "print",
						Usage: "Print all user tidal playlists",
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							err := commands.PrintTidalPlaylists(cCtx.Context, c, jsonConfig)
							if err != nil {
								log.Fatal().Err(err).Msg("error printing tidal playlists")
							}

							return nil
						},
					},
					{
						Name:  "links",
						Usage: "Print all user tidal playlist links",
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							err := commands.PrintTidalPlaylistsLinks(cCtx.Context, c, jsonConfig)
							if err != nil {
								log.Fatal().Err(err).Msg("error printing tidal playlists")
							}

							return nil
						},
					},
				},
			},
			{
				Name:  "spotify",
				Usage: "Spotify related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "save",
						Usage: "Save all user spotify playlists to a JSON file",
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							err := commands.SaveSpotifyPlaylists(cCtx.Context, c, jsonConfig)
							if err != nil {
								log.Fatal().Err(err).Msg("error saving spotify playlists")
							}

							return nil
						},
					},
					{
						Name:  "print",
						Usage: "Print all user spotify playlist IDs",
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							err := commands.PrintSpotifyPlaylists(cCtx.Context, c, jsonConfig)
							if err != nil {
								log.Fatal().Err(err).Msg("error printing spotify playlists")
							}

							return nil
						},
					},
					{
						Name:  "create-tidal-playlists",
						Usage: "Create Spotify playlists on Tidal. Used in conjunction with https://github.com/spotify2tidal/spotify_to_tidal",
						Action: func(cCtx *cli.Context) error {
							c, jsonConfig := initialize()

							err := commands.CreateSpotifyPlaylistsOnTidal(cCtx.Context, c, jsonConfig)
							if err != nil {
								log.Fatal().Err(err).Msg("error printing spotify playlists")
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}
