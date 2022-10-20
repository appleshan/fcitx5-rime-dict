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

	// flag: 1 只有汉字，2 汉字+注音，3 汉字+注音+权重，4 汉字+权重。
	// 为 sogou、ext、tencent 这些词库中没有权重的词条自动上权重
	rime.AddWeight(rime.SogouPath, 3, 1)
	rime.AddWeight(rime.ExtPath, 4, 1)
	rime.AddWeight(rime.TencentPath, 4, 1)

	// 通用检查
	rime.Check(rime.HanziPath, 3)
	rime.Check(rime.AVPath, 3)
	rime.Check(rime.MainPath, 3)
	rime.Check(rime.SogouPath, 3)
	rime.Check(rime.ExtPath, 4)
	rime.Check(rime.TencentPath, 4)

	areYouOK()

	// 排序，顺便去重，包括自身重复和其他词库的重复
	rime.Sort(rime.HanziPath, 3)
	rime.Sort(rime.AVPath, 3)
	rime.Sort(rime.MainPath, 3)
	rime.Sort(rime.SogouPath, 3)   // 对 main 中已经有的，去重
	rime.Sort(rime.ExtPath, 4)     // 对 main、sogou 中已经有的，去重
	rime.Sort(rime.TencentPath, 4) // 对 main、sogou、ext 中已经有的，去重
}

func areYouOK() {
	fmt.Println("是否继续：（\"ok\"）")
	var isOK string
	_, _ = fmt.Scanf("%s", &isOK)
	if strings.ToLower(isOK) != "ok" {
		os.Exit(123)
	}
	fmt.Println("--------------------------------------------------")
}
