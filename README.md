# music-utils

A collection of miscellaneous utilities for my music setup.

## Installation

A Spotify and Tidal developer application is required and the client ID and secret for both.

Use the provided `compose.yml` file to setup the container. Then run `docker compose run music-utils` to interact with the CLI.

## Commands

- `tidal`
  - `save` - Save all user's Tidal playlists to a JSON file
  - `print` - Print all user's Tidal playlists
  - `links` - Print all user's Tidal playlist links
- `spotify`
  - `save` - Save all user's Spotify playlists to a JSON file
  - `print` - Print all user's Spotify playlists
  - `create-tidal-playlists` - Creates **empty** Tidal playlists from Spotify playlists. Use [spotify-to-tidal](https://github.com/spotify2tidal/spotify_to_tidal) to convert the tracks.
