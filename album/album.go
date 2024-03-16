package album

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
)

type Track struct {
	Filename  string   `toml:"filename"`
	Num       int      `toml:"num"`
	Name      string   `toml:"name"`
	Date      string   `toml:"date"`
	Artist    string   `toml:"artist"`
	Tags      []string `toml:"tags"`
	Featuring []string `toml:"featuring,omitempty"`
}

type AlbumConfig struct {
	Album    string  `toml:"album"`
	Artist   string  `toml:"artist"`
	CoverArt string  `toml:"coverart"`
	Tracks   []Track `toml:"tracks"`
}

func Execute(src string, dst string) error {
	// TODO(patrik): Add force flag

	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	conf := path.Join(src, "album.toml")

	data, err := os.ReadFile(conf)
	if err != nil {
		return err
	}

	var config AlbumConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	err = ExecuteConfig(config, src, dst)
	if err != nil {
		return err
	}

	return nil
}

func ExecuteConfig(config AlbumConfig, src, dst string) error {
	artistName := strings.TrimSpace(config.Artist)

	safeArtistName, err := utils.SafeName(artistName)
	if err != nil {
		return err
	}

	dstDir := path.Join(dst, safeArtistName)
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}

	if artistName != safeArtistName {
		err := os.WriteFile(path.Join(dstDir, "override.txt"), []byte(artistName), 0644)
		if err != nil {
			return err
		}
	}

	pretty.Println(config)

	albumName := config.Album

	safeAlbumName, err := utils.SafeName(albumName)
	if err != nil {
		return err
	}

	dir := path.Join(dstDir, safeAlbumName)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	if albumName != safeAlbumName {
		// TODO(patrik): Dont override
		err := os.WriteFile(path.Join(dir, "override.txt"), []byte(albumName), 0644)
		if err != nil {
			return err
		}
	}

	if config.CoverArt != "" {
		srcCoverArt := path.Join(src, config.CoverArt)
		ext := path.Ext(srcCoverArt)
		_, err = utils.Copy(srcCoverArt, path.Join(dir, "cover"+ext))
		if err != nil {
			return err
		}
	}

	// TODO(patrik): Check albumName for forward slashes and other illegal
	// filesystem characters
	for _, track := range config.Tracks {
		args := []string{}

		trackPath := path.Join(src, track.Filename)

		inputExt := path.Ext(trackPath)
		outputExt := ".opus"

		if inputExt == ".wav" {
			outputExt = ".flac"
		}

		copyMode := inputExt == outputExt

		args = append(args, "-i", trackPath, "-vn", "-map_metadata", "-1")

		artist := config.Artist
		if track.Artist != "" {
			artist = track.Artist
		}

		args = append(args, "-metadata", fmt.Sprintf("title=%s", track.Name))
		args = append(args, "-metadata", fmt.Sprintf("artist=%s", artist))
		args = append(args, "-metadata", fmt.Sprintf("album_artist=%s", config.Artist))
		args = append(args, "-metadata", fmt.Sprintf("album=%s", albumName))
		args = append(args, "-metadata", fmt.Sprintf("track=%d", track.Num))

		if len(track.Tags) > 0 {
			args = append(args, "-metadata", fmt.Sprintf("tags=%s", strings.Join(track.Tags, ",")))
		}

		if track.Date != "" {
			args = append(args, "-metadata", fmt.Sprintf("date=%s", track.Date))
		}

		if len(track.Featuring) > 0 {
			args = append(args, "-metadata", fmt.Sprintf("featuring=%s", strings.Join(track.Featuring, ",")))
		}

		if !copyMode && outputExt == ".opus" {
			args = append(args, "-vbr", "on", "-b:a", "128k")
		}

		if copyMode {
			args = append(args, "-codec", "copy")
		}

		outputName := fmt.Sprintf("%02v - %s%s", track.Num, strings.TrimSpace(track.Name), outputExt)
		safeOutputName, err := utils.SafeName(outputName)
		if err != nil {
			return err
		}
		args = append(args, path.Join(dir, safeOutputName))

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
