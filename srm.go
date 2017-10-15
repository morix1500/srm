package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CLI -- command line interface
type CLI struct {
	outStream, errStream io.Writer
}

const (
	// BackupDirName -- backup directory name
	BackupDirName = ".srm"
)
const (
	// ExitCodeOK -- success code
	ExitCodeOK = iota
	// ExitCodeErr -- error code
	ExitCodeErr
)

func encodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func decodeBase64(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func remove(source, backupDir string) error {
	_, err := os.Stat(source)
	if err != nil {
		return err
	}
	filename := encodeBase64(source)

	// Create backup
	err = Compress(source, backupDir, filename)
	if err != nil {
		return err
	}

	return os.RemoveAll(source)
}

func restore(source, backupDir string) error {
	filename := encodeBase64(source) + ".tar.gz"
	backupFile := filepath.Join(backupDir, filename)

	_, err := os.Stat(backupFile)
	if err != nil {
		return errors.New(source + ": not found backup file")
	}
	output, _ := filepath.Split(source)
	err = UnCompress(backupFile, output)
	if err != nil {
		return err
	}

	return nil
}

func list(backupDir string) ([]string, error) {
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return []string{}, err
	}
	arr := make([]string, len(files))
	for i := 0; i < len(files); i++ {
		name := strings.Split(files[i].Name(), ".")[0]
		name, err := decodeBase64(name)
		if err != nil {
			return []string{}, err
		}
		arr[i] = name
	}
	sort.Strings(arr)
	return arr, nil
}

func createBackupDir() (string, error) {
	home := os.Getenv("HOME")
	backupDir := filepath.Join(home, BackupDirName)
	_, err := os.Stat(backupDir)
	if err == nil {
		return backupDir, nil
	}

	err = os.Mkdir(backupDir, 0700)
	if err != nil {
		return "", errors.Wrap(err, "Failed make backup directory")
	}

	return backupDir, nil
}

func (c *CLI) run(args []string) int {
	var isRestore bool
	var isList bool
	var isVersion bool

	flags := flag.NewFlagSet("srm", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&isRestore, "restore", false, "Restore deleted files(directory).")
	flags.BoolVar(&isRestore, "r", false, "Restore deleted files(directory).")
	flags.BoolVar(&isList, "list", false, "Display a list of deleted files(directory) in the past.")
	flags.BoolVar(&isList, "l", false, "Display a list of deleted files(directory) in the past.")
	flags.BoolVar(&isVersion, "v", false, "Display version.")
	flags.BoolVar(&isVersion, "version", false, "Display version.")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeErr
	}
	if isVersion {
		fmt.Fprintln(c.outStream, Version)
		return ExitCodeOK
	}

	backupDir, err := createBackupDir()
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeErr
	}
	if isList {
		files, err := list(backupDir)
		if err != nil {
			fmt.Fprintln(c.errStream, err)
			return ExitCodeErr
		}
		for i := 0; i < len(files); i++ {
			fmt.Fprintln(c.outStream, files[i])
		}
		return ExitCodeOK
	}

	targetFiles := flags.Args()
	for _, file := range targetFiles {
		abs, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintln(c.errStream, err)
			return ExitCodeErr
		}
		if isRestore {
			if err := restore(abs, backupDir); err != nil {
				fmt.Fprintln(c.errStream, err)
				return ExitCodeErr
			}
		} else {
			if err := remove(abs, backupDir); err != nil {
				fmt.Fprintln(c.errStream, err)
				return ExitCodeErr
			}
		}
	}

	return ExitCodeOK
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.run(os.Args))
}
