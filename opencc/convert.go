package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var basePath = getCurrentAbPath()
var mapTXT = filepath.Join(basePath, "mapping.txt")
var emojiTXT = filepath.Join(basePath, "emoji.txt")

type OrderedMap struct {
	keys []string
	m    map[string][]string
}

func main() {
	// fmt.Println(basePath)

	// 读取 map.txt
	mapTxt, err := os.Open(mapTXT)
	if err != nil {
		log.Fatalln(err)
	}
	defer mapTxt.Close()

	om := new(OrderedMap)
	om.m = make(map[string][]string)

	sc := bufio.NewScanner(mapTxt)
	lineNumber := 0
	for sc.Scan() {
		lineNumber++
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") {
			log.Fatalln("unexpected space", line)
		}

		arr := strings.Split(line, " ")
		if len(arr) < 2 {
			log.Fatalln("line error:", line, "\nline number:", lineNumber)
		}

		for _, zh := range arr[1:] {
			if !in(om.keys, zh) {
				om.keys = append(om.keys, zh)
			}
			om.m[zh] = append(om.m[zh], arr[0])
		}
	}

	// fmt.Println(om)

	// 写入 emoji.txt
	emojiTxt, err := os.OpenFile(emojiTXT, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer emojiTxt.Close()

	for _, key := range om.keys {
		emojis := strings.Join(om.m[key], " ")
		line := key + "\t" + key + " " + emojis + "\n"
		emojiTxt.WriteString(line)
	}
	emojiTxt.Sync()

	// 重新部署
	fmt.Println("done. 开始重新部署...")
	cmd := exec.Command("/Library/Input Methods/Squirrel.app/Contents/MacOS/Squirrel", "--reload")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// slice 是否包含某个元素
func in(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// https://zhuanlan.zhihu.com/p/363714760
// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
