package cmd

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/nanoteck137/slurpuff/album"
	"github.com/nanoteck137/slurpuff/utils"
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

		genres, _ := cmd.Flags().GetString("genres")
		tags, _ := cmd.Flags().GetString("tags")
		dateOverride, _ := cmd.Flags().GetString("date")

		entries, err := os.ReadDir(src)
		if err != nil {
			log.Fatal(err)
		}

		albumArtist := ""
		albumName := ""

		defaultGenres := strings.Split(genres, ",")
		defaultTags := strings.Split(tags, ",")

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
			if utils.IsValidTrackExt(ext[1:]) {
				info, err := utils.CheckFile(p)
				if err != nil {
					log.Fatal(err)
				}

				if albumName == "" {
					if name, exists := info.Tags["album"]; exists {
						albumName = name
					}
				}

				if albumArtist == "" {
					if name, exists := info.Tags["album_artist"]; exists {
						albumArtist = name
					}
				}

				artist := ""
				if value, exists := info.Tags["artist"]; exists {
					artist = value
				}

				track := info.Number
				if value, exists := info.Tags["track"]; exists {
					if track == 0 {
						t, _ := strconv.ParseInt(value, 10, 64)
						track = int(t)
					}
				}

				name := info.Name
				if value, exists := info.Tags["title"]; exists {
					name = value
				} else {
					if name == "" {
						name = entry.Name()
					}
				}

				date := ""
				if dateOverride != "" {
					date = dateOverride
				} else {
					if value, exists := info.Tags["date"]; exists {
						date = value
					}
				}

				var genres []string = make([]string, len(defaultGenres))
				copy(genres, defaultGenres)
				if value, exists := info.Tags["genres"]; exists {
					genres = strings.Split(value, ",")
					for i := range genres {
						genres[i] = strings.TrimSpace(genres[i])
					}
				}

				artists := strings.Split(artist, ",")
				for i := range artists {
					artists[i] = strings.TrimSpace(artists[i])
				}

				tracks = append(tracks, album.Track{
					Filename:  entry.Name(),
					Num:       int(track),
					Name:      name,
					Date:      date,
					Artist:    artists[0],
					Tags:      defaultTags,
					Genres:    genres,
					Featuring: artists[1:],
				})
			}
		}

		if albumArtist == "" && len(tracks) > 0 {
			albumArtist = tracks[0].Artist
		}

		albumCover := utils.FindFirstValidImage(src)

		config := album.AlbumConfig{
			Album:    albumName,
			Artist:   albumArtist,
			CoverArt: albumCover,
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
	initAlbumCmd.Flags().String("genres", "", "set genres (comma seperated list)")
	initAlbumCmd.Flags().String("tags", "", "set tags (comma seperated list)")
	initAlbumCmd.Flags().String("date", "", "override date")

	initCmd.AddCommand(initAlbumCmd)
	rootCmd.AddCommand(initCmd)
}
