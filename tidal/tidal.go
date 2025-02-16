package tidal

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zibbp/music-utils/config"
	"github.com/zibbp/music-utils/navidrome"
)

var (
	apiURL  = "https://api.tidal.com/v1"
	apiURL2 = "https://listen.tidal.com/v2"
)

type Service struct {
	ClientId          string
	ClientSecret      string
	AccessToken       string // device-flow user resources access token
	ClientAccessToken string // application client for accessing Tidal API non-user resources
	UserID            string
	Config            *config.JsonConfigService
}

func Initialize(clientId, clientSecret string, config *config.JsonConfigService) (*Service, error) {
	if clientId == "" || clientSecret == "" {
		log.Fatal().Msg("Tidal client ID and secret not provided, check env vars")
	}

	var s Service
	s.ClientId = clientId
	s.ClientSecret = clientSecret
	s.Config = config

	if s.Config.Get().Tidal.AccessToken != "" {
		s.AccessToken = s.Config.Get().Tidal.AccessToken
	}

	// client auth is non-interactive so run always
	clientAccessToken, err := s.clientAuth(s.ClientId, s.ClientSecret)
	if err != nil {
		return nil, err
	}

	s.ClientAccessToken = clientAccessToken

	err = s.DeviceAuthenticate()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to authenticate with Tidal")
	}

	return &s, nil
}

// Perform device authentication with Tidal to access user resources
func (s *Service) DeviceAuthenticate() error {
	if s.Config.Get().Tidal.AccessToken == "" || s.Config.Get().Tidal.RefreshToken == "" {
		log.Info().Msg("No Tidal access token found")

		deviceCode, err := s.getDeviceCode()
		if err != nil {
			return err
		}

		log.Info().Msgf("Please visit the following URL to authorize this application: https://%v", deviceCode.VerificationURIComplete)

		// start poll for authorization
		for {
			loginResponse, err := s.tokenLogin(*deviceCode)
			if err != nil {
				// continue polling
				log.Debug().Msg("Failed to login with Tidal")

			}

			if (loginResponse != nil && AuthLogin{} == loginResponse.AuthLogin) {
				if loginResponse.AuthError.Error == "expired_token" {
					log.Fatal().Msg("Tidal auth failed - device code expired")
				}
			} else {
				s.Config.JsonConfig.Tidal.UserID = strconv.Itoa(int(loginResponse.AuthLogin.User.UserID))
				s.Config.JsonConfig.Tidal.AccessToken = loginResponse.AuthLogin.AccessToken
				s.Config.JsonConfig.Tidal.RefreshToken = loginResponse.AuthLogin.RefreshToken
				s.Config.Save()
				break
			}

			d := time.Duration(deviceCode.Interval) * time.Second
			log.Debug().Msgf("Waiting %d seconds before trying again.", deviceCode.Interval)
			time.Sleep(d)

		}
	} else {
		log.Info().Msg("Tidal access token found")
		_, err := s.checkSession(s.Config.Get().Tidal.AccessToken)
		if err != nil {
			// failed probably need to refresh
			log.Info().Msg("Tidal access token expired")
			refresh, err := s.refreshSession(s.Config.Get().Tidal.RefreshToken)
			if err != nil {
				return err
			}

			s.Config.JsonConfig.Tidal.AccessToken = refresh.AccessToken
			s.Config.Save()

		}

		log.Info().Msg("Tidal access token valid")
	}

	s.AccessToken = s.Config.Get().Tidal.AccessToken
	s.UserID = s.Config.Get().Tidal.UserID

	return nil
}

func (s *Service) ToNavidromePlaylist(playlist *Playlist) (navidrome.Playlist, error) {
	var np navidrome.Playlist

	np.SourceId = playlist.UUID
	np.Name = playlist.Title
	np.Description = playlist.Description
	np.Tracks = make([]navidrome.Track, 0)

	for _, track := range playlist.Tracks {
		np.Tracks = append(np.Tracks, navidrome.Track{
			ID:       strconv.FormatInt(track.ID, 10),
			Title:    track.Title,
			Artist:   track.Artists[0].Name,
			Album:    track.Album.Title,
			Duration: track.Duration,
			ISRC:     track.Isrc,
		})
	}

	return np, nil
}
