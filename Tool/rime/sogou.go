package rime

// 搜狗转 Rime 词库：https://github.com/lewangdev/scel2txt

import (
	"bufio"
	"fmt"
	"github.com/commander-cli/cmd"
	mapset "github.com/deckarep/golang-set/v2"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var filterMark = "# *_*"
var filterList mapset.Set[string] // 不需要的搜狗流行词词汇

// 初始化时弄好过滤词汇列表
func init() {
	fmt.Println("sogou.go init")
	filterList = mapset.NewSet[string]()

	file, err := os.Open(SogouPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	isFilterMark := false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line == mark {
			break
		}
		if !isFilterMark {
			if strings.Contains(line, filterMark) {
				isFilterMark = true
			}
			continue
		}
		// 判断一些可能出错的情况
		if !strings.HasPrefix(line, "# ") {
			log.Fatal("sogou 无效行：", line)
		}
		text := strings.TrimPrefix(line, "# ")
		if strings.ContainsAny(text, " \t") {
			log.Fatal("sogou 包含空字符，line:", line)
		}
		// 加入过滤词列表
		filterList.Add(text)
	}

	// 初始化结束
	defer printTimeCost(initStart)
}

// Sogou 下载搜狗流行词，加入到现有词库末尾，过滤、去重，打印新增词汇
func Sogou() {
	// 控制台输出
	fmt.Println("搜狗流行词：")
	defer printTimeCost(time.Now())

	downloadAndWrite()
	checkAndWrite()
	PrintNewWords()

	// 弄完了删除文件，否则 VSCode 搜索词汇时会搜索到
	err := os.Remove("./scel2txt/scel/sogou.scel")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("./scel2txt/out/luna_pinyin.sogou.dict.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("./scel2txt/out/sogou.txt")
	if err != nil {
		log.Fatal(err)
	}
}

// downloadAndWrite 下载搜狗流行词，转换为 Rime 格式，然后加入到现有词库的后面
func downloadAndWrite() {
	// 下载
	url := "https://pinyin.sogou.com/d/dict/download_cell.php?id=4&name=%E7%BD%91%E7%BB%9C%E6%B5%81%E8%A1%8C%E6%96%B0%E8%AF%8D%E3%80%90%E5%AE%98%E6%96%B9%E6%8E%A8%E8%8D%90%E3%80%91&f=detail"
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create("./scel2txt/scel/sogou.scel")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	// 用python进行转换
	c := cmd.NewCommand("python3 scel2txt.py", cmd.WithWorkingDir("./scel2txt"))
	err = c.Execute()
	if err != nil {
		fmt.Println(c.Stderr())
		log.Fatal(err)
	}
	fmt.Printf(c.Stdout())

	// 加入到现有词库的末尾
	sogouFile, err := os.OpenFile(SogouPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer sogouFile.Close()
	download, err := os.ReadFile("./scel2txt/out/sogou.txt")
	if err != nil {
		panic(err)
	}
	_, err = sogouFile.Write(download)
	if err != nil {
		panic(err)
	}
}

// checkAndWrite 检查 sogou，过滤、排序、去重
func checkAndWrite() {
	// 控制台输出
	fmt.Println("检查 sogou，过滤、排序、去重、重新写入：")
	defer printTimeCost(time.Now())

	// 打开文件
	file, err := os.OpenFile(SogouPath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 前缀内容和词库切片，前者原封不动直接写入，后者会检查、排序之类的
	prefixContents := make([]string, 0) // 前置内容切片
	contents := make([]lemma, 0)        // 词库切片

	isMark := false
	twoWordCount := 0
	wrongWordCount := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if !isMark {
			prefixContents = append(prefixContents, line)
			if line == mark {
				isMark = true
			}
			continue
		}
		sp := strings.Split(line, "\t")
		text, code := sp[0], sp[1]
		// 检查错误
		if len(sp) != 2 {
			log.Fatal("sogou 分割错误：", line)
		}
		// 计数
		if utf8.RuneCountInString(text) <= 2 {
			twoWordCount++
			continue
		}
		if filterList.Contains(text) {
			wrongWordCount++
			continue
		}
		contents = append(contents, lemma{text: text, code: code})
	}
	fmt.Println("两个字及以下的 count：", twoWordCount)
	fmt.Println("过滤掉的词汇 count：", wrongWordCount)

	// 排序
	sort.Slice(contents, func(i, j int) bool {
		if contents[i].code != contents[j].code {
			return contents[i].code < contents[j].code
		}
		if contents[i].text != contents[j].text {
			return contents[i].text < contents[j].text
		}
		return false
	})

	// 准备写入
	err = file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal()
	}
	// 写入前缀
	for _, content := range prefixContents {
		_, err := file.WriteString(content + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	// 写入词库，顺便去重
	duplicateCount := 0
	set := mapset.NewSet[string]()
	set = set.Union(MainSet)
	for _, content := range contents {
		if set.Contains(content.text) {
			duplicateCount++
			continue
		}
		set.Add(content.text)
		_, err := file.WriteString(content.text + "\t" + content.code + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("删除重复 count：", duplicateCount)
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

// PrintNewWords 打印新增词汇，注音也打出来，肉眼校对一边
func PrintNewWords() {
	newSogouSet := readAndSet(SogouPath)
	newWords := newSogouSet.Difference(SogouSet)
	fmt.Println("新增词汇：")
	// for _, word := range newWords.ToSlice() {
	// 	fmt.Println(word)
	// }
	file, err := os.Open(SogouPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	isMark := false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if !isMark {
			if line == mark {
				isMark = true
			}
			continue
		}
		if newWords.Contains(strings.Split(line, "\t")[0]) {
			fmt.Println(line)
		}
	}

	fmt.Println("count: ", newWords.Cardinality())

	// 更新全局的 set，方便后面检查
	SogouSet = newSogouSet
}
