package cmd

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/parasect"
	"github.com/nanoteck137/slurpuff/types"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func Convert(p string) {
	albumPath := path.Join(p, "album.toml")
	log.Printf("Converting '%s'", albumPath) 

	data, err := os.ReadFile(albumPath)
	if err != nil {
		log.Fatal(err)
	}

	var old types.OldAlbumMetadata
	err = toml.Unmarshal(data, &old)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(old)

	var tracks []types.TrackMetadata
	for _, t := range old.Tracks {
		year, err := strconv.Atoi(t.Date)
		if err != nil {
			log.Fatal(err)
		}

		trackFile := path.Join(p, t.Filename)

		info, err := parasect.GetTrackInfo(trackFile)
		if err != nil {
			log.Fatal(err)
		}

		lossless := t.Filename
		lossy := ""

		ext := path.Ext(trackFile)

		if utils.IsLossyFormatExt(ext) {
			lossless = ""
			lossy = t.Filename
		} else {
			name := strings.TrimSuffix(t.Filename, ext) + ".opus"
			dst := path.Join(p, name)

			// TODO(patrik): Add options for this
			err = parasect.RunFFmpeg(true, "-y", "-i", trackFile, "-vbr", "on", "-b:a", "128k", dst)
			if err != nil {
				log.Fatal(err)
			}

			lossy = name
		}

		tracks = append(tracks, types.TrackMetadata{
			Num:       t.Num,
			Name:      t.Name,
			Duration:  info.Duration,
			Artist:    t.Artist,
			Year:      year,
			Tags:      t.Tags,
			Genres:    t.Genres,
			Featuring: t.Featuring,
			File: types.TrackFile{
				Lossless: lossless,
				Lossy:    lossy,
			},
		})
	}

	var metadata types.AlbumMetadata
	metadata.Album = old.Album
	metadata.Artist = old.Artist
	metadata.CoverArt = old.CoverArt
	metadata.Tracks = tracks

	pretty.Println(metadata)

	d, err := toml.Marshal(old)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(path.Join(p, "old_album.toml"), d, 0644)
	if err != nil {
		log.Fatal(err)
	}

	d, err = toml.Marshal(metadata)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(albumPath, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

var convertCmd = &cobra.Command{
	Use: "convert",
	Run: func(cmd *cobra.Command, args []string) {
		recursive, _ := cmd.Flags().GetBool("recursive")

		if recursive {
			filepath.WalkDir(".", func(p string, d fs.DirEntry, err error) error {
				if d.Name() == "album.toml" {
					Convert(path.Dir(p))
				}

				return nil
			})
		} else {
			Convert(".")
		}
	},
}

func init() {
	convertCmd.Flags().BoolP("recursive", "r", false, "Recursively search for 'album.toml' to convert")

	rootCmd.AddCommand(convertCmd)
}
