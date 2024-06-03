package types

type TrackMetadata struct {
	Filename  string   `toml:"filename"`
	Num       int      `toml:"num"`
	Name      string   `toml:"name"`
	Artist    string   `toml:"artist"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Genres    []string `toml:"genres"`
	Featuring []string `toml:"featuring"`
}

type AlbumMetadata struct {
	Album    string          `toml:"album"`
	Artist   string          `toml:"artist"`
	CoverArt string          `toml:"coverart"`
	Tracks   []TrackMetadata `toml:"tracks"`
}
