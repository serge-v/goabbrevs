-- go abbreviations plugin
-- triggers by space or '='.
-- see gen.go for abbreviations list.

local config = import("micro/config")
local util = import("micro/util")

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
		str = ""
		return false
	end
	if util.IsWordChar(r) or r == ";" or r == "." then
		str = str .. r
	else
		str = ""
	end
	return false
end
