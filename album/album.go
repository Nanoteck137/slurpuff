package album

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/kr/pretty"
	"github.com/nanoteck137/slurpuff/utils"
	"github.com/pelletier/go-toml/v2"
)

type Track struct {
	Filename  string   `toml:"filename"`
	Num       int      `toml:"num"`
	Name      string   `toml:"name"`
	Artist    string   `toml:"artist"`
	Date      string   `toml:"date"`
	Tags      []string `toml:"tags"`
	Genres    []string `toml:"genres"`
	Featuring []string `toml:"featuring"`
}

type AlbumConfig struct {
	Album    string  `toml:"album"`
	Artist   string  `toml:"artist"`
	CoverArt string  `toml:"coverart"`
	Tracks   []Track `toml:"tracks"`
}

func Execute(mode, src, dst string) error {
	// TODO(patrik): Add force flag

	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	conf := path.Join(src, "album.toml")

	data, err := os.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("%s: %w", conf, err)
	}

	var config AlbumConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("%s: %w", conf, err)
	}

	err = ExecuteConfig(config, mode, src, dst)
	if err != nil {
		return fmt.Errorf("%s: %w", conf, err)
	}

	return nil
}

const (
	ModeDwebble = "dwebble"
	ModeOpus    = "opus"
	ModeMp3     = "mp3"
	ModeMap     = "map"
)

func ExecuteConfig(config AlbumConfig, mode, src, dst string) error {
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

	wg := sync.WaitGroup{}

	plock := sync.Mutex{}

	// TODO(patrik): Check albumName for forward slashes and other illegal
	// filesystem characters
	for _, track := range config.Tracks {
		args := []string{}

		trackPath := path.Join(src, track.Filename)

		inputExt := path.Ext(trackPath)
		outputExt := ""

		switch mode {
		case ModeDwebble:
			outputExt = inputExt
			if inputExt == ".wav" {
				outputExt = ".flac"
			}
		case ModeMp3:
			// TODO(patrik): Add params
			outputExt = ".mp3"
		case ModeOpus:
			// TODO(patrik): Add params
			outputExt = ".opus"
		case ModeMap:
			outputExt = inputExt
		default:
			log.Fatal("Unknown mode:", mode)
		}

		copyMode := inputExt == outputExt

		args = append(args, "-i", trackPath, "-vn", "-map_metadata", "-1")

		artist := track.Artist
		if artist == "" {
			artist = config.Artist
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

		if len(track.Genres) > 0 {
			args = append(args, "-metadata", fmt.Sprintf("genre=%s", strings.Join(track.Genres, ",")))
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

		wg.Add(1)

		go func() {
			plock.Lock()
			fmt.Println("Processing:", trackPath)
			plock.Unlock()

			cmd := exec.Command("ffmpeg", args...)
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			plock.Lock()
			fmt.Println("Done Processing:", trackPath)
			plock.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
