package lu

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type (
	handleFunc func(*fasthttp.RequestCtx, func(error))
	errorsFunc func(error, *fasthttp.RequestCtx, func(error))
	midware    struct {
		route      []byte
		index      int
		middleWare handleFunc
		errorsWare errorsFunc
	}
	stack []midware
)

func (stack *stack) push(m midware) {
	*stack = append(*stack, m)
}

type Lu struct {
	stack    stack
	errStack stack
}

func New() *Lu {
	return &Lu{}
}

func (lu *Lu) Use(route string, handler interface{}) {
	if route == "" || route[0] != '/' {
		panic("The first params of Use func must be a string which start with '/'")
	}

	if route == "/" {
		route = ""
	}

	midHandle, mOk := handler.(func(*fasthttp.RequestCtx, func(error)))
	errHandle, eOk := handler.(func(error, *fasthttp.RequestCtx, func(error)))

	if !mOk && !eOk {
		panic("The second params of Use func must be a" +
			"\n\tfunc(*fasthttp.RequestCtx, func(error)) or" +
			"\n\tfunc(error, *fasthttp.RequestCtx, func(error)) type")
	}

	if mOk {
		lu.stack.push(midware{
			route:      []byte(route),
			index:      len(lu.errStack),
			middleWare: midHandle})
		return
	}

	if eOk {
		lu.errStack.push(midware{
			route:      []byte(route),
			index:      len(lu.stack),
			errorsWare: errHandle})
		return
	}
}

func (lu *Lu) Post(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("POST"), handler)
}

func (lu *Lu) Put(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("PUT"), handler)
}

func (lu *Lu) Get(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("GET"), handler)
}

func (lu *Lu) Delete(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("DELETE"), handler)
}

func (lu *Lu) Options(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("OPTIONS"), handler)
}

func (lu *Lu) Patch(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("PATCH"), handler)
}

func (lu *Lu) Head(route string, handler handleFunc) {
	lu.httpMethod(route, []byte("HEAD"), handler)
}

func (lu *Lu) Listen(port string) error {
	return fasthttp.ListenAndServe(port, lu.Handler)
}

func (lu *Lu) Handler(ctx *fasthttp.RequestCtx) {
	var (
		index    = 0
		errIndex = 0
		nxt      func(error)
		err      error = nil
	)
	nxt = func(err error) {
		var m midware

		if err != nil {
			if errIndex >= len(lu.errStack) {
				nonMiddlwareResponese(ctx)
				return
			}
			m = lu.errStack[errIndex]
			errIndex++
			index = m.index
		} else {
			if index >= len(lu.stack) {
				nonMiddlwareResponese(ctx)
				return
			}
			m = lu.stack[index]
			index++
			errIndex = m.index
		}

		handle(err, ctx, m, nxt)
	}

	nxt(err)
}

func nonMiddlwareResponese(ctx *fasthttp.RequestCtx) {
	ctx.ResetBody()
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBodyString("Not Found")
}

func (lu *Lu) httpMethod(route string, method []byte, handler handleFunc) {
	lu.Use(route, func(ctx *fasthttp.RequestCtx, next func(error)) {
		if sliceCompare(ctx.Method(), method) {
			handler(ctx, next)
			return
		}
		next(nil)
	})
}

func sliceCompare(src, dest []byte) bool {
	if len(src) != len(dest) {
		return false
	}

	return sliceDiff(src, dest)
}

func sliceContains(src, dest []byte) bool {
	if len(src) < len(dest) {
		return false
	}

	return sliceDiff(src, dest)
}

func sliceDiff(src, dest []byte) bool {
	for i, w := range dest {
		if src[i] != w {
			return false
		}
	}
	return true
}

func handle(err error, ctx *fasthttp.RequestCtx, m midware, n func(error)) {
	url := ctx.Path()
	urlLen := len(url)
	rouLen := len(m.route)

	if !sliceContains(url, m.route) {
		n(err)
		return
	}

	if urlLen > rouLen && url[rouLen] != '/' && url[rouLen] != '?' {
		n(err)
		return
	}

	if err != nil {
		m.errorsWare(err, ctx, n)
		return
	}

	m.middleWare(ctx, n)
	return
}

func init() {
	version := "0.0.1"
	fmt.Printf(`  	  ___              
	 |  /      ..    )))   
	 | |     .    . (((
	 | |    .    ||~~~~||
	 | |___ .    | \__/ |   
	 \_____/      \____/    ` + "version: " + version + "\n\n")
}
