package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type Track struct {
	Filename  string   `toml:"filename"`
	CoverArt  string   `toml:"coverart"`
	Name      string   `toml:"name"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Featuring []string `toml:"featuring"`
}

type Config struct {
	Artist string  `toml:"artist"`
	Tracks []Track `toml:"tracks"`
}

var convertCmd = &cobra.Command{
	Use:  "convert <DEST_DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Convert", args[0])
		srcDir, _ := cmd.Flags().GetString("src")

		dst := args[0]
		err := os.MkdirAll(dst, 0755)
		if err != nil {
			log.Fatal(err)
		}

		conf := path.Join(srcDir, "tracks.toml")

		data, err := os.ReadFile(conf)
		if err != nil {
			log.Fatal(err)
		}

		var config Config
		err = toml.Unmarshal(data, &config)
		if err != nil {
			log.Fatal(err)
		}

		artistName := strings.TrimSpace(config.Artist)

		dstDir := path.Join(dst, artistName)
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(config)

		// TODO(patrik): Check albumName for forward slashes and other illegal
		// filesystem characters
		for _, track := range config.Tracks {
			albumName := track.Name + " (Single)"

			dir := path.Join(dstDir, albumName)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}

			srcCoverArt := path.Join(srcDir, track.CoverArt)
			ext := path.Ext(srcCoverArt)
			_, err = utils.Copy(srcCoverArt, path.Join(dir, "cover"+ext))
			if err != nil {
				log.Fatal(err)
			}

			args := []string{}

			args = append(args, "-i", path.Join(srcDir, track.Filename))

			args = append(args, "-metadata", fmt.Sprintf(`title=%s`, track.Name))
			args = append(args, "-metadata", fmt.Sprintf(`artist=%s`, config.Artist))
			args = append(args, "-metadata", fmt.Sprintf(`album_artist=%s`, config.Artist))
			args = append(args, "-metadata", fmt.Sprintf(`album=%s`, albumName))

			if len(track.Tags) > 0 {
				args = append(args, "-metadata", fmt.Sprintf(`tags=%s`, strings.Join(track.Tags, ",")))
			}

			if track.Date != "" {
				args = append(args, "-metadata", fmt.Sprintf(`date=%s`, track.Date))
			}

			if len(track.Featuring) > 0 {
				args = append(args, "-metadata", fmt.Sprintf(`featuring=%s`, strings.Join(track.Featuring, ",")))
			}

			args = append(args, path.Join(dir, "01 - "+strings.TrimSpace(track.Name)+".flac"))

			cmd := exec.Command("ffmpeg", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	convertCmd.Flags().StringP("src", "s", ".", "directory with tracks.toml")

	rootCmd.AddCommand(convertCmd)
}
