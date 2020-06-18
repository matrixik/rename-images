// No configuration camera picture sorting.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanoberholster/exiftool"
	"github.com/go-errors/errors"
)

var nameStarts = []string{
	// Names could start with `DSC_` but also with `_DSC` (like Sony AdobeRGB),
	// it's handled in the code.
	// Source: https://en.wikipedia.org/wiki/MediaWiki:Filename-prefix-blacklist
	"CIMG", // Casio
	"DSC",  // Nikon, Sony
	"DSCF", // Fuji
	"DSCN", // Nikon
	"DUW",  // some mobile phones
	"IMAG", // Many companies
	"IMG",  // generic
	"JD",   // Jenoptik
	"KIF",  // Kyocera
	"MGP",  // Pentax
	"S700", // Samsung
	"PICT", // misc.
}

var emptyFolders []string

func main() {
	err := processImages("./")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println("Success")
}

func processImages(path string) error {
	files, err := imagesInFolder(path)
	if err != nil {
		return err
	}

	filesMap, err := proposeRename(files)
	if err != nil {
		return err
	}

	return moveFiles(filesMap)
}

func imagesInFolder(root string) (files []string, err error) {
	if _, err = os.Stat(root); os.IsNotExist(err) {
		return nil, err
	}
	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			// Select only files that need to have changed name
			if !info.IsDir() &&
				isSupportedFile(path) &&
				hasDefaultName(path) {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return
}

func isSupportedFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case
		// Sony
		".arw",
		// Nikon
		".nef",
		// Canon
		// ".crw",
		".cr2",
		// Adobe
		".dng",
		// Jpegs
		".jpg",
		".jpeg":
		return true
	}
	return false
}

func hasDefaultName(path string) bool {
	cleanFilename :=
		strings.TrimPrefix(filepath.Base(strings.ToUpper(path)), "_")
	for _, pref := range nameStarts {
		if strings.HasPrefix(cleanFilename, pref) {
			return true
		}
	}
	return false
}

func cleanName(filename string) string {
	clean := strings.ReplaceAll(
		strings.TrimPrefix(strings.ToUpper(filename), "_"), "JPEG", "JPG")
	for _, pref := range nameStarts {
		if strings.HasPrefix(clean, pref) {
			return strings.TrimPrefix(strings.TrimPrefix(clean, pref), "_")
		}
	}
	return filename
}

func imageCreationDate(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Date(2020, time.June, 01, 0, 0, 0, 0, time.UTC), err
	}
	defer f.Close()

	eh, err := exiftool.SearchExifHeader(f)
	if err != nil {
		return time.Date(2020, time.June, 01, 0, 0, 0, 0, time.UTC),
			errors.Errorf("File: %v, error: %v", path, err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return time.Date(2020, time.June, 01, 0, 0, 0, 0, time.UTC), err
	}

	buf, _ := ioutil.ReadAll(f)
	e, err := eh.ParseExif(bytes.NewReader(buf))
	if err != nil {
		return time.Date(2020, time.June, 01, 0, 0, 0, 0, time.UTC), err
	}

	return e.DateTime()
}

func proposeRename(files []string) (map[string]string, error) {
	renames := make(map[string]string)
	for _, file := range files {

		fileDate, err := imageCreationDate(file)
		if err != nil {
			return nil, err
		}

		newFilename := strings.ToLower(
			filepath.Join(
				fileDate.Format("2006"),
				fileDate.Format("2006-01-02"),
				fileDate.Format("20060102-150405")+
					"_"+
					cleanName(filepath.Base(file))))
		renames[file] = newFilename

		// Copy also sidecar files
		xmpFilenameSmall := file + ".xmp"
		if _, err = os.Stat(xmpFilenameSmall); !os.IsNotExist(err) {
			renames[xmpFilenameSmall] = newFilename + ".xmp"
		}
		xmpFilenameBig := file + ".XMP"
		if _, err = os.Stat(xmpFilenameBig); !os.IsNotExist(err) {
			renames[xmpFilenameBig] = newFilename + ".xmp"
		}
	}
	return renames, nil
}

func ensureDir(folder string) error {
	err := os.MkdirAll(folder, os.ModePerm)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func isEmpty(dir string) bool {
	// Source: https://stackoverflow.com/a/30708914/1722542

	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	return err == io.EOF
}

func moveFiles(filesMap map[string]string) error {
	for src, dest := range filesMap {
		destPath := filepath.Dir(dest)
		err := ensureDir(destPath)
		if err != nil {
			return err
		}

		log.Println("Move file", src, "to", dest)
		err = os.Rename(src, dest)
		if err != nil {
			return err
		}

		srcPath := filepath.Dir(src)
		if isEmpty(srcPath) {
			emptyFolders = append(emptyFolders, srcPath)
		}
	}
	log.Println("Moved", len(filesMap), "files.")

	if len(emptyFolders) > 0 {
		log.Println("Left", len(emptyFolders), "empty folder(s):",
			emptyFolders)
	}
	return nil
}
