package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/kr/pretty"
	"github.com/pelletier/go-toml/v2"
)

type Track struct {
	Filename string `toml:"filename"`
	Name     string `toml:"name"`
}

type Config struct {
	Artist string `toml:"artist"` 
	Tracks []Track `toml:"tracks"`
}

func main() {
	dir := "./result"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	srcDir := "/Volumes/media/musicraw/Divide Music"

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

	for _, track := range config.Tracks {
		cmd := exec.Command("ffmpeg", 
			"-i", path.Join(srcDir, track.Filename), 
			"-metadata", fmt.Sprintf(`title=%s`, track.Name), 
			"-metadata", fmt.Sprintf(`artist=%s`, config.Artist), 
			"-metadata", fmt.Sprintf(`album_artist=%s`, config.Artist), 
			"-metadata", fmt.Sprintf(`album=%s`, track.Name + " (Single)"), 
			path.Join(dir, track.Name+".flac"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
