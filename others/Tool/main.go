package main

import (
	"RimeTool/rime"
	"fmt"
	"os"
	"strings"
)

func main() {

	// 单项检查或更新
	rime.Emoji()       // 检查 Emoji 与 main、sogou、ext 的差集，确保 Emoji 的词都在词库中
	rime.UpdateSogou() // 自动下载搜狗流行词，加入到现有词库末尾，过滤、去重，打印新增词汇，搜狗转 Rime 词库：https://github.com/lewangdev/scel2txt
	// doYouContinue()

	// 通用检查
	rime.Check(rime.HanziPath)
	rime.Check(rime.AVPath)
	rime.Check(rime.MainPath)
	rime.Check(rime.SogouPath)
	rime.Check(rime.ExtPath)
	rime.Check(rime.TencentPath)
	doYouContinue()

	// 排序，顺便去重，包括自身重复和其他词库的重复
	// 列数，1 只有汉字，2 汉字+注音，3 汉字+注音+权重
	rime.Sort(rime.HanziPath, 3)
	rime.Sort(rime.AVPath, 3)
	rime.Sort(rime.MainPath, 3)
	rime.Sort(rime.SogouPath, 2)   // 对 main 中已经有的，去重
	rime.Sort(rime.ExtPath, 1)     // 对 main、sogou 中已经有的，去重
	rime.Sort(rime.TencentPath, 1) // 对 main、sogou、ext 中已经有的，去重
}

func doYouContinue() {
	fmt.Println("是否继续：（\"ok\"）")
	var isOK string
	_, _ = fmt.Scanf("%s", &isOK)
	if strings.ToLower(isOK) != "ok" {
		os.Exit(123)
	}
	fmt.Println("--------------------------------------------------")
}
