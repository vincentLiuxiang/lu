package controller

import (
	"encoding/json"
	"fmt"
	"lu/example/lib"

	"github.com/valyala/fasthttp"
)

type LoginMiddleware struct {
	cookieKey string
}

type Login struct {
	username string
}

type _loginBody struct {
	Username string `json:username`
	Password string `json:password`
}

func NewLoginMiddleware(cookieKey string) *LoginMiddleware {
	return &LoginMiddleware{cookieKey}
}

func (mw *LoginMiddleware) Middleware() func(ctx *fasthttp.RequestCtx, next func(error)) {
	return func(ctx *fasthttp.RequestCtx, next func(error)) {
		sid := string(ctx.Request.Header.Cookie(mw.cookieKey))
		session := ctx.UserValue(sid)
		if session == nil {
			ctx.Redirect("/static/login/", 302)
			return
		}
		next(nil)
	}
}

func (mw *LoginMiddleware) Check(redis *lib.Redis) func(ctx *fasthttp.RequestCtx, next func(error)) {
	return func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(200)
		loginBody := &_loginBody{}
		err := json.Unmarshal(ctx.PostBody(), &loginBody)

		if err != nil {
			ctx.SetBodyString("{\"success\": false,\"noteMsg\":\"post body error\"}")
			return
		}

		fmt.Println(loginBody.Password, loginBody.Username)

		username := loginBody.Username
		password := loginBody.Password

		if username != "lu" || password != "fasthttp" {
			ctx.SetBodyString("{\"success\": false,\"noteMsg\":\"username or password invalid\"}")
			return
		}

		sid := "xhxhh3h8xhjs92jsj2qz==="
		ctx.Response.Header.Set("Set-Cookie", "sid="+sid+"; max-age=10000;")
		redis.Set(sid, "liuxiang")
		ctx.SetBodyString("{\"success\": true, \"redirect\":\"/static/home/index.html\"}")
	}
}
