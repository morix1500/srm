package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unGzip(gzipPath string) (string, error) {
	reader, err := os.Open(gzipPath)

	if err != nil {
		return "", err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	tarfile := strings.Replace(archive.Name, ".gz", "", 1)
	writer, err := os.Create(tarfile)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return tarfile, err
}

func genGzip(tarPath string) error {
	reader, err := os.Open(tarPath)
	if err != nil {
		return err
	}

	filename := tarPath + ".gz"
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	if err != nil {
		return err
	}

	return os.Remove(tarPath)
}

func unTar(source, target string) error {
	defer func() {
		os.Remove(source)
	}()
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func genTar(source, outputDir, target string) (string, error) {
	filename := target + ".tar"
	output := filepath.Join(outputDir, filename)
	tarfile, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return "", err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)

			return err
		})

	return output, nil
}

func Compress(source, backupDir, target string) error {
	tarPath, err := genTar(source, backupDir, target)
	if err != nil {
		return err
	}
	err = genGzip(tarPath)
	if err != nil {
		return err
	}
	os.Remove(tarPath)

	return nil
}

func UnCompress(source, target string) error {
	tarPath, err := unGzip(source)
	if err != nil {
		return err
	}
	err = unTar(tarPath, target)
	if err != nil {
		return err
	}
	os.Remove(source)
	os.Remove(tarPath)

	return nil
}
