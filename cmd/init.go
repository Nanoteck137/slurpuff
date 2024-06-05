package cmd

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/nanoteck137/parasect"
	"github.com/nanoteck137/slurpuff/types"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",

	Run: func(cmd *cobra.Command, args []string) {
		src, _ := cmd.Flags().GetString("dir")

		outputFile, _ := cmd.Flags().GetString("output")

		genres, _ := cmd.Flags().GetString("genres")
		tags, _ := cmd.Flags().GetString("tags")
		yearOverride, _ := cmd.Flags().GetInt("year")

		entries, err := os.ReadDir(src)
		if err != nil {
			log.Fatal(err)
		}

		albumArtist := ""
		albumName := ""

		var defaultGenres []string
		if genres != "" {
			defaultGenres = strings.Split(genres, ",")
			for i, genre := range defaultGenres {
				defaultGenres[i] = strings.TrimSpace(genre)
			}
		}

		var defaultTags []string
		if tags != "" {
			defaultTags = strings.Split(tags, ",")
			for i, tag := range defaultTags {
				defaultTags[i] = strings.TrimSpace(tag)
			}
		}

		var tracks []types.TrackMetadata
		for _, entry := range entries {
			if entry.Name()[0] == '.' {
				continue
			}

			p := path.Join(src, entry.Name())
			ext := path.Ext(entry.Name())

			if ext == "" || !utils.IsValidTrackExt(ext) {
				continue
			}

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

			year := time.Now().Year()

			if yearOverride != 0 {
				year = yearOverride
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

			lossless := entry.Name()
			lossy := ""

			if utils.IsLossyFormatExt(ext) {
				lossless = ""
				lossy = entry.Name()
			} else {
				dst := strings.TrimSuffix(entry.Name(), ext) + ".opus"

				// TODO(patrik): Add options for this
				err = parasect.RunFFmpeg(true, "-y", "-i", entry.Name(), "-vbr", "on", "-b:a", "128k", dst)
				if err != nil {
					log.Fatal(err)
				}

				lossy = dst
			}

			tracks = append(tracks, types.TrackMetadata{
				Num:       int(track),
				Name:      name,
				Duration:  info.Duration,
				Artist:    artist,
				Year:      year,
				Tags:      defaultTags,
				Genres:    genres,
				Featuring: artists[1:],
				File: types.TrackFile{
					Lossless: lossless,
					Lossy:    lossy,
				},
			})
		}

		if albumArtist == "" && len(tracks) > 0 {
			albumArtist = tracks[0].Artist
		}

		albumCover := utils.FindFirstValidImage(src)

		config := types.AlbumMetadata{
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

	initCmd.Flags().StringP("output", "o", "album.toml", "output file")
	initCmd.Flags().String("genres", "", "set genres (comma seperated list)")
	initCmd.Flags().String("tags", "", "set tags (comma seperated list)")
	initCmd.Flags().Int("year", 0, "override year")

	rootCmd.AddCommand(initCmd)
}
