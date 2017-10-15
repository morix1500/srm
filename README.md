# srm
srm is a safe remove command.

## Description
srm creates a backup before deleting files/directories.  
In addtion, it has a function to restore deleted files.  

## Example
```
$ mkdir -p /tmp/test1/test2/test3
$ echo "foo" >> /tmp/test1/test2/foo.txt
$ echo "bar" >> /tmp/test1/test2/test3/bar.txt
$ tree /tmp/test1
 /tmp/test1/
└── test2
    ├── foo.txt
    └── test3
        └── bar.txt

2 directories, 2 files

# Delete "test2" directory
$ srm /tmp/test1/test2
$ tree /tmp/test1
/tmp/test1/

0 directories, 0 files

# Backups are stored in the following directories
$ ls ~/.srm
L3RtcC90ZXN0MS90ZXN0Mg==.tar.gz

# Restore "test2" directory
$ srm -r /tmp/test1/test2
$ tree /tmp/test1
/tmp/test1
└── test2
    ├── foo.txt
    └── test3
        └── bar.txt

2 directories, 2 files
```

## Usage
```
$ srm -h
Usage of srm:
  -l    Display a list of deleted files(directory) in the past.
  -list
        Display a list of deleted files(directory) in the past.
  -r    Restore deleted files(directory).
  -restore
        Restore deleted files(directory).
  -v    Display version.
  -version
        Display version.
```

### Option: list
Display deleted files/directories in a list.  
If restored, it will be deleted from the list.  

```
$ cd /tmp
$ touch foo.txt
$ touch bar.txt

$ srm foo.txt bar.txt
$ srm -l
/tmp/bar.txt
/tmp/foo.txt
```

### Option: restore
Restore deleted file/directory.  

```
$ cd /tmp
$ echo "foo" > foo.txt
$ srm foo.txt
$ ls -l foo.txt
ls: cannot access 'foo.txt': No such file or directory

$ srm -r foo.txt
$ cat foo.txt
foo
```

## Installtion
```
$ wget https://github.com/morix1500/srm/releases/download/v1.0.0/srm_linux_amd64 -O /usr/local/bin/srm
$ chmod u+x /usr/local/bin/srm
```

## License
Please see the [LICENSE](./LICENSE) file for details.  

## Author
Shota Omori(Morix)  
https://github.com/morix1500
