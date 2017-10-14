package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

type CLI struct {
	outStream, errStream io.Writer
}

const (
	BackupDirName = ".srm"
)
const (
	ExitCodeOK = iota
	ExitCodeErr
)

func getHash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

func remove(source, backupDir string) error {
	_, err := os.Stat(source)
	if err != nil {
		return err
	}
	filename := getHash(source)

	// Create backup
	err = Compress(source, backupDir, filename)
	if err != nil {
		return err
	}

	return os.RemoveAll(source)
}

func restore(source, backupDir string) error {
	filename := getHash(source) + ".tar.gz"
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

func (c *CLI) Run(args []string) int {
	var isRestore bool
	flags := flag.NewFlagSet("srm", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&isRestore, "restore", false, "restore")
	flags.BoolVar(&isRestore, "r", false, "restore")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeErr
	}
	backupDir, err := createBackupDir()
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeErr
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
	os.Exit(cli.Run(os.Args))
}
