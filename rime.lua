-- Rime lua 扩展：https://github.com/hchunhui/librime-lua
-------------------------------------------------------------
-- 日期时间
-- 提高权重的原因：因为在方案中设置了大于 1 的 initial_quality，导致 rq sj xq dt ts 产出的候选项在所有词语的最后。
function date_translator(input, seg)
    -- 日期
    if (input == "rq") then
        local cand = Candidate("date", seg.start, seg._end, os.date("%Y-%m-%d"), "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("date", seg.start, seg._end, os.date("%Y/%m/%d"), "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("date", seg.start, seg._end, os.date("%Y.%m.%d"), "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("date", seg.start, seg._end, os.date("%Y 年 %m 月 %d 日"), "")
        cand.quality = 100
        yield(cand)
    end
    -- 时间
    if (input == "sj") then
        local cand = Candidate("time", seg.start, seg._end, os.date("%H:%M"), "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("time", seg.start, seg._end, os.date("%H:%M:%S"), "")
        cand.quality = 100
        yield(cand)
    end
    -- 星期
    if (input == "xq") then
        local weakTab = {'日', '一', '二', '三', '四', '五', '六'}
        local cand = Candidate("week", seg.start, seg._end, "周" .. weakTab[tonumber(os.date("%w") + 1)], "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("week", seg.start, seg._end, "星期" .. weakTab[tonumber(os.date("%w") + 1)], "")
        cand.quality = 100
        yield(cand)
    end
    -- ISO 8601/RFC 3339 的时间格式 （固定东八区）（示例 2022-01-07T20:42:51+08:00）
    if (input == "dt") then
        local cand = Candidate("datetime", seg.start, seg._end, os.date("%Y-%m-%dT%H:%M:%S+08:00"), "")
        cand.quality = 100
        yield(cand)
        local cand = Candidate("time", seg.start, seg._end, os.date("%Y%m%d%H%M%S"), "")
        cand.quality = 100
        yield(cand)
    end
    -- 时间戳（十位数，到秒，示例 1650861664）
    if (input == "ts") then
        local cand = Candidate("datetime", seg.start, seg._end, os.time(), "")
        cand.quality = 100
        yield(cand)
    end
end
-------------------------------------------------------------
-- 以词定字
-- https://github.com/BlindingDark/rime-lua-select-character
local function utf8_sub(s, i, j)
    i = i or 1
    j = j or -1

    if i < 1 or j < 1 then
        local n = utf8.len(s)
        if not n then
            return nil
        end
        if i < 0 then
            i = n + 1 + i
        end
        if j < 0 then
            j = n + 1 + j
        end
        if i < 0 then
            i = 1
        elseif i > n then
            i = n
        end
        if j < 0 then
            j = 1
        elseif j > n then
            j = n
        end
    end

    if j < i then
        return ""
    end

    i = utf8.offset(s, i)
    j = utf8.offset(s, j + 1)

    if i and j then
        return s:sub(i, j - 1)
    elseif i then
        return s:sub(i)
    else
        return ""
    end
end

local function first_character(s)
    return utf8_sub(s, 1, 1)
end

local function last_character(s)
    return utf8_sub(s, -1, -1)
end

function select_character(key, env)
    local engine = env.engine
    local context = engine.context
    local commit_text = context:get_commit_text()
    local config = engine.schema.config

    local first_key = config:get_string('key_binder/select_first_character') or 'bracketleft'
    local last_key = config:get_string('key_binder/select_last_character') or 'bracketright'

    if (key:repr() == first_key and commit_text ~= "") then
        engine:commit_text(first_character(commit_text))
        context:clear()

        return 1 -- kAccepted
    end

    if (key:repr() == last_key and commit_text ~= "") then
        engine:commit_text(last_character(commit_text))
        context:clear()

        return 1 -- kAccepted
    end

    return 2 -- kNoop
end
-------------------------------------------------------------
-- 长词优先（提升「西安」「提案」「图案」「饥饿」等词汇的优先级）
-- https://github.com/tumuyan/rime-melt
-- 修改：不提升英文和中英混输的
-- 目前是将3个词插入到第2、3、4位，我想改成插入到第3、4、5位，不知道怎么改。。。
function long_word_filter(input)
    local l = {}
    local length = 0 -- 记录第一个候选词的长度，提前的候选词至少要比第一个候选词长
    -- local s1 = 0 -- 记录筛选了多少个英语词条(只提升3个词的权重，并且对comment长度过长的候选进行过滤)
    local s2 = 0 -- 记录筛选了多少个汉语词条(只提升3个词的权重)
    for cand in input:iter() do
        leng = utf8.len(cand.text)
        if (length < 1) then
            length = leng
            yield(cand)
        -- 不知道这两行是干嘛用的，似乎注释掉也没有影响。
        -- elseif #table > 30 then
        --     table.insert(l, cand)
        -- elseif ((leng > length) and (s1 < 2)) and (string.find(cand.text, "^[%w%p%s]+$")) then
        --     s1 = s1 + 1
        --     if (string.len(cand.text) / string.len(cand.comment) > 1.5) then
        --         yield(cand)
        --     end
        elseif ((leng > length) and (s2 < 3)) and (string.find(cand.text, "[%w%p%s]+") == nil) then
            yield(cand)
            s2 = s2 + 1
        else
            table.insert(l, cand)
        end
    end
    for i, cand in ipairs(l) do
        yield(cand)
    end
end
-------------------------------------------------------------
-- 因为英文方案的 initial_quality 大于 1，导致输入「va」时，候选项是「van vain。。。」
-- 单字优先，候选项应改为「ā á ǎ à」
--
-- 不知道这个方法为什么不行啊？？？
-- function v_single_char_first_filter(input)
--     if (string.find(input, "v") == 1 and string.len(input) == 2) then
--         local l = {}
--         for cand in input:iter() do
--             if (utf8.len(cand.text) == 1) then
--                 yield(cand)
--             else
--                 table.insert(l, cand)
--             end
--         end
--         for cand in ipairs(l) do
--             yield(cand)
--         end
--     end
-- end
--
-- 反正是解决了，不知道怎么就解决了，就是最后多一个候选项，没多大影响。
function v_single_char_first_filter(input, seg)
    if (string.find(input, "v") == 1 and string.len(input) == 2) then
        yield(Candidate("", seg.start, seg._end, "", ""))
    end
end
-------------------------------------------------------------
