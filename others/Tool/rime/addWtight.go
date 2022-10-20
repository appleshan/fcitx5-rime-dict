package rime

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// AddWeight 为 sogou、ext、tencent 这些词库中没有权重的词条自动加上权重
func AddWeight(dictPath string, flag int, weight int) {
	// 控制台输出
	fmt.Println("权重+1 " + path.Base(dictPath) + ":")
	defer printTimeCost(time.Now())

	// 读取文件到 lines 数组
	file, err := os.ReadFile(dictPath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")

	// 逐行遍历，加上 weight
	isMark := false
	for i, line := range lines {
		if !isMark {
			if line == mark {
				isMark = true
			}
			continue
		}

		// 过滤空行
		if line == "" {
			continue
		}

		// 修改行
		parts := strings.Split(line, "\t")
		if flag == 3 && len(parts) == 2 { // 汉字+注音+权重：sogou
			lines[i] = line + "\t" + strconv.Itoa(weight)
			// fmt.Printf("add: %d: %q\n", i, lines[i])
		}

		if flag == 4 && len(parts) == 1 { // 汉字+权重: ext tencent
			lines[i] = line + "\t" + strconv.Itoa(weight)
			// fmt.Printf("add: %d: %q\n", i, lines[i])
		}
	}

	// 重新写入
	resultString := strings.Join(lines, "\n")
	err = os.WriteFile(dictPath, []byte(resultString), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
