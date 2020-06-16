package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanoberholster/exiftool"
	"github.com/go-errors/errors"
)

var nameStarts = []string{
	"DSC",
	"CRW",
	"IMG",
}

func main() {
	fmt.Println(imagesInFolder("./"))
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
	clean := strings.TrimPrefix(strings.ToUpper(filename), "_")
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
				fileDate.Format("2006-01"),
				fileDate.Format("02"),
				fileDate.Format("20060102")+
					"_"+
					cleanName(filepath.Base(file))))
		renames[file] = newFilename
		xmpFilename := file + ".xmp"
		if _, err = os.Stat(xmpFilename); !os.IsNotExist(err) {
			renames[xmpFilename] = newFilename + ".xmp"
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

func safeRename(src, dest string) error {
	err := os.Link(src, dest)
	if err != nil {
		return err
	}

	return os.Remove(src)
}

func moveFiles(filesMap map[string]string) error {
	for src, dest := range filesMap {
		destPath := filepath.Dir(dest)
		fmt.Println(destPath)
		err := ensureDir(destPath)
		if err != nil {
			return err
		}

		err = safeRename(src, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
