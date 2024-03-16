package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"
	dwebbleutils "github.com/nanoteck137/dwebble/utils"
	"github.com/nanoteck137/slurpuff/album"
	"github.com/nanoteck137/slurpuff/single"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(1),
}

var singleCmd = &cobra.Command{
	Use:  "singles <DEST_DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		src, _ := cmd.Flags().GetString("src")

		err := single.Execute(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var albumCmd = &cobra.Command{
	Use:  "album <DEST_DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		src, _ := cmd.Flags().GetString("src")

		err := album.Execute(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var allCmd = &cobra.Command{
	Use: "all",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		src, _ := cmd.Flags().GetString("src")

		fmt.Printf("dst: %v\n", dst)
		fmt.Printf("src: %v\n", src)

		var albums []string
		var singles []string

		filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
			if d.Name() == "album.toml" {
				albums = append(albums, p)
			}

			if d.Name() == "singles.toml" {
				singles = append(singles, p)
			}

			return nil
		})

		fmt.Printf("albums: %v\n", albums)
		fmt.Printf("singles: %v\n", singles)

		for _, p := range albums {
			src := path.Dir(p)
			err := album.Execute(src, dst)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, p := range singles {
			src := path.Dir(p)
			err := single.Execute(src, dst)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var createCmd = &cobra.Command{
	Use:  "create <SRC>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]

		outputFile, _ := cmd.Flags().GetString("out")
		if outputFile == "" {
			outputFile = path.Join(src, "album.toml")
		}

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
	singleCmd.Flags().StringP("src", "s", ".", "directory with singles.toml")

	albumCmd.Flags().StringP("src", "s", ".", "directory with album.toml")

	allCmd.Flags().StringP("src", "s", ".", "directory to search for music")

	createCmd.Flags().StringP("out", "o", "", "output file (defaults to \"<SRC_DIR>/album.toml\"")

	convertCmd.AddCommand(singleCmd)
	convertCmd.AddCommand(albumCmd)
	convertCmd.AddCommand(allCmd)
	convertCmd.AddCommand(createCmd)

	rootCmd.AddCommand(convertCmd)
}
