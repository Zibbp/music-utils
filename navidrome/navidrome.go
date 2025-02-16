package navidrome

type Playlist struct {
	SourceId      string  `json:"source_id"`
	DestinationId string  `json:"destination_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Tracks        []Track `json:"tracks"`
}

type Track struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Duration int64  `json:"duration"`
	ISRC     string `json:"isrc"`
}
