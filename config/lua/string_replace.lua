
-- USAGE:
-- string.replace("mystring", "my", "our")
-- or
-- local teststr = "weird[]str%ing"
-- teststr2 = teststr:replace("weird[]", "cool(%1)")

-- Warning: add your own \0 char handling if you need it!

do
	local function regexEscape(str)
		return str:gsub("[%(%)%.%%%+%-%*%?%[%^%$%]]", "%%%1")
	end
	-- you can use return and set your own name if you do require() or dofile()
	
	-- like this: str_replace = require("string-replace")
	-- return function (str, this, that) -- modify the line below for the above to work
	string.replace = function (str, this, that)
		return str:gsub(regexEscape(this), that:gsub("%%", "%%%%")) -- only % needs to be escaped for 'that'
	end
end