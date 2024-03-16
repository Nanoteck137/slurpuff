package cmd

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/kr/pretty"
	dwebbleutils "github.com/nanoteck137/dwebble/utils"
	"github.com/nanoteck137/slurpuff/album"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
}

var initAlbumCmd = &cobra.Command{
	Use: "album",
	Run: func(cmd *cobra.Command, args []string) {
		src, _ := cmd.Flags().GetString("dir")

		outputFile, _ := cmd.Flags().GetString("output")

		entries, err := os.ReadDir(src)
		if err != nil {
			log.Fatal(err)
		}

		albumArtist := ""
		albumName := ""

		tracks := []album.Track{}

		for _, entry := range entries {
			if entry.Name()[0] == '.' {
				continue
			}

			p := path.Join(src, entry.Name())
			ext := path.Ext(entry.Name())

			if ext == "" {
				continue
			}

			// TODO(patrik): Change this IsValidTrackExt
			if dwebbleutils.IsValidTrackExt(ext[1:]) {
				res, err := dwebbleutils.CheckFile(p)
				if err != nil {
					log.Fatal(err)
				}

				pretty.Println(res)

				if albumName == "" {
					if name, exists := res.Tags["album"]; exists {
						albumName = name
					}
				}

				if albumArtist == "" {
					if name, exists := res.Tags["album_artist"]; exists {
						albumArtist = name
					}
				}

				artist := ""
				if value, exists := res.Tags["artist"]; exists {
					artist = value
				}

				name := res.Name
				if value, exists := res.Tags["title"]; exists {
					name = value
				}

				date := ""
				if value, exists := res.Tags["date"]; exists {
					date = value
				}

				artists := strings.Split(artist, ",")
				for i := range artists {
					artists[i] = strings.TrimSpace(artists[i])
				}

				tracks = append(tracks, album.Track{
					Filename:  entry.Name(),
					Num:       res.Number,
					Name:      name,
					Date:      date,
					Artist:    artists[0],
					Tags:      []string{},
					Featuring: artists[1:],
				})
			}
		}

		config := album.AlbumConfig{
			Album:    albumName,
			Artist:   albumArtist,
			CoverArt: "",
			Tracks:   tracks,
		}

		data, err := toml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initCmd.PersistentFlags().StringP("dir", "d", ".", "album directory")

	initAlbumCmd.Flags().StringP("output", "o", "album.toml", "output file")

	initCmd.AddCommand(initAlbumCmd)
	rootCmd.AddCommand(initCmd)
}
