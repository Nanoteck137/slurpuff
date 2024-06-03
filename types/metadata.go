package types

type OldTrackMetadata struct {
	Filename  string   `toml:"filename"`
	Num       int      `toml:"num"`
	Name      string   `toml:"name"`
	Artist    string   `toml:"artist"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Genres    []string `toml:"genres"`
	Featuring []string `toml:"featuring"`
}

type OldAlbumMetadata struct {
	Album    string             `toml:"album"`
	Artist   string             `toml:"artist"`
	CoverArt string             `toml:"coverart"`
	Tracks   []OldTrackMetadata `toml:"tracks"`
}

type TrackFile struct {
	Lossless string `toml:"lossless"`
	Lossy    string `toml:"lossy"`
}

type TrackMetadata struct {
	Num       int       `toml:"num"`
	Name      string    `toml:"name"`
	Duration  int       `toml:"duration"`
	Artist    string    `toml:"artist"`
	Year      int       `toml:"year"`
	Tags      []string  `toml:"tags"`
	Genres    []string  `toml:"genres"`
	Featuring []string  `toml:"featuring"`
	File      TrackFile `toml:"file,inline"`
}

type AlbumMetadata struct {
	Album    string          `toml:"album"`
	Artist   string          `toml:"artist"`
	CoverArt string          `toml:"coverart"`
	Tracks   []TrackMetadata `toml:"tracks"`
}
