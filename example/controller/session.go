package controller

import "github.com/valyala/fasthttp"

import "lu/example/lib"

type SessionMiddleware struct {
	cookieKey string
}

type Session struct {
	username string
}

func NewSessionMiddleware(cookieKey string) *SessionMiddleware {
	return &SessionMiddleware{cookieKey}
}

func (mw *SessionMiddleware) Middleware(redis *lib.Redis) func(ctx *fasthttp.RequestCtx, next func(error)) {
	return func(ctx *fasthttp.RequestCtx, next func(error)) {
		sid := string(ctx.Request.Header.Cookie(mw.cookieKey))
		val, _ := redis.Get(sid)
		if val == nil {
			next(nil)
			return
		}
		ctx.SetUserValue(sid, val)
		next(nil)
	}
}
