package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/pelletier/go-toml/v2"
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

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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

		srcCoverArt := path.Join(srcDir, track.CoverArt)
		ext := path.Ext(srcCoverArt)
		_, err = copy(srcCoverArt, path.Join(dir, "cover"+ext))
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
}
