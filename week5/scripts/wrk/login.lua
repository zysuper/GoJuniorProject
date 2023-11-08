wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"
-- 让其可以通过用户指纹校验.
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 这个要改为你的注册的数据
wrk.body='{"email":"12347@qq.com", "password": "hello#world123"}'


function response(status, headers, body)
    print(headers["X-Jwt-Token"])
end