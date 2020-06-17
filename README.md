# rename-images

No configuration image renaming.

This tool is intended to rename image files downloaded straight from camera
before any processing is done.

Output folder structure:

```bash
.
├── 2019
│   ├── 2019-06-10
│   │   ├── 20190610-042345_0101.arw
│   │   ├── 20190610-042345_0101.jpg
│   │   └── 20190610-042345_0203.arw
│   └── 2019-11-23
│       ├── 20191123-042345_0101.arw
│       ├── 20191123-042345_0101.arw.xmp
│       └── 20190610-042345_0102.arw
└── 2020
    └── 2020-03-08
        ├── 20200308-042345_101.arw
        ├── 20200308-042346_00101.nef
        └── 20200308-042347_02.cr2
```

This way even if you shoot with multiple cameras all pictures will be sorted
properly (if you synchronize time on all cameras).

Next you can change subfolders names to contain more useful info like
`2020-03-08_Alice_portrait`. Avoid spaces in names.

You need to remove old and empty folders by hand (it leave then just in case).

Warning: if destination file with same name already exists it will
be overwritten.

## Run it

Go to folder where you want to rename all image files and run program

```bash
$ rename-images
```

## Other info

Source of folder name structure idea:

https://www.photographyessentials.net/image-file-management-chronological-folder-names/

Additional good info:

https://composeclick.com/how-to-name-and-organize-your-photos/

https://www.scanyourentirelife.com/what-everybody-ought-know-when-naming-your-scanned-photos-part-1/
