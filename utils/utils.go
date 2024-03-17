package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

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
	"m4a",
	"flac",
	"mp3",
	"opus",
	"ogg",
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
