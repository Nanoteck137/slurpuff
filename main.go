package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/pelletier/go-toml/v2"
)

type Track struct {
	Filename string   `toml:"filename"`
	CoverArt string   `toml:"coverart"`
	Name     string   `toml:"name"`
	Date     string   `toml:"date"`
	Tags     []string `toml:"tags"`
}

type Config struct {
	Artist string  `toml:"artist"`
	Tracks []Track `toml:"tracks"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <SOURCE_DIR> <DEST_DIR>\n", os.Args[0])
		os.Exit(-1)
	}

	srcDir := os.Args[1]
	dstDir := os.Args[2]

	data, err := os.ReadFile(path.Join(srcDir, "tracks.toml"))
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(config)

	// TODO(patrik): Check albumName for forward slashes and other illegal
	// filesystem characters
	for _, track := range config.Tracks {
		fmt.Printf("track.Name: %v\n", track.Name)

		albumName := track.Name + " (Single)"

		dir := path.Join(dstDir, albumName)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		src := path.Join(srcDir, track.CoverArt)
		srcFile, err := os.Open(src)
		if err != nil {
			log.Fatal(err)
		}

		_ = srcFile

		args := []string{}

		args = append(args, "-i", path.Join(srcDir, track.Filename))

		args = append(args, "-metadata", fmt.Sprintf(`title=%s`, track.Name))
		args = append(args, "-metadata", fmt.Sprintf(`artist=%s`, config.Artist))
		args = append(args, "-metadata", fmt.Sprintf(`album_artist=%s`, config.Artist))
		args = append(args, "-metadata", fmt.Sprintf(`album=%s`, albumName))
		args = append(args, "-metadata", fmt.Sprintf(`tags=%s`, strings.Join(track.Tags, ",")))
		if track.Date != "" {
			args = append(args, "-metadata", fmt.Sprintf(`date=%s`, track.Date))
		}

		args = append(args, path.Join(dir, track.Name+".flac"))

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
