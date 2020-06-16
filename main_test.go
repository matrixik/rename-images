package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/otiai10/copy"
)

func Test_imagesInFolder(t *testing.T) {
	type args struct {
		folder string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			"Test recursive folders with images",
			args{folder: "assets"},
			[]string{
				// "assets/CRW_1446.CRW",
				"assets/_DSC3262.ARW",
				"assets/_DSC3262.JPG",
				"assets/test1/DSC_0976.NEF",
				"assets/test1/IMG_9526.CR2",
			},
			false},
		{
			"Test not existing folder",
			args{folder: "testowe"},
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := imagesInFolder(tt.args.folder)
			if (err != nil) != tt.wantErr {
				t.Errorf("imagesInFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("imagesInFolder() = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}

func Test_imageCreationDate(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			"Test ARW image data creation",
			args{path: "assets/_DSC3262.ARW"},
			time.Date(2020, time.June, 13, 17, 46, 29, 0, time.UTC),
			false,
		},
		{
			"Test JPG image data creation",
			args{path: "assets/_DSC3262.JPG"},
			time.Date(2020, time.June, 13, 17, 46, 29, 0, time.UTC),
			false,
		},
		{
			"Test NEF image data creation",
			args{path: "assets/test1/DSC_0976.NEF"},
			time.Date(2019, time.June, 29, 03, 11, 23, 0, time.UTC),
			false,
		},
		{
			"Test CR2 image data creation",
			args{path: "assets/test1/IMG_9526.CR2"},
			time.Date(2019, time.May, 29, 12, 32, 11, 0, time.UTC),
			false,
		},
		{
			"Test CRW image data creation",
			args{path: "assets/CRW_1446.CRW"},
			time.Date(2020, time.June, 01, 0, 0, 0, 0, time.UTC),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := imageCreationDate(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageCreationDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageCreationDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_proposeRename(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			"Test image rename",
			args{files: []string{
				// "assets/CRW_1446.CRW",
				"assets/_DSC3262.ARW",
				"assets/_DSC3262.JPG",
				"assets/test1/DSC_0976.NEF",
				"assets/test1/IMG_9526.CR2",
			}},
			map[string]string{
				// "assets/CRW_1446.CRW":           "2018-01/01/20180101_1446.crw",
				"assets/_DSC3262.ARW":           "2020-06/13/20200613_3262.arw",
				"assets/_DSC3262.JPG":           "2020-06/13/20200613_3262.jpg",
				"assets/test1/DSC_0976.NEF":     "2019-06/29/20190629_0976.nef",
				"assets/test1/DSC_0976.NEF.xmp": "2019-06/29/20190629_0976.nef.xmp",
				"assets/test1/IMG_9526.CR2":     "2019-05/29/20190529_9526.cr2",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proposeRename(tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("proposeRename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("proposeRename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_moveFiles(t *testing.T) {
	currentDir, _ := os.Getwd()
	defer os.Chdir(currentDir)

	dir, err := ioutil.TempDir("", "assets")
	if err != nil {
		t.Errorf("Creating temp dir error: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	err = copy.Copy("assets/", dir)
	if err != nil {
		t.Errorf("Copy dir error: %v", err)
	}

	_ = os.Chdir(dir)

	type args struct {
		filesMap map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test moving files",
			args{filesMap: map[string]string{
				// "CRW_1446.CRW":           "2018-01/01/20180101_1446.crw",
				"_DSC3262.ARW":           "2020-06/13/20200613_3262.arw",
				"_DSC3262.JPG":           "2020-06/13/20200613_3262.jpg",
				"test1/DSC_0976.NEF":     "2019-06/29/20190629_0976.nef",
				"test1/DSC_0976.NEF.xmp": "2019-06/29/20190629_0976.nef.xmp",
				"test1/IMG_9526.CR2":     "2019-05/29/20190529_9526.cr2",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := moveFiles(tt.args.filesMap); (err != nil) != tt.wantErr {
				t.Errorf("moveFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
