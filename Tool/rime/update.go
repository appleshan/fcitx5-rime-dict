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

// UpdateSogou 1.下载搜狗流行词加入到现有词库末尾，2.过滤去重排序，3.打印新增词汇
func UpdateSogou() {
	// 控制台输出
	fmt.Println("搜狗流行词：")
	defer printTimeCost(time.Now())

	downloadAndWrite()       // 1.下载搜狗流行词加入到现有词库末尾
	checkAndWrite(SogouPath) // 2.过滤去重排序
	PrintNewWords(SogouPath) // 3.打印新增词汇

	// 弄完了删除临时用的文件，否则 VSCode 全局搜索词汇时会搜索到，影响体验
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

// UpdateWiki 和搜狗的逻辑是一样的，只是没有自动化下载更新
func UpdateWiki() {
	// 控制台输出
	fmt.Println("维基词库：")
	defer printTimeCost(time.Now())

	checkAndWrite(WikiPath)
	PrintNewWords(WikiPath)
}

// 弄好过滤词汇列表
func makeFilterList(dictPath string) mapset.Set[string] {
	filterList := mapset.NewSet[string]()

	file, err := os.Open(dictPath)
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

	return filterList
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

// checkAndWrite 过滤、排序、去重、重新写入
func checkAndWrite(dictPath string) {
	// 先制作好过滤词列表，在此列表的不加入进词库
	filterList := makeFilterList(dictPath)

	// 打开文件
	file, err := os.OpenFile(dictPath, os.O_RDWR, 0644)
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
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			log.Fatal("分割错误：", line)
		}
		text, code := parts[0], parts[1]
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
func PrintNewWords(dictPath string) {
	newSet := readAndSet(dictPath)
	newWords := mapset.NewSet[string]()
	if dictPath == SogouPath {
		newWords = newSet.Difference(SogouSet)
	}
	if dictPath == WikiPath {
		newWords = newSet.Difference(WikiSet)
	}
	fmt.Println("新增词汇：")

	// 没有注音
	// for _, word := range newWords.ToSlice() {
	// 	fmt.Println(word)
	// }

	// 把注音也打出来，方便直接校对
	file, err := os.Open(dictPath)
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
	if dictPath == SogouPath {
		SogouSet = newSet
	}
	if dictPath == WikiPath {
		WikiSet = newSet
	}
}