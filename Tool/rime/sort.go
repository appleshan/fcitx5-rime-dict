package rime

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Sort(dictPath string) {
	// 控制台输出
	fmt.Println("开始排序 ", path.Base(dictPath), "：")
	defer printTimeCost(time.Now())

	// 是否有任何变动
	oldSha1 := getSha1(dictPath)
	defer func(oldSha1 string) {
		newSha1 := getSha1(dictPath)
		if newSha1 != oldSha1 {
			fmt.Println("sorted")
			updateVersion(dictPath)
		}
	}(oldSha1)

	// 打开文件
	file, err := os.OpenFile(dictPath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 前缀内容和词库切片，前者原封不动直接写入，后者做排序后再写入
	prefixContents := make([]string, 0) // 前置内容切片
	contents := make([]lemma, 0)        // 词库切片
	selfSet := mapset.NewSet[string]()  // 用作去除和自身重复的
	isMark := false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		// mark 之前的，写入 prefixContents
		if !isMark {
			prefixContents = append(prefixContents, line)
			if line == mark {
				isMark = true
			}
			continue
		}

		// 分割为 词汇text 编码code 权重weight
		parts := strings.Split(line, "\t")
		text, code, weight := parts[0], "", ""

		// 将 main 中注释了但没删除的词汇权重调为 0
		if dictPath == MainPath && strings.HasPrefix(line, "# ") {
			line = parts[0] + "\t" + parts[1] + "\t" + "0"
		}

		// mark 之后的，写入到 contents
		// 自身重复的直接排除，不重复的写入
		switch len(parts) {
		case 1: // ext tencent 是一列
			if selfSet.Contains(text) {
				fmt.Println("重复：", line)
				break
			}
			selfSet.Add(text)
			contents = append(contents, lemma{text: text})
		case 2: // sogou 是两列
			text, code = parts[0], parts[1]
			if selfSet.Contains(text + code) {
				fmt.Println("重复：", line)
				break
			}
			selfSet.Add(text + code)
			contents = append(contents, lemma{text: text, code: code})
		case 3: // 字表 main av 是三列
			text, code, weight = parts[0], parts[1], parts[2]
			if selfSet.Contains(text + code) {
				fmt.Println("重复：", line)
				break
			}
			selfSet.Add(text + code)
			weight, _ := strconv.Atoi(weight)
			contents = append(contents, lemma{text: text, code: code, weight: weight})
		default:
			log.Fatal("分割错误：", line)
		}
	}

	// 排序：拼音升序、权重降序、最后直接按Unicode编码排序
	sort.Slice(contents, func(i, j int) bool {
		if contents[i].code != contents[j].code {
			return contents[i].code < contents[j].code
		}
		if contents[i].weight != contents[j].weight {
			return contents[i].weight > contents[j].weight
		}
		if contents[i].text != contents[j].text {
			return contents[i].text < contents[j].text
		}
		return false
	})

	// 下面开始写入，顺便从其他词库中去重
	// 准备写入
	err = file.Truncate(0)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalln(err)
	}

	// 写入前缀
	for _, content := range prefixContents {
		_, err := file.WriteString(content + "\n")
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 字表、main、av，直接写入，不需要从其他词库去重
	if dictPath == HanziPath || dictPath == MainPath || dictPath == AVPath {
		for _, line := range contents {
			_, err := file.WriteString(line.text + "\t" + line.code + "\t" + strconv.Itoa(line.weight) + "\n")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	// 其他词库需要从一个或多个词库中去重
	if dictPath == SogouPath || dictPath == ExtPath || dictPath == TencentPath {
		var intersect mapset.Set[string]
		switch dictPath {
		case SogouPath:
			// sogou 不和 main 有重复
			intersect = SogouSet.Intersect(MainSet)
		case ExtPath:
			// ext 不和 mian+sogou 有重复
			intersect = ExtSet.Intersect(MainSet.Union(SogouSet))
		case TencentPath:
			// tencent 不和 main+sogou+ext 有重复
			intersect = TencenSet.Intersect(MainSet.Union(SogouSet).Union(ExtSet))
		default:
			log.Fatal("？？？")
		}

		count := 0 // 重复个数
		for _, line := range contents {
			if intersect.Contains(line.text) {
				count++
				fmt.Printf("%s 重复于其他词库：%s\n", strings.Split(path.Base(dictPath), ".")[0], line.text)
				continue
			}
			str := ""
			if dictPath == ExtPath || dictPath == TencentPath {
				str = line.text + "\n"
			}
			if dictPath == SogouPath {
				str = line.text + "\t" + line.code + "\n"
			}
			_, err := file.WriteString(str)
			if err != nil {
				log.Fatal(err)
			}
		}
		if count != 0 {
			fmt.Println("重复个数：", count)
		}
	}

	// 同步
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func getSha1(dictPath string) string {
	f, err := os.Open(dictPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sha1Handle := sha1.New()
	if _, err := io.Copy(sha1Handle, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(sha1Handle.Sum(nil))
}

// 如果词库发生了变动，顺便把日期也给改了。
func updateVersion(dictPath string) {
	// 打开文件
	file, err := os.OpenFile(dictPath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 修改并读取到 arr
	arr := make([]string, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "version:") {
			s := fmt.Sprintf("version: \"%s\"", time.Now().Format("2006-01-02"))
			arr = append(arr, s)
		} else {
			arr = append(arr, line)
		}
	}

	// 重新写入
	err = file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal()
	}
	for _, line := range arr {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
