package rime

import (
	"bufio"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var specialWords = mapset.NewSet[string]()   // 特殊词汇列表，不进行任何检查
var polyphoneWords = mapset.NewSet[string]() // 需要注音的多音字列表
var wrongWords = mapset.NewSet[string]()     // 异形词列表
var filterWords = mapset.NewSet[string]()    // 与异形词列表同时使用，过滤掉，一般是包含异形词但不是异形词的，比如「作爱」是异形词，但「叫作爱」是正确的。

// 初始化特殊词汇列表、多音字列表、异形词列表
func init() {
	fmt.Println("init...")
	fmt.Println("check.go init")
	// 特殊词汇列表
	specialWords.Add("狄尔斯–阿尔德反应")
	specialWords.Add("特里斯坦–达库尼亚")
	specialWords.Add("特里斯坦–达库尼亚群岛")
	specialWords.Add("茱莉亚·路易斯-德瑞弗斯")
	specialWords.Add("科科斯（基林）群岛")
	specialWords.Add("刚果（金）")
	specialWords.Add("刚果（布）")
	specialWords.Add("赛博朋克：边缘行者")
	specialWords.Add("赛博朋克：边缘跑手")
	specialWords.Add("赛博朋克：命运之轮")

	// 需要注音的多音字列表
	polyphoneFile, err := os.Open("rime/多音字.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer polyphoneFile.Close()
	sc := bufio.NewScanner(polyphoneFile)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		polyphoneWords.Add(line)
	}

	// 异形词的两个列表
	file, err := os.Open("rime/异形词.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sc = bufio.NewScanner(file)
	isMark := false
	for sc.Scan() {
		line := sc.Text()

		if strings.Contains(line, "# -_-") {
			isMark = true
			continue
		}

		if !isMark {
			wrongWords.Add(line)
		} else {
			filterWords.Add(line)
		}
	}
}

// Check 对传入的词库文件进行检查
func Check(dictPath string) {
	// 控制台输出
	fmt.Println("检查 " + path.Base(dictPath) + ":")
	defer printTimeCost(time.Now())

	// 打开文件
	file, err := os.Open(dictPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 开始检查！
	lineNumber := 0
	isMark := false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lineNumber++
		line := sc.Text()
		if !isMark {
			if line == mark {
				isMark = true
			}
			continue
		}

		// 忽略注释，main 里有很多被注释了的词汇，暂时没有删除
		if strings.HasPrefix(line, "#") {
			continue
		}

		// 注释以"#"开头，但不是以"# "开头（没有空格）（强迫症晚期）
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# ") {
			fmt.Println("has # but not #␣", line)
		}

		// 是否有空行
		if strings.TrimSpace(line) == "" {
			fmt.Println("空行，行号：", lineNumber)
		}

		// 开头结尾是否有空字符
		if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") ||
			strings.HasPrefix(line, "\t") || strings.HasSuffix(line, "\t") {
			fmt.Printf("开头或结尾有空格：%q\n", line)
		}

		// +---------------------------------------------------------------
		// | 通用检查：分割为：词汇text, 编码code, 权重weight
		// +---------------------------------------------------------------
		parts := strings.Split(line, "\t")
		var text, code, weight string
		switch len(parts) {
		case 1: // ext tencent 是一列
			text = parts[0]
		case 2: // sogou 是两列
			text, code = parts[0], parts[1]
		case 3: // 字表 main av 是三列
			text, code, weight = parts[0], parts[1], parts[2]
		default:
			log.Fatal("分割错误：", line)
		}

		// 过滤特殊词条
		if specialWords.Contains(text) {
			continue
		}

		// 检查：text 和 weight 不应该含有空格
		if strings.Contains(text, " ") || strings.Contains(weight, " ") {
			fmt.Println("text 和 weight 不应该含有空格：", line)
		}

		// 检查：code 前后不应该含有空格
		if strings.HasPrefix(code, " ") || strings.HasSuffix(code, " ") {
			fmt.Println("code 前后不应该含有空格：", line)
		}

		// 检查：编码是否含有非字母，或没有小写
		for _, r := range code {
			if string(r) != " " && !(r >= 97 && r <= 122) {
				fmt.Println("编码含有非字母或大写字母：", line)
				break
			}
		}

		// 检查：汉字部分是否含有非汉字内容
		for _, c := range text {
			if strings.Contains(text, "·") {
				break
			}
			if !unicode.Is(unicode.Han, c) {
				fmt.Println("含有非汉字内容：", line, c)
				break
			}
		}

		// 除了字表，其他词库不应该含有单个的汉字
		if dictPath != HanziPath && utf8.RuneCountInString(text) == 1 {
			fmt.Println("意外的单个汉字：", line)
		}

		// 除了 main ，其他词库不应该含有两个字的词汇
		if dictPath != MainPath && utf8.RuneCountInString(text) == 2 {
			fmt.Println("意外的两字词：", line)
		}

		// 汉字个数应该与拼音个数相等
		if code != "" {
			codeCount := len(strings.Split(code, " "))
			textCount := utf8.RuneCountInString(text)
			if strings.Contains(text, "·") {
				textCount -= strings.Count(text, "·")
			}
			if strings.HasPrefix(text, "# ") {
				textCount -= 2
			}
			if textCount != codeCount {
				fmt.Println("汉字个数和拼音个数不相等：", text, code)
			}
		}

		// +---------------------------------------------------------------
		// | 多音字注音问题检查
		// +---------------------------------------------------------------
		if dictPath == ExtPath || dictPath == TencentPath {
			for _, word := range polyphoneWords.ToSlice() {
				if strings.Contains(text, word) {
					fmt.Printf("多音字注音问题：%q\n", line)
				}
			}
		}

		// +---------------------------------------------------------------
		// | 异形词检查
		// +---------------------------------------------------------------
		if dictPath != HanziPath && !filterWords.Contains(text) {
			for _, wrongWord := range wrongWords.ToSlice() {
				if strings.Contains(text, wrongWord) {
					fmt.Printf("异形词汇: %s - %s\n", wrongWord, text)
				}
			}
		}
	}
}
