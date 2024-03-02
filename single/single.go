package single

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

type Single struct {
	Filename  string   `toml:"filename"`
	CoverArt  string   `toml:"coverart"`
	Name      string   `toml:"name"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Featuring []string `toml:"featuring"`
}

type SingleConfig struct {
	Artist  string   `toml:"artist"`
	Singles []Single `toml:"singles"`
}

func Execute(src string, dst string) error {
	// srcDir, _ := cmd.Flags().GetString("src")

	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	conf := path.Join(src, "singles.toml")

	data, err := os.ReadFile(conf)
	if err != nil {
		return err
	}

	var config SingleConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	artistName := strings.TrimSpace(config.Artist)

	dstDir := path.Join(dst, artistName)
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}

	pretty.Println(config)

	// TODO(patrik): Check albumName for forward slashes and other illegal
	// filesystem characters
	for _, track := range config.Singles {
		albumName := track.Name + " (Single)"

		dir := path.Join(dstDir, albumName)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		srcCoverArt := path.Join(src, track.CoverArt)
		ext := path.Ext(srcCoverArt)
		_, err = utils.Copy(srcCoverArt, path.Join(dir, "cover"+ext))
		if err != nil {
			return err
		}

		args := []string{}

		trackPath := path.Join(src, track.Filename)

		args = append(args, "-i", trackPath)

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
			return err
		}
	}

	return nil
}
