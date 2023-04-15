# 自用 Rime 配置

![demo](./others/demo.webp)

## 基本套路

- 简体 | 全拼 | 双拼
- 主要功能
    -   「[袖珍简化字方案](https://github.com/rime/rime-pinyin-simp)」作为基础
    -   [easy_en](https://github.com/BlindingDark/rime-easy-en) 英文输入
    -   简繁切换
    -   日期、时间、星期
    -   自整理的 Emoji
    -   [以词定字](https://github.com/BlindingDark/rime-lua-select-character)
    -   [长词优先](https://github.com/tumuyan/rime-melt/blob/master/lua/melt.lua)
    -   所有标点符号直接上屏，「/」模式改为「v」模式，「/」直接上屏
    -   增加了许多拼音纠错
- 简体字表、词库
    -   [《通用规范汉字表》](https://github.com/iDvel/The-Table-of-General-Standard-Chinese-Characters)
    -   [华宇野风系统词库](http://bbs.pinyin.thunisoft.com/forum.php?mod=viewthread&tid=30049)
    -   [清华大学开源词库](https://github.com/thunlp/THUOCL)
    -   [《现代汉语常用词表》](https://gist.github.com/indiejoseph/eae09c673460aa0b56db)
    -   [《现代汉语词典》](https://forum.freemdict.com/t/topic/12102)
    -   [《同义词词林》](https://forum.freemdict.com/t/topic/1211)
    -   [《新华成语大词典》](https://forum.freemdict.com/t/topic/11407)
    -   [搜狗网络流行新词](https://pinyin.sogou.com/dict/detail/index/4)
    -   [腾讯词向量](https://ai.tencent.com/ailab/nlp/en/download.html)
- 词库修订
    - 校对大量异形词、错别字、错误注音

<br>

## 长期维护词库

因为没有找到一份比较好的词库，干脆自己维护一个。综合了几个不错的词库，精心调教了很多。

主要维护的词库：

- `8105` 字表。
- `base` 基础词库。
- `sogou` 搜狗流行词。
- `ext` 扩展词库，小词库。
- `tencent` 扩展词库，大词库。
- Emoji

维护内容主要是异形词、错别字的校对，错误注音的修正，缺失的常用词汇的增添，词频的调整。

欢迎在词库方面提 issue，我会及时更新修正。

<br>

## 使用说明

### 1. 选项菜单
在输入状态时，<kbd>F4</kbd> 或者 <kbd>control</kbd> + <kbd>`</kbd> 弹出菜单

### 2. 菜单内容
弹出的菜单中，处于第一位的是当前使用的输入法方案，其后跟着是该方案中的输入法菜单，有【袖珍简化字拼音】、【简 ──> 繁繁繁】、【大写数字】、【特殊字符】等常见功能菜单，再后面是其它可选的输入法方案，对应 [`default.custom.yaml`](https://github.com/appleshan/fcitx5-rime-dict/blob/main/default.custom.yaml) 中 `schema_list` 字段内容

### 3. 默认二三候选
默认的二三候选是 <kbd>;</kbd> <kbd>'</kbd> 两个键

### 4. 候选翻页
方向 <kbd>上</kbd><kbd>下</kbd>、<kbd>-</kbd> <kbd>=</kbd>

### 5. 支持 简入繁出
是以切换输入方案的形式实现的，使用时，调出菜单，选择 `简 ──> 繁繁繁` 方案即可
简繁转换的功能能实现：
- 转繁体
> 以不切换文字的形式使用只是暂时转繁，换个程序就会恢复简体了。如果你想一直使用简入繁出就选择 「简入繁出」这个方案

### 6. 系统 `时间`、`日期` 和 `星期`
输入对应词，获取当前日期和时间
- `rq` 输出日期，格式 `2019年06月19日` `2019-06-19`
- `sj` 输出时间，格式 `10:00` `10:00:00`
- `xq` 输出星期，格式 `周四` `星期四`

### 7. 支持大写数字输入：壹贰叁肆伍陆
本库中包含一个可以输入大写数字的方案，名叫 `大写数字`，呼出菜单选择该方案即可。
在这个模式下：具体可以看源文件 [`numbers.schema.yaml`](https://github.com/appleshan/fcitx5-rime-dict/blob/main/numbers.schema.yaml)


| 键           | 对应值             | | 键 (按住 shift) | 对应值            |
|-------------|--------------------|---|-----------|-------------------|
| 1234567890  | 壹贰叁肆伍陆柒捌玖零  | | 1234567890 | 一二三四五六七八九〇  |
| wqbsjfd.    | 万仟佰拾角分第点     | | wqbsjfd.   | 万千百十角分点       |
| z           | 整之               | | z          | 整之               |
| y           | 元月亿             | | y          | 元月亿             |

### 8. 特殊字符快捷输入
默认是开启的，具体可以查看 wiki [如何启用 ` /fh` 这种特殊符号输入](https://github.com/KyleBing/rime-wubi86-jidian/wiki/%E5%A6%82%E4%BD%95%E5%90%AF%E7%94%A8-%60--fh%60-%E8%BF%99%E7%A7%8D%E7%89%B9%E6%AE%8A%E7%AC%A6%E5%8F%B7%E8%BE%93%E5%85%A5)

使用方法改成类似 `vfh` 按键，即以 `v` 作为前缀按键。

### 9. IntelliJ 系列软件的 IDE 中输入框位置不正确
此问题的根本原因是 IDE 附带的 JBR 不正确，要处理此问题，需要：
1. 前往 [Release 17.0.6-b829.5](https://github.com/RikudouPatrickstar/JetBrainsRuntime-for-Linux-x64/releases) 下载 jbr 并解压到任意路径
2. 按照 [此指导](https://intellij-support.jetbrains.com/hc/en-us/articles/206544879-Selecting-the-JDK-version-the-IDE-will-run-under) 更改 IDE 的 JBR

## 长期维护词库

主要维护的词库：

- `8105.dict.yaml` 字表。
- `base.dict.yaml` 基础词库。
- `sogou.dict.yaml` 搜狗流行词。
- `ext.dict.yaml` 扩展词库。
- `tencent.dict.yaml` 扩展词库，大词库。

维护内容主要是异形词、错别字的校对，错误注音的修正，缺失的常用词汇的增添，词频的调整。

<br>

## 参考资料
    - [优化英文输入体验](https://dvel.me/posts/make-rime-en-better/)

---

详细介绍：[我的 Rime 配置 2022](https://dvel.me/posts/my-rime-setting-2022/)
