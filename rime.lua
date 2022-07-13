-- 日期时间
function date_translator(input, seg)
    -- 日期
    if (input == "rq") then
        --- Candidate(type, start, end, text, comment)
        yield(Candidate("date", seg.start, seg._end, os.date("%Y-%m-%d"), ""))
        yield(Candidate("date", seg.start, seg._end, os.date("%Y/%m/%d"), ""))
        yield(Candidate("date", seg.start, seg._end, os.date("%Y.%m.%d"), ""))
        -- yield(Candidate("date", seg.start, seg._end, os.date("%Y年%m月%d日"), ""))
        yield(Candidate("date", seg.start, seg._end, os.date("%Y 年 %m 月 %d 日"), "")) -- 加些空格
        -- yield(Candidate("date", seg.start, seg._end, os.date("%m-%d"), ""))
        -- yield(Candidate("date", seg.start, seg._end, os.date("%m-%d-%Y"), ""))
    end

    -- 时间
    if (input == "sj") then
        --- Candidate(type, start, end, text, comment)
        yield(Candidate("time", seg.start, seg._end, os.date("%H:%M"), ""))
        yield(Candidate("time", seg.start, seg._end, os.date("%H:%M:%S"), ""))
    end

    -- 星期  @JiandanDream  https://github.com/KyleBing/rime-wubi86-jidian/issues/54
    if (input == "xq") then
        local weakTab = {'日', '一', '二', '三', '四', '五', '六'}
        yield(Candidate("week", seg.start, seg._end, "周"..weakTab[tonumber(os.date("%w")+1)], ""))
        yield(Candidate("week", seg.start, seg._end, "星期"..weakTab[tonumber(os.date("%w")+1)], ""))
    end

    -- ISO 8601/RFC 3339 的时间格式 （固定东八区）（示例 2022-01-07T20:42:51+08:00）
    if (input == "dt") then
        yield(Candidate("datetime", seg.start, seg._end, os.date("%Y-%m-%dT%H:%M:%S+08:00"), ""))
        yield(Candidate("time", seg.start, seg._end, os.date("%Y%m%d%H%M%S"), ""))
    end

    -- 时间戳（十位数，到秒，示例 1650861664）
    if (input == "ts") then
        yield(Candidate("datetime", seg.start, seg._end, os.time(), ""))
    end
end


-- 以词定字  https://github.com/BlindingDark/rime-lua-select-character
-- 设定的快捷键请看 select_character() 方法下面的注释
local function utf8_sub(s, i, j)
    i = i or 1
    j = j or -1

    if i < 1 or j < 1 then
         local n = utf8.len(s)
         if not n then return nil end
         if i < 0 then i = n + 1 + i end
         if j < 0 then j = n + 1 + j end
         if i < 0 then i = 1 elseif i > n then i = n end
         if j < 0 then j = 1 elseif j > n then j = n end
    end

    if j < i then return "" end

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
