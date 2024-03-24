// Copyright (c) 2020, Dobrosław Żybort
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gosimple/hashdir"
	"github.com/otiai10/copy"
)

// fp returns path with system dependend separators.
func fp(paths ...string) string {
	return filepath.Join(paths...)
}

func Test_imagesInFolder(t *testing.T) {
	t.Parallel()

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
				fp("assets", "20200606", "20200606_2188.ARW"),
				fp("assets", "20200606", "DSC_02188.JPEG"),
				// fp("assets", "CRW_1446.CRW"),
				fp("assets", "PXL_20240113_192329323.NIGHT.jpg"),
				fp("assets", "_DSC3262.ARW"),
				fp("assets", "_DSC3262.JPG"),
				fp("assets", "test1", "DSC_0976.NEF"),
				fp("assets", "test1", "IMG_9526.CR2"),
			},
			false,
		},
		{
			"Test not existing folder",
			args{folder: "testowe"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := imagesInFolder(tt.args.folder)
			if (err != nil) != tt.wantErr {
				t.Errorf("imagesInFolder() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotFiles, tt.wantFiles) {
				t.Errorf("imagesInFolder() =\ndiff=\n%v",
					cmp.Diff(gotFiles, tt.wantFiles))
			}
		})
	}
}

func Test_imageCreationDate(t *testing.T) {
	t.Parallel()

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
			args{path: fp("assets", "_DSC3262.ARW")},
			time.Date(2020, time.June, 13, 17, 46, 29, 0, time.FixedZone("UTC+2", 2*60*60)),
			false,
		},
		{
			"Test JPG image data creation",
			args{path: fp("assets", "_DSC3262.JPG")},
			time.Date(2020, time.June, 13, 17, 46, 29, 0, time.FixedZone("UTC+2", 2*60*60)),
			false,
		},
		{
			"Test NEF image data creation",
			args{path: fp("assets", "test1", "DSC_0976.NEF")},
			time.Date(2019, time.June, 29, 0o3, 11, 23, 14000000, time.FixedZone("UTC-7", -7*60*60)),
			false,
		},
		{
			"Test CR2 image data creation",
			args{path: fp("assets", "test1", "IMG_9526.CR2")},
			time.Date(2019, time.May, 29, 12, 32, 11, 0, time.UTC),
			false,
		},
		{
			"Test CRW image data creation",
			args{path: fp("assets", "CRW_1446.CRW")},
			time.Date(2020, time.June, 0o1, 0, 0, 0, 0, time.UTC),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := imageCreationDate(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageCreationDate() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("imageCreationDate() =\ndiff=\n%v",
					cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_proposeRename(t *testing.T) {
	t.Parallel()

	type args struct {
		photoFile string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			"Test image rename 1",
			args{photoFile: fp("assets", "_DSC3262.ARW")},
			map[string]string{
				fp("assets", "_DSC3262.ARW"): fp("2020", "2020-06-13", "20200613-174629__dsc3262.arw"),
			},
			false,
		},
		{
			"Test image rename 2",
			args{photoFile: fp("assets", "_DSC3262.JPG")},
			map[string]string{
				fp("assets", "_DSC3262.JPG"): fp("2020", "2020-06-13", "20200613-174629__dsc3262.jpg"),
			},
			false,
		},
		{
			"Test image rename 3",
			args{photoFile: fp("assets", "20200606", "DSC_02188.JPEG")},
			map[string]string{
				fp("assets", "20200606", "DSC_02188.JPEG"): fp("2020", "2020-06-06", "20200606-080244_dsc_02188.jpg"),
			},
			false,
		},
		{
			"Test image rename 4",
			args{photoFile: fp("assets", "test1", "DSC_0976.NEF")},
			map[string]string{
				fp("assets", "test1", "DSC_0976.NEF"):     fp("2019", "2019-06-29", "20190629-031123_dsc_0976.nef"),
				fp("assets", "test1", "DSC_0976.NEF.xmp"): fp("2019", "2019-06-29", "20190629-031123_dsc_0976.nef.xmp"),
			},
			false,
		},
		{
			"Test image rename 5",
			args{photoFile: fp("assets", "test1", "IMG_9526.CR2")},
			map[string]string{
				fp("assets", "test1", "IMG_9526.CR2"): fp("2019", "2019-05-29", "20190529-123211_img_9526.cr2"),
			},
			false,
		},
		{
			"Test image rename 6",
			args{photoFile: fp("assets", "PXL_20240113_192329323.NIGHT.jpg")},
			map[string]string{
				fp("assets", "PXL_20240113_192329323.NIGHT.jpg"): fp("2024", "2024-01-13", "20240113-202329_pxl_20240113_192329323.night.jpg"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proposeRename(tt.args.photoFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("proposeRename() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("proposeRename() = \ndiff=\n%v",
					cmp.Diff(got, tt.want))
			}
		})
	}
}

// FIXME: not testing anything useful
// func Test_moveFiles(t *testing.T) {
// 	currentDir, _ := os.Getwd()
// 	defer func() {
// 		if chdirErr := os.Chdir(currentDir); chdirErr != nil {
// 			t.Errorf("Change dir error: %v", chdirErr)
// 		}
// 	}()

// 	dir, err := os.MkdirTemp("", "assets")
// 	if err != nil {
// 		t.Errorf("Creating temp dir error: %v", err)
// 	}
// 	defer os.RemoveAll(dir) // clean up

// 	err = copy.Copy("assets", dir)
// 	if err != nil {
// 		t.Errorf("Copy dir error: %v", err)
// 	}

// 	_ = os.Chdir(dir)

// 	type args struct {
// 		filesMap   map[string]string
// 		startingNr int
// 		nrWidth    int
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			"Test moving files",
// 			args{
// 				filesMap: map[string]string{
// 					"CRW_1446.CRW":                  fp("2018", "20180101_CRW_1446.crw"),
// 					"_DSC3262.ARW":                  fp("2020", "2020-06-13", "20200613-174629-_dsc3262.arw"),
// 					"_DSC3262.JPG":                  fp("2020", "2020-06-13", "20200613-174629-_dsc3262.jpg"),
// 					fp("test1", "DSC_0976.NEF"):     fp("2019", "2019-06-29", "20190629-031123-dsc_0976.nef"),
// 					fp("test1", "DSC_0976.NEF.xmp"): fp("2019", "2019-06-29", "20190629-031123-dsc_0976.nef.xmp"),
// 					fp("test1", "IMG_9526.CR2"):     fp("2019", "2019-05-29", "20190529-123211-img_9526.cr2"),
// 				},
// 				startingNr: 0,
// 				nrWidth:    1,
// 			},
// 			false,
// 		},
// 		{
// 			"Test moving files error",
// 			args{
// 				filesMap: map[string]string{
// 					"CRW_1446.CRW":              fp("2018", "20180101_CRW_1446.crw"),
// 					"_DSC3262.ARW":              fp("2020", "2020-06-13", "20200613-174629----_dsc3262.arw"),
// 					fp("test1", "IMG_9526.CR2"): fp("2019", "2019-05-29", "20190529-123211-img_9526.cr2"),
// 				},
// 				startingNr: 0,
// 				nrWidth:    1,
// 			},
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := moveFiles(
// 				tt.args.filesMap, tt.args.startingNr, tt.args.nrWidth)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("moveFiles() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_processImages(t *testing.T) {
	currentDir, _ := os.Getwd()
	defer func() {
		if chdirErr := os.Chdir(currentDir); chdirErr != nil {
			t.Errorf("Change dir error: %v", chdirErr)
		}
	}()

	dir, err := os.MkdirTemp("", "assets")
	if err != nil {
		t.Errorf("Creating temp dir error: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	err = copy.Copy("assets", dir)
	if err != nil {
		t.Errorf("Copy dir error: %v", err)
	}

	_ = os.Chdir(dir)

	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantWindows string
		wantErr     bool
	}{
		{
			"Test whole images processing",
			args{path: dir},
			// On all unix platforms
			"fb98508f35c383a3f5a3a70e7a3266a66d3db58b07fdd40e1d7e86427b68c02b",
			// On Windows (file path separator is different)
			"f63daf706deb41d9024558d58f00459090322a189d98f3070d618b680605c12d",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processImages(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("processImages() error = %v, wantErr %v",
					err, tt.wantErr)
			}

			hash, err := hashdir.Make("./", "sha256")
			if (err != nil) != tt.wantErr {
				t.Errorf("hashdir.Create() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}

			if !(runtime.GOOS == "windows") {
				if !cmp.Equal(hash, tt.want) {
					t.Errorf("proposeRename() =\ndiff=\n%v",
						cmp.Diff(hash, tt.want))
				}
			} else {
				// Different file path separator on Windows so hash also will
				// be different.
				if !cmp.Equal(hash, tt.wantWindows) {
					t.Errorf("proposeRename() =\ndiff=\n%v",
						cmp.Diff(hash, tt.wantWindows))
				}
			}
		})
	}
}

func Test_isFilenameDateTimePrefixed(t *testing.T) {
	type args struct {
		path string
	}
	tests := map[string]struct {
		args args
		want bool
	}{
		"Test 1": {args{path: fp("assets", "20200606", "20200606_2188.ARW")}, false},
		"Test 2": {args{path: fp("assets", "20200606", "20200606-214434_img234.arw")}, true},
		"Test 3": {args{path: fp("assets", "20200606-2144d4-img234.arw")}, false},
		"Test 4": {args{path: fp("assets", "200606-214434_img234.arw")}, false},
	}
	for name, tt := range tests {
		tt := tt
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := isFilenameDateTimePrefixed(tt.args.path); got != tt.want {
				t.Errorf("isFilenameDateTimePrefixed() = %v, want %v", got, tt.want)
			}
		})
	}
}
