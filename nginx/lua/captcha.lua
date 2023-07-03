local ngx_token = ngx.var.token or ''
local ngx_hash = ngx.var.hash or ''
local ngx_ttl = ngx.var.ttl or ''
local ngx_remote_ip = ngx.var.remote_addr or ''

-- expire time to complete captcha in xxx seconds
local captcha_expire = 150
local captcha_server = '127.0.0.1'

local ck = require "resty.cookie"
local http = require "resty.http"

if ngx_token:len() < 1 or ngx_hash :len() < 1 or ngx_ttl :len() < 1 or ngx_remote_ip :len() < 1 then
  -- set cookie_name  cookie
  local cookie, err = ck:new()
  if not cookie then
      ngx.log(ngx.ERR, err)
      return
  end
  local cookie_name_cookieUrl = 'cookieUrl'
  -- get cookie_name_cookieUrl
  local cookie_name_cookieUrl_value, err = cookie:get(cookie_name_cookieUrl)
  if not cookie_name_cookieUrl_value then
    local control_cookie = ngx.var.scheme .. '://' .. ngx.var.host .. ngx.var.request_uri
    local ok, err = cookie:set({
        key = cookie_name_cookieUrl,
        value = control_cookie,
        domain = ngx.var.http_host,
        path = "/",
        httponly = true,
        expires = ngx.cookie_time(ngx.time()+captcha_expire),
    })
    if not ok then
        ngx.log(ngx.ERR, err)
    end
  end
  -- get cookie_name_cookieIP
  local cookie_name_cookieIP = 'cookieIP'
  -- local cookie_name_cookieIP_value, err = cookie:get(cookie_name_cookieIP)
  if not cookie_name_cookieIP_value then
    local control_cookie = ngx_remote_ip
    local ok, err = cookie:set({
        key = cookie_name_cookieIP,
        value = control_cookie,
        domain = ngx.var.http_host,
        path = "/",
        httponly = true,
        expires = ngx.cookie_time(ngx.time()+captcha_expire),
    })
    if not ok then
        ngx.log(ngx.ERR, err)
    end
  end

  ngx.redirect("http://" .. ngx.var.http_host .. "/captcha", ngx.HTTP_MOVED_TEMPORARILY);
else
  -- check token
  local hc = http:new()
  local gtimeout = 2000

  local ok, code, headers, status, body  = hc:request {
    url = "http://" .. captcha_server .. "/token/check/" .. ngx_token .. "/" .. ngx_hash .. "/" .. ngx.var.remote_addr .. "/" .. ngx_ttl,
    timeout = gtimeout,
      method = "GET",
      headers = { Host = ngx.var.http_host },
  }
  -- ngx.log(ngx.STDERR, code)
  if code ~= 200 then
    -- ngx.status = code
    -- ngx.say(body)
    return ngx.exit(code)
  end
end
