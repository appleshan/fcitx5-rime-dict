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
	rime.UpdateWiki()  // 更新维基词库，手动录入后运行。
	doYouContinue()

	// 通用检查
	rime.Check(rime.HanziPath)
	rime.Check(rime.AVPath)
	rime.Check(rime.MainPath)
	rime.Check(rime.SogouPath)
	rime.Check(rime.ExtPath)
	rime.Check(rime.WikiPath)
	rime.Check(rime.TencentPath)
	doYouContinue()

	// 排序，顺便去重，包括自身重复和其他词库的重复
	rime.Sort(rime.HanziPath)
	rime.Sort(rime.AVPath)
	rime.Sort(rime.MainPath)
	rime.Sort(rime.SogouPath)   // 对 main 中已经有的，去重
	rime.Sort(rime.ExtPath)     // 对 main、sogou 中已经有的，去重
	rime.Sort(rime.WikiPath)    // 对 main、sogou、ext 中已经有的，去重
	rime.Sort(rime.TencentPath) // 对 main、sogou、ext、wiki 中已经有的，去重
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