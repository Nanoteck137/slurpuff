package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/nanoteck137/slurpuff/album"
	"github.com/nanoteck137/slurpuff/single"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use: "convert",
}

var singleCmd = &cobra.Command{
	Use:  "singles <DEST_DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		src, _ := cmd.Flags().GetString("src")
		mode, _ := cmd.Flags().GetString("mode")

		err := single.Execute(mode, src, dst)
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
		mode, _ := cmd.Flags().GetString("mode")

		err := album.Execute(mode, src, dst)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var allCmd = &cobra.Command{
	Use:  "all",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		src, _ := cmd.Flags().GetString("src")
		mode, _ := cmd.Flags().GetString("mode")

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

		wg := sync.WaitGroup{}
		for _, p := range albums {
			wg.Add(1)
			go func() {
				src := path.Dir(p)
				err := album.Execute(mode, src, dst)
				if err != nil {
					log.Fatal(err)
				}
				wg.Done()
			}()
		}

		for _, p := range singles {
			src := path.Dir(p)
			err := single.Execute(mode, src, dst)
			if err != nil {
				log.Fatal(err)
			}
		}

		wg.Wait()
	},
}

func init() {
	singleCmd.Flags().StringP("src", "s", ".", "directory with singles.toml")
	albumCmd.Flags().StringP("src", "s", ".", "directory with album.toml")
	allCmd.Flags().StringP("src", "s", ".", "directory to search for music")

	convertCmd.PersistentFlags().StringP("mode", "m", "dwebble", "convertion mode (valid dwebble,opus,mp3,map)")

	convertCmd.AddCommand(singleCmd)
	convertCmd.AddCommand(albumCmd)
	convertCmd.AddCommand(allCmd)

	rootCmd.AddCommand(convertCmd)
}
