package single

import (
	"os"
	"path"

	"github.com/nanoteck137/slurpuff/album"
	"github.com/pelletier/go-toml/v2"
)

type Single struct {
	Filename  string   `toml:"filename"`
	CoverArt  string   `toml:"coverart"`
	Name      string   `toml:"name"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Featuring []string `toml:"featuring"`
}

type SingleConfig struct {
	Artist  string   `toml:"artist"`
	Singles []Single `toml:"singles"`
}

func Execute(mode, src, dst string) error {
	// srcDir, _ := cmd.Flags().GetString("src")

	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	conf := path.Join(src, "singles.toml")

	data, err := os.ReadFile(conf)
	if err != nil {
		return err
	}

	var config SingleConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	for _, single := range config.Singles {
		_ = single

		albumConfig := album.AlbumConfig{
			Album:    single.Name + " (Single)",
			Artist:   config.Artist,
			CoverArt: single.CoverArt,
			Tracks:   []album.Track{
				{
					Filename:  single.Filename,
					Num:       1,
					Name:      single.Name,
					Date:      single.Date,
					Tags:      single.Tags,
					Featuring: single.Featuring,
				},
			},
		}

		err := album.ExecuteConfig(albumConfig, mode, src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}
