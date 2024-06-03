package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flytam/filenamify"
	"github.com/nanoteck137/parasect"
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

var validTrackExts []string = []string{
	"wav",
	"m4a",
	"flac",
	"mp3",
	"opus",
	"ogg",
}

func IsValidTrackExt(ext string) bool {
	return parasect.IsValidExt(validTrackExts, ext)
}

var validCoverExts []string = []string{
	"png",
	"jpg",
	"jpeg",
}

func IsValidCoverExt(ext string) bool {
	return parasect.IsValidExt(validCoverExts, ext)
}


var lossyFormatExts = []string{
	"opus",
	"mp3",
}

func IsLossyFormatExt(ext string) bool {
	return parasect.IsValidExt(lossyFormatExts, ext)
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
		if d.Name()[0] != '.' && IsValidCoverExt(filepath.Ext(p)) {
			found = d.Name()
			return filepath.SkipAll
		}

		return nil
	})

	return found
}
