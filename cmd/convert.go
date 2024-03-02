package cmd

import (
	"log"

	"github.com/nanoteck137/slurpuff/album"
	"github.com/nanoteck137/slurpuff/single"
	"github.com/spf13/cobra"
)


var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(1),
}

var singleCmd = &cobra.Command{
	Use: "singles <DEST_DIR>",
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
	Use: "album <DEST_DIR>",
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

func init() {
	singleCmd.Flags().StringP("src", "s", ".", "directory with singles.toml")

	albumCmd.Flags().StringP("src", "s", ".", "directory with album.toml")

	convertCmd.AddCommand(singleCmd)
	convertCmd.AddCommand(albumCmd)

	rootCmd.AddCommand(convertCmd)
}
