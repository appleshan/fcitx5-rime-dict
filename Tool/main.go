package main

import (
	"RimeTool/rime"
	"fmt"
	"strings"
)

func main() {
	rime.Emoji() // 检查 Emoji 与 main、sogou、ext 的差集，确保 Emoji 的词都在词库中
	rime.Sogou() // 下载搜狗流行词，加入到现有词库末尾，过滤、去重，打印新增词汇，搜狗转 Rime 词库：https://github.com/lewangdev/scel2txt
	doYouContinue()

	// 检查
	rime.Check(rime.HanziPath)
	rime.Check(rime.MainPath)
	rime.Check(rime.SogouPath)
	rime.Check(rime.ExtPath)
	rime.Check(rime.AVPath)
	rime.Check(rime.TencentPath)
	doYouContinue()

	// 去重、排序
	rime.Sort(rime.HanziPath)
	rime.Sort(rime.MainPath)
	rime.Sort(rime.SogouPath)
	rime.Sort(rime.ExtPath)
	rime.Sort(rime.AVPath)
	rime.Sort(rime.TencentPath)
}

func doYouContinue() {
	fmt.Println("是否继续：（\"ok\"）")
	var isOK string
	_, _ = fmt.Scanf("%s", &isOK)
	if strings.ToLower(isOK) != "ok" {
		return
	}
	fmt.Println("--------------------------------------------------")
}
