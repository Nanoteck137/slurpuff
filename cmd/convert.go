package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"

	"github.com/nanoteck137/slurpuff/album"
	"github.com/nanoteck137/slurpuff/single"
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

func init() {
	singleCmd.Flags().StringP("src", "s", ".", "directory with singles.toml")

	albumCmd.Flags().StringP("src", "s", ".", "directory with album.toml")

	allCmd.Flags().StringP("src", "s", ".", "directory to search for music")

	convertCmd.AddCommand(singleCmd)
	convertCmd.AddCommand(albumCmd)
	convertCmd.AddCommand(allCmd)

	rootCmd.AddCommand(convertCmd)
}
