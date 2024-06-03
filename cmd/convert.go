package cmd

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/parasect"
	"github.com/nanoteck137/slurpuff/types"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use: "convert",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile("album.toml")
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

			info, err := parasect.GetTrackInfo(t.Filename)
			if err != nil {
				log.Fatal(err)
			}

			lossless := t.Filename
			lossy := ""

			ext := path.Ext(t.Filename)

			if utils.IsLossyFormatExt(ext) {
				lossless = ""
				lossy = t.Filename
			} else {
				dst := strings.TrimSuffix(t.Filename, ext) + ".opus"

				// TODO(patrik): Add options for this
				err = parasect.RunFFmpeg(true, "-y", "-i", t.Filename, "-vbr", "on", "-b:a", "128k", dst)
				if err != nil {
					log.Fatal(err)
				}

				lossy = dst
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

		err = os.WriteFile("old_album.toml", d, 0644)
		if err != nil {
			log.Fatal(err)
		}

		d, err = toml.Marshal(metadata)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile("album.toml", d, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// var singleCmd = &cobra.Command{
// 	Use:  "singles <DEST_DIR>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dst := args[0]
// 		src, _ := cmd.Flags().GetString("src")
// 		mode, _ := cmd.Flags().GetString("mode")
//
// 		err := single.Execute(mode, src, dst)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	},
// }
//
// var albumCmd = &cobra.Command{
// 	Use:  "album <DEST_DIR>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dst := args[0]
// 		src, _ := cmd.Flags().GetString("src")
// 		mode, _ := cmd.Flags().GetString("mode")
//
// 		err := album.Execute(mode, src, dst)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	},
// }
//
// var allCmd = &cobra.Command{
// 	Use:  "all <DEST_DIR>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dst := args[0]
// 		src, _ := cmd.Flags().GetString("src")
// 		mode, _ := cmd.Flags().GetString("mode")
//
// 		fmt.Printf("dst: %v\n", dst)
// 		fmt.Printf("src: %v\n", src)
//
// 		var albums []string
// 		var singles []string
//
// 		filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
// 			if d.Name() == "album.toml" {
// 				albums = append(albums, p)
// 			}
//
// 			if d.Name() == "singles.toml" {
// 				singles = append(singles, p)
// 			}
//
// 			return nil
// 		})
//
// 		fmt.Printf("albums: %v\n", albums)
// 		fmt.Printf("singles: %v\n", singles)
//
// 		for _, p := range albums {
// 			src := path.Dir(p)
// 			err := album.Execute(mode, src, dst)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
//
// 		for _, p := range singles {
// 			src := path.Dir(p)
// 			err := single.Execute(mode, src, dst)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	},
// }

func init() {
	// singleCmd.Flags().StringP("src", "s", ".", "directory with singles.toml")
	// albumCmd.Flags().StringP("src", "s", ".", "directory with album.toml")
	// allCmd.Flags().StringP("src", "s", ".", "directory to search for music")

	// convertCmd.PersistentFlags().StringP("mode", "m", "dwebble", "convertion mode (valid dwebble,opus,mp3,map)")

	// convertCmd.AddCommand(singleCmd)
	// convertCmd.AddCommand(albumCmd)
	// convertCmd.AddCommand(allCmd)

	rootCmd.AddCommand(convertCmd)
}
