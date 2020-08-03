package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	exts []string = []string{"png", "jpg"}
	list []string
)

// 判断是否多倍图
func isOnePicture(file string) bool {
	match, _ := regexp.MatchString(`\d+\.0x\\\w`, file)
	if (match) {
		return true
	}
	for _, ext := range(exts) {
		match, _ = regexp.MatchString(`\.` + ext + `$`, file)
		if (match) {
			return false
		}
	}
	return true
}

func GenerateAssetsClass(path string, savefile string) []string {
	var files []string

	println(path)
	filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if (!info.IsDir() && !isOnePicture(dir)) {
			files = append(files, dir)
		}
		return nil
	})
	for _, f := range files {
		println(f)
	}
	writeClass(savefile, files)
	return files
}

// 写入资源类
func writeClass(savefile string, files []string)  {
	var str []string = []string{"class KAssets {"}
	for _, file := range(files) {
		file = strings.ReplaceAll(file, "\\", "/")
		str = append(str, "static final String " + getMethodName(file) + " = \"" + file + "\";")
	}
	str = append(str, "}")

	dir := filepath.Dir(savefile)
	if (!PathExists(dir)) {
		err := os.MkdirAll(dir, 0755)
		if (err != nil) {
			fmt.Println(err)
			os.Exit(500)
		}
	}
	fs, err := os.OpenFile(savefile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fs.Close()
	n, err := fs.Write([]byte(strings.Join(str, "\r\n")))
	if err == nil && n < len(str) {
		fmt.Println(err)
		return
	}
}

func getMethodName(file string) string {
	suffix := strings.Join(exts, "|")
	reg := regexp.MustCompile(`/(.+?)\.(` + suffix + `)$`)
	if reg == nil {
		return ""
	}
	match := reg.FindStringSubmatch(file)
	if len(match) == 0 {
		return ""
	}
	reg = regexp.MustCompile(`/|_|-`)
	file = reg.ReplaceAllString(match[1], "/")
	reg = regexp.MustCompile(`/(\w)`)
	rep := reg.ReplaceAllStringFunc(file, func(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, "/", ""))
	})

	return rep
}

//
//func main() {
//	files := walk()
//	GenerateAssetsClass(files)
//	// println(getMethodName("assets/images/user/manor_wait.png"))
//}
