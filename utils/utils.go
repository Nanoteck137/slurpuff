package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/flytam/filenamify"
)

func Copy(src, dst string) (int64, error) {
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

func RunFFprobe(args ...string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)
	if true {
		cmd.Stderr = os.Stderr
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

type ProbeResult struct {
	Artist      string
	AlbumArtist string
	Title       string
	Album       string
	Track       int
	Disc        int
}

type FileResult struct {
	Path   string
	Number int
	Name   string

	Probe ProbeResult
}

type probeFormat struct {
	BitRate string `json:"bit_rate"`
	Tags    struct {
		Album       string `json:"album"`
		AlbumArtist string `json:"album_artist"`
		Artist      string `json:"artist"`
		Disc        string `json:"disc"`
		Title       string `json:"title"`
		Track       string `json:"track"`
		// "encoder": "Lavf58.29.100",
	} `json:"tags"`

	// "filename": "/Volumes/media/music/Various Artists/Cyberpunk 2077/cd1/19 - P.T. Adamczyk - Rite Of Passage.mp3",
	// "nb_streams": 2,
	// "nb_programs": 0,
	// "format_name": "mp3",
	// "format_long_name": "MP2/3 (MPEG audio layer 2/3)",
	// "start_time": "0.025056",
	// "duration": "334.915918",
	// "size": "13898147",
	// "probe_score": 51,
}

type probeStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`

	// Video
	Width  int `json:"width"`
	Height int `json:"height"`

	Disposition struct {
		AttachedPic int `json:"attached_pic"`
		// "default": 0,
		// "dub": 0,
		// "original": 0,
		// "comment": 0,
		// "lyrics": 0,
		// "karaoke": 0,
		// "forced": 0,
		// "hearing_impaired": 0,
		// "visual_impaired": 0,
		// "clean_effects": 0,
		// "timed_thumbnails": 0,
		// "captions": 0,
		// "descriptions": 0,
		// "metadata": 0,
		// "dependent": 0,
		// "still_image": 0
	} `json:"disposition"`

	Tags struct {
		Comment string `json:"comment"`
		// "comment": "Cover (front)"
	} `json:"tags"`

	// "codec_long_name": "PNG (Portable Network Graphics) image",
	// "codec_tag_string": "[0][0][0][0]",
	// "codec_tag": "0x0000",
	// "coded_width": 512,
	// "coded_height": 512,
	// "closed_captions": 0,
	// "film_grain": 0,
	// "has_b_frames": 0,
	// "pix_fmt": "rgba",
	// "level": -99,
	// "color_range": "pc",
	// "refs": 1,
	// "r_frame_rate": "90000/1",
	// "avg_frame_rate": "0/0",
	// "time_base": "1/90000",
	// "start_pts": 2255,
	// "start_time": "0.025056",
	// "duration_ts": 30142433,
	// "duration": "334.915922",
}

type probe struct {
	Streams []probeStream `json:"streams"`
	Format  probeFormat   `json:"format"`
}

func getNumberFromFormatString(s string) int {
	if strings.Contains(s, "/") {
		s = strings.Split(s, "/")[0]
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}

	return num
}

// TODO(patrik): Update to not include file extentions
var test1 = regexp.MustCompile(`(^\d+)[-\s]*(.+)\.`)
var test2 = regexp.MustCompile(`track(\d+).+`)

// TODO(patrik): Fix this function
func CheckFile(filepath string) (FileResult, error) {
	// ffprobe -v quiet -print_format json -show_format -show_streams input
	data, err := RunFFprobe("-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filepath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return FileResult{}, err
	}

	// fmt.Printf("string(data): %v\n", string(data))

	var probe probe
	err = json.Unmarshal(data, &probe)
	if err != nil {
		return FileResult{}, err
	}

	// fmt.Printf("probe: %+v\n", probe)
	// probe.Format.Tags.Track

	track := getNumberFromFormatString(probe.Format.Tags.Track)
	disc := getNumberFromFormatString(probe.Format.Tags.Disc)

	probeResult := ProbeResult{
		Artist:      probe.Format.Tags.Artist,
		AlbumArtist: probe.Format.Tags.AlbumArtist,
		Title:       probe.Format.Tags.Title,
		Album:       probe.Format.Tags.Album,
		Track:       track,
		Disc:        disc,
	}

	name := path.Base(filepath)
	res := test1.FindStringSubmatch(name)
	if res == nil {
		res := test2.FindStringSubmatch(name)
		if res == nil {
			return FileResult{}, fmt.Errorf("No result")
		}

		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		return FileResult{
			Path:   filepath,
			Number: num,
			Name:   "",
			Probe:  probeResult,
		}, nil
	} else {
		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		name := string(res[2])
		return FileResult{
			Path:   filepath,
			Number: num,
			Name:   name,
			Probe:  probeResult,
		}, nil
	}
}

func IsValidExt(exts []string, ext string) bool {
	if len(ext) == 0 {
		return false
	}

	if ext[0] == '.' {
		ext = ext[1:]
	}

	for _, valid := range exts {
		if valid == ext {
			return true
		}
	}

	return false
}

var validTrackExts []string = []string{
	"wav",
	"flac",
	"mp3",
}

func IsValidTrackExt(ext string) bool {
	return IsValidExt(validTrackExts, ext)
}

var validCoverExts []string = []string{
	"png",
	"jpg",
	"jpeg",
}

func IsValidCoverExt(ext string) bool {
	return IsValidExt(validCoverExts, ext)
}

func SafeName(name string) (string, error) {
	replacementSpace := func(options *filenamify.Options) { 
		options.Replacement = "" 
	}

	return filenamify.FilenamifyV2(name, replacementSpace)
}

func FindFirstValidImage(dir string) string {
	found := ""

	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if IsValidCoverExt(filepath.Ext(p)) {
			found = d.Name()
			return filepath.SkipAll
		}

		return nil
	})

	return found
}
