-- go abbreviations plugin
-- triggers by space or '='.
-- see gen.go for abbreviations list.

local config = import("micro/config")
local util = import("micro/util")
local micro = import("micro")

abbrevs = loadfile(os.getenv("HOME") .. "/.config/micro/plug/goabbrevs/abbrevs.lua")
abbrevs()

local str = ""

function replace(bp, repl, back)
	bp:SelectWordLeft()
	bp.Cursor:DeleteSelection()
	bp.Buf:Insert(-bp.Cursor.Loc, repl)
	loc = bp.Cursor:Move(back, bp.Buf)
	bp.Cursor:GotoLoc(loc)
end

micro.Log("Start")

function onRune(bp, r)
	local ft = bp.Buf:FileType()
	if ft ~= "go" then
		return false
	end

	if r == " " or r == "=" then
		local repl, back = expand(str)
		if repl ~= "" then
			replace(bp, repl, back)
		end
		micro.Log("str:", str, "r:", r)
		str = ""
		return false
	end
	if util.IsWordChar(r) or r == ";" or r == "." then
		str = str .. r
	else
		str = ""
	end
	-- micro.Log("str:", str, "r:", r)
	return false
end

function reset(bp)
	local ft = bp.Buf:FileType()
	if ft ~= "go" then
		return false
	end
	str = ""
	return false
end

function onCursorUp(bp)
	reset(bp)
end

function onCursorDown(bp)
	reset(bp)
end

function onCursorLeft(bp)
	reset(bp)
end


function onCursorRight(bp)
	reset(bp)
end

function onBackspace(bp)
	reset(bp)
end
