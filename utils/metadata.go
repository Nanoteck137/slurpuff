package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func RunFFprobe(args ...string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)
	cmd.Stderr = os.Stderr

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

	Duration int
	Tags     map[string]string
}

type probeFormat struct {
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Duration string            `json:"duration"`

	// "filename": "/Volumes/media/music/Various Artists/Cyberpunk 2077/cd1/19 - P.T. Adamczyk - Rite Of Passage.mp3",
	// "nb_streams": 2,
	// "nb_programs": 0,
	// "format_name": "mp3",
	// "format_long_name": "MP2/3 (MPEG audio layer 2/3)",
	// "start_time": "0.025056",
	// "size": "13898147",
	// "probe_score": 51,
}

type probeStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`

	Duration string `json:"duration"`

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

	Tags     map[string]string `json:"tags"`

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

func convertMapKeysToLowercase(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

var test1 = regexp.MustCompile(`(^\d+)[-\s]*(.+)\.`)
var test2 = regexp.MustCompile(`track(\d+).+`)

type Info struct {
	Tags map[string]string
	Duration int
}

func GetInfo(filepath string) (Info, error) {
	// ffprobe -v quiet -print_format json -show_format -show_streams input
	data, err := RunFFprobe("-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filepath)
	if err != nil {
		return Info{}, err
	}

	var probe probe
	err = json.Unmarshal(data, &probe)
	if err != nil {
		return Info{}, err
	}

	var tags map[string]string

	duration := 0
	for _, s := range probe.Streams {
		if s.CodecType == "audio" {
			dur, err := strconv.ParseFloat(s.Duration, 32)
			if err != nil {
				return Info{}, err
			}

			duration = int(dur)
			tags = convertMapKeysToLowercase(s.Tags)
		}
	}

	return Info{
		Tags:     tags,
		Duration: duration,
	}, nil
}

func CheckFile(filepath string) (FileResult, error) {
	info, err := GetInfo(filepath)
	if err != nil {
		return FileResult{}, err
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
			Path:     filepath,
			Number:   num,
			Name:     "",
			Duration: info.Duration,
			Tags:     info.Tags,
		}, nil
	} else {
		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		name := string(res[2])
		return FileResult{
			Path:     filepath,
			Number:   num,
			Name:     name,
			Duration: info.Duration,
			Tags:     info.Tags,
		}, nil
	}
}
