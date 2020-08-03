package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

var (
	sep string = string(os.PathSeparator)
)

func isOnexPicture(str string) bool {
	return strings.Index(str, "@") > -1
}

func getFiles(path string) []string {
	var files []string

	filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if (!info.IsDir()) {
			files = append(files, dir)
		}
		return nil
	})
	return files
}

func replaceFile(path string, file string) string {
	reg := regexp.MustCompile(`[\\|/](\w+)@(\dx)\.(\w+)$`)
	str := reg.ReplaceAllString(file, sep + "$2" + sep + "$1.$3")
	reg = regexp.MustCompile(`(\d)x`)
	str = reg.ReplaceAllString(str, "$1.0x")
	return strings.Replace(str, path, "", 1)
}

func movefile(oldpath, newpath string) error { //跨卷移动
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return syscall.MoveFile(from, to)//windows API

}

func Move(path string, to string)  {
	path = strings.ReplaceAll(path, "/", sep)
	fs := getFiles(path)
	for _, f := range(fs) {
		tofile := filepath.Join(to, replaceFile(path, f))
		dir := filepath.Dir(tofile)
		if (!PathExists(dir)) {
			err := os.MkdirAll(dir, 0755)
			if (err != nil) {
				fmt.Println(err)
				os.Exit(500)
			}
		}
		err := movefile(f, tofile)
		if (err != nil) {
			fmt.Println(err)
			os.Exit(500)
		}
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
