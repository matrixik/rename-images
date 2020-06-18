# sort-camera-pics

No configuration camera picture sorting.

This tool is intended to rename pictures files downloaded straight from camera
before any processing is done. It's searching for all pictures having filename
prefix like `DSC`, `_DSC`, `IMG` or similar (so already renamed pictures aren't
touch). Then extract `EXIF` creation date and time and use it for new name.
Input folder structure does not matter. All supported files will be moved
to new place.

`some_folder/placeA/_DSC7890.arw => 2020/2020-06-18/20200618-121314_7890.arw`

Sidecar files like `jpg`, `jpeg` or `xmp` also will be sorted (`xmp` will use
picture creation date and time).

Output folder structure:

```bash
.
├── 2019
│   ├── 2019-06-10
│   │   ├── 20190610-042345_0101.arw
│   │   ├── 20190610-042345_0101.jpg
│   │   └── 20190610-042345_0203.arw
│   └── 2019-11-23
│       ├── 20191123-234512_0101.arw
│       ├── 20191123-234512_0101.arw.xmp
│       └── 20190610-234512_0102.arw
└── 2020
    └── 2020-03-08
        ├── 20200308-123445_101.arw
        ├── 20200308-123446_00101.nef
        └── 20200308-123447_02.cr2
```

This way even if you shoot with multiple cameras all pictures will be sorted
properly (if you synchronize time on all cameras).

Next you can change subfolders names to contain more useful info like
`2020-03-08_Alice_portrait`. This way if you want to move/copy this folder
somewhere else it will always have unique name. Avoid spaces in names.

You need to remove old and empty folders by hand (it leave them just in case).

Warning: if destination file with same name already exists it will
be overwritten. Thanks to using time with seconds in new file names this should
almost never happen with different pictures.

## Usage

```bash
go get -u -v github.com/matrixik/sort-camera-pics
```

Go to folder where you want to rename all image files and run program

```bash
sort-camera-pics
```

## Other info

Source of folder name structure idea:

https://www.photographyessentials.net/image-file-management-chronological-folder-names/

Additional good info:

https://composeclick.com/how-to-name-and-organize-your-photos/

https://www.scanyourentirelife.com/what-everybody-ought-know-when-naming-your-scanned-photos-part-1/
