function date_translator(input, seg)
    if (input == "rq") then
        --- Candidate(type, start, end, text, comment)
        yield(Candidate("date", seg.start, seg._end, os.date("%Y-%m-%d"), ""))
        yield(Candidate("date", seg.start, seg._end, os.date("%Y/%m/%d"), ""))
        -- yield(Candidate("date", seg.start, seg._end, os.date("%Y年%m月%d日"), ""))
        yield(Candidate("date", seg.start, seg._end, os.date("%Y 年 %m 月 %d 日"), "")) -- 加些空格
        -- yield(Candidate("date", seg.start, seg._end, os.date("%m-%d"), ""))
        -- yield(Candidate("date", seg.start, seg._end, os.date("%m-%d-%Y"), ""))
    end
    if (input == "sj") then
        --- Candidate(type, start, end, text, comment)
        yield(Candidate("time", seg.start, seg._end, os.date("%H:%M"), ""))
        yield(Candidate("time", seg.start, seg._end, os.date("%H:%M:%S"), ""))
        yield(Candidate("time", seg.start, seg._end, os.date("%Y%m%d%H%M%S"), ""))
    end

		-- 增加一个 ISO 8601 的时间格式 （示例 2022-01-07T20:42:51+08:00）
		if (input == "dt") then
			yield(Candidate("datetime", seg.start, seg._end, os.date("%Y-%m-%dT%H:%M:%S+08:00"), ""))
		end

    -- @JiandanDream
    -- https://github.com/KyleBing/rime-wubi86-jidian/issues/54

    if (input == "xq") then
        local weakTab = {'日', '一', '二', '三', '四', '五', '六'}
        yield(Candidate("week", seg.start, seg._end, "周"..weakTab[tonumber(os.date("%w")+1)], ""))
        yield(Candidate("week", seg.start, seg._end, "星期"..weakTab[tonumber(os.date("%w")+1)], ""))
    end
end
