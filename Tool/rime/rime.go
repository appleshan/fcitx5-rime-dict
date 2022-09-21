package rime

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

const (
	mark        = "# +_+"
	HanziPath   = "/Users/dvel/Library/Rime/cn_dicts/8105.dict.yaml"
	MainPath    = "/Users/dvel/Library/Rime/cn_dicts/main.dict.yaml"
	SogouPath   = "/Users/dvel/Library/Rime/cn_dicts/sogou.dict.yaml"
	ExtPath     = "/Users/dvel/Library/Rime/cn_dicts/ext.dict.yaml"
	AVPath      = "/Users/dvel/Library/Rime/cn_dicts/av.dict.yaml"
	TencentPath = "/Users/dvel/Library/Rime/cn_dicts/tencent.dict.yaml"
	WikiPath    = "/Users/dvel/Library/Rime/cn_dicts/zhwiki.dict.yaml"
	EmojiPath   = "/Users/dvel/Library/Rime/opencc/mapping.txt"
)

var (
	// HanziSet  mapset.Set[string]
	MainSet   mapset.Set[string]
	SogouSet  mapset.Set[string]
	ExtSet    mapset.Set[string]
	TencenSet mapset.Set[string]
	WikiSet mapset.Set[string]
	initStart time.Time
)

// 一个词条的组成部分：汉字、编码、权重
type lemma struct {
	text   string
	code   string
	weight int
}

func init() {
	fmt.Println("rime.go init")
	// HanziSet = readAndSet(HanziPath)
	MainSet = readAndSet(MainPath)
	SogouSet = readAndSet(SogouPath)
	ExtSet = readAndSet(ExtPath)
	TencenSet = readAndSet(TencentPath)
	WikiSet = readAndSet(WikiPath)
}

// readAndSet 读取词库文件为 set
func readAndSet(dictPath string) mapset.Set[string] {
	set := mapset.NewSet[string]()

	file, err := os.Open(dictPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	isMark := false
	for sc.Scan() {
		line := sc.Text()
		if !isMark {
			if strings.Contains(line, mark) {
				isMark = true
			}
			continue
		}
		parts := strings.Split(line, "\t")
		set.Add(parts[0])
	}

	return set
}

// printTimeCost 打印耗时时间
func printTimeCost(start time.Time) {
	fmt.Printf("耗时：%.2fs\n", time.Since(start).Seconds())
	fmt.Println("--------------------------------------------------")
}