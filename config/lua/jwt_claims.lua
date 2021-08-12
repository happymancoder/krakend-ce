function replace_path_with_jwt_claims( req )

    local authHeader = req:headers("Authorization")
    print(authHeader)
    local first = string.find(authHeader, "%.")
    local last = string.find(authHeader, "%.", first+1)
    local rawData = string.sub(authHeader, first+1, last-1)
    local decoded = from_base64(rawData)
    local jwtData = json_parse(decoded)

    local u = req:url()
    req:url(u:gsub(".JWT.sub", jwtData["sub"]) )
end

function replace_path_with_jwt_claims2( req )

    local u = req:url()
        print(u)
end

-- version: 0.0.1
-- code: Ketmar // Avalon Group
-- public domain

-- expand $var and ${var} in string
-- ${var} can call Lua functions: ${string.rep(' ', 10)}
-- `$' can be screened with `\'
-- `...': args for $<number>
-- if `...' is just a one table -- take it as args
function ExpandVars (s, ...)
  local args = {...};
  args = #args == 1 and type(args[1]) == "table" and args[1] or args;
  -- return true if there was an expansion
  local function DoExpand (iscode)
    local was = false;
    local mask = iscode and "()%$(%b{})" or "()%$([%a%d_]*)";
    local drepl = iscode and "\\$" or "\\\\$";
    s = s:gsub(mask, function (pos, code)
      if s:sub(pos-1, pos-1) == "\\" then return "$"..code;
      else was = true; local v, err;
        if iscode then code = code:sub(2, -2);
        else local n = tonumber(code);
          if n then v = args[n]; end;
        end;
        if not v then
          v, err = loadstring("return "..code); if not v then error(err); end;
          v = v();
        end;
        if v == nil then v = ""; end;
        v = tostring(v):gsub("%$", drepl);
        return v;
      end;
    end);
    if not (iscode or was) then s = s:gsub("\\%$", "$"); end;
    return was;
  end;

  repeat DoExpand(true); until not DoExpand(false);
  return s;
end;

function interp(s, tab)
  return (s:gsub('($%b{})', function(w) return tab[w:sub(3, -2)] or w end))
end

function url_decode(str)
   str = str:gsub("+", " ")
   str = str:gsub("%%(%x%x)", function(h)
      return string.char(tonumber(h,16))
   end)
   str = str:gsub("\r\n", "\n")
   return str
end

function replace_path_with_jwt_claims3( req )
    local userId = req:headers("sub")
    print(userId)
    local customerId = req:headers("customer_id")
    print(customerId)
    local countryCode = req:headers("country_code")
    print(countryCode)
    
        local authHeader = req:headers("Content-Type")
    print(authHeader)
    
        local u = req:url()
        print(u)
        
        local decoded = url_decode(u)
        print(decoded)
        
        local output1 = decoded:gsub("{{.JWT.sub}}", userId)
        local output2 = output1:gsub("{{.JWT.customer_id}}", customerId)
        local output3 = output2:gsub("{{.JWT.country_code}}", countryCode)
        
        print(output3)
req:url(output3)



--local uri = interp("http:test.com/${JWT.sub}/permissions", {JWT.sub = "77752"})
--print(uri)
        --local n = string.replace(u, "{JWT.sub}", authHeader)
        --local uri = u:gsub("%{.JWT.sub}%", authHeader)
    --req:url(uri)
end
