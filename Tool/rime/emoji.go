package rime

import (
	"bufio"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

func Emoji() {
	// 控制台输出
	fmt.Println("检查 Emoji 和 main、sogou、ext 的差集：")
	defer printTimeCost(time.Now())

	// 打开文件
	file, err := os.Open(EmojiPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Emoji 的 set
	set := mapset.NewSet[string]()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		// Emoji 是单独的逻辑
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "\t") {
			log.Fatal("mapping.txt 含有 Tab：", line)
		}
		sp := strings.Split(line, " ")
		for _, hanzi := range sp[1:] {
			set.Add(hanzi)
		}
		continue
	}

	count := 0
	for _, item := range set.Difference(MainSet.Union(SogouSet).Union(ExtSet)).ToSlice() {
		// 去除英文字母
		if match, _ := regexp.MatchString(`[A-Za-z]+`, item); match {
			continue
		}
		// 去除一个字的
		if utf8.RuneCountInString(item) <= 1 {
			continue
		}
		count++
		fmt.Println(item)
	}
	fmt.Println("差集个数：", count)
}
