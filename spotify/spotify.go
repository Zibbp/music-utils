package spotify

import (
	"context"
	"fmt"

	"github.com/zibbp/music-utils/config"

	spotifyPkg "github.com/zmb3/spotify/v2"
)

type Service struct {
	client            *spotifyPkg.Client
	config            *config.JsonConfigService
	clientId          string
	clientSecret      string
	clientRedirectUri string
}

func Initialize(clientId, clientSecret, clientRedirectUri string, config *config.JsonConfigService) (*Service, error) {
	if clientId == "" || clientSecret == "" || clientRedirectUri == "" {
		return nil, fmt.Errorf("spotify client ID, secret and redirect URI not provided, check env vars")
	}

	var s Service
	s.clientId = clientId
	s.clientSecret = clientSecret
	s.clientRedirectUri = clientRedirectUri
	s.config = config

	err := s.Authenticate()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Service) Authenticate() error {
	client, err := s.authFlow()
	if err != nil {
		return err
	}

	s.client = client

	return nil
}

func (s *Service) GetUserPlaylists() ([]spotifyPkg.SimplePlaylist, error) {
	playlists, err := s.client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		return nil, err
	}

	var allPlaylists []spotifyPkg.SimplePlaylist
	for page := 1; ; page++ {
		allPlaylists = append(allPlaylists, playlists.Playlists...)
		if playlists.Next == "" {
			break
		}

		err = s.client.NextPage(context.Background(), playlists)
		if err != nil {
			return nil, err
		}
	}

	return allPlaylists, nil
}

func (s *Service) GetPlaylistTracks(id spotifyPkg.ID) ([]*spotifyPkg.FullTrack, error) {
	items, err := s.client.GetPlaylistItems(context.Background(), id)
	if err != nil {
		return nil, err
	}

	var allTracks []*spotifyPkg.FullTrack
	for page := 1; ; page++ {
		for _, track := range items.Items {
			allTracks = append(allTracks, track.Track.Track)
		}
		if items.Next == "" {
			break
		}

		err = s.client.NextPage(context.Background(), items)
		if err != nil {
			return nil, err
		}
	}

	return allTracks, nil
}
