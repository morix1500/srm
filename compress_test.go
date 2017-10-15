package main

import(
	"fmt"
	"testing"
	"os/exec"
	"os"
	"path/filepath"
)

func TestCompress(t *testing.T) {
	base := "/tmp/srm-test/testcompress"
	source := filepath.Join(base, "source")
	sourceFile := filepath.Join(source, "test1", "test.txt")
	backupDir := filepath.Join(base, "backup")
	target := "test_compress"

	err := exec.Command("mkdir", "-p", filepath.Dir(sourceFile)).Run()
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("touch", sourceFile).Run()
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("mkdir", "-p", backupDir).Run()
	if err != nil {
		fmt.Println(err)
	}

	err = Compress(source, backupDir, target)
	if err != nil {
		t.Errorf("Output error", err)
	}

	want := filepath.Join(backupDir, target + ".tar.gz")
	if _, err := os.Stat(want); os.IsNotExist(err) {
		t.Errorf("faild create backup file. %s", want)
	}

	want = filepath.Join(backupDir, target + ".tar")
	if _, err := os.Stat(want); !os.IsNotExist(err) {
		t.Errorf("faild delete tar file. %s", want)
	}

	os.RemoveAll(base)
}

func TestUnCompress(t *testing.T) {
	base := "/tmp/srm-test/testuncompress"
	source := filepath.Join(base, "source")
	sourceFile := filepath.Join(source, "test1", "test.txt")
	backupDir := filepath.Join(base, "backup")
	target := "test_uncompress"

	err := exec.Command("mkdir", "-p", filepath.Dir(sourceFile)).Run()
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("touch", sourceFile).Run()
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("mkdir", "-p", backupDir).Run()
	if err != nil {
		fmt.Println(err)
	}

	err = Compress(source, backupDir, target)
	if err != nil {
		t.Errorf("Output error", err)
	}
	os.RemoveAll(source)

	err = UnCompress(filepath.Join(backupDir, target + ".tar.gz"), base)
	if err != nil {
		fmt.Println(err)
	}

	want := filepath.Join(backupDir, target + ".tar.gz")
	if _, err := os.Stat(want); !os.IsNotExist(err) {
		t.Errorf("faild delete backup file. %s", want)
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		t.Errorf("faild uncompress")
	}

	os.RemoveAll(base)
}
