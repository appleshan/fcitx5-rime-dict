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
		// 过滤注释
		if strings.HasPrefix(line, "#") {
			continue
		}
		// 检查是否含有 Tab
		if strings.Contains(line, "\t") {
			fmt.Println("mapping.txt 含有 Tab：", line)
		}
		// 加入 set，顺便用 testSet 检查是否有重复
		sp := strings.Split(line, " ")
		testSet := mapset.NewSet[string]()
		for _, word := range sp[1:] {
			set.Add(word)
			if testSet.Contains(word) {
				fmt.Println("此行有重复项：", line)
			} else {
				testSet.Add(word)
			}
		}
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