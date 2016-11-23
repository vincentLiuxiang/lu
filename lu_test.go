package lu

import (
	"bytes"
	"net"
	"testing"

	"fmt"

	"strings"

	"errors"

	"github.com/valyala/fasthttp"
)

func Test_RouterError(t *testing.T) {
	defer func() {
		if err := recover(); err == "The first params of Use func must be a string which start with '/'" {
			t.Log("OK")
		} else {
			t.Error("ERROR")
		}
	}()
	app := New()
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
	})
	app.Use("hello", func(ctx *fasthttp.RequestCtx, next func(error)) {
	})
}

func Test_MiddleWareError(t *testing.T) {
	defer func() {
		if err := recover(); strings.Contains(err.(string), "The second params of Use func must be a") {
			t.Log("OK")
		} else {
			fmt.Println(err)
			t.Error("ERROR")
		}
	}()
	app := New()
	app.Use("/", func() {
	})
}

func Test_Use(t *testing.T) {
	var (
		Get, Post, Head, Delete, Put, Patch, Options bool
	)
	app := New()
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Get = true
		Post = true
		Head = true
		Delete = true
		Put = true
		Patch = true
		Options = true
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// Get
	Get = false
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Get {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Post
	Post = false
	rw.rbuf.WriteString("POST / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Post {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Head
	Head = false
	rw.rbuf.WriteString("HEAD / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Head {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Delete
	Delete = false
	rw.rbuf.WriteString("DELETE / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Delete {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Put
	Put = false
	rw.rbuf.WriteString("PUT / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Put {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Patch
	Patch = false
	rw.rbuf.WriteString("PATCH / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Patch {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Options
	Options = false
	rw.rbuf.WriteString("OPTIONS / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Options {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

}

func Test_HttpMethod(t *testing.T) {
	var (
		Get, Post, Head, Delete, Put, Patch, Options bool
	)
	app := New()
	app.Get("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Get = true
	})

	app.Post("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Post = true
	})

	app.Head("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Head = true
	})

	app.Delete("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Delete = true
	})

	app.Put("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Put = true
	})

	app.Patch("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Patch = true
	})

	app.Options("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		Options = true
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// Get
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Get {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Post
	rw.rbuf.WriteString("POST / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Post {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Head
	rw.rbuf.WriteString("HEAD / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Head {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Delete
	rw.rbuf.WriteString("DELETE / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Delete {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Put
	rw.rbuf.WriteString("PUT / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Put {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Patch
	rw.rbuf.WriteString("PATCH / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Patch {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Options
	rw.rbuf.WriteString("OPTIONS / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if Options {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_Router(t *testing.T) {
	var allTest, test1, test2, test3 bool
	app := New()

	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		allTest = true
		test1 = false
		test2 = false
		test3 = false
		// next(errors.New("error"))
		next(nil)
	})
	app.Use("/test1", func(ctx *fasthttp.RequestCtx, next func(error)) {
		test1 = true
		test2 = false
		test3 = false
	})

	app.Use("/test2", func(ctx *fasthttp.RequestCtx, next func(error)) {
		test1 = false
		test2 = true
		test3 = false
	})

	app.Use("/test3", func(ctx *fasthttp.RequestCtx, next func(error)) {
		test1 = false
		test2 = false
		test3 = true
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// non-router
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if allTest && !test1 && !test2 && !test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// /test1
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test1 HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if allTest && test1 && !test2 && !test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// /test1
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test1/hello HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if allTest && test1 && !test2 && !test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// /test2
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test2/hello/world HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if allTest && !test1 && test2 && !test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// /test3
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test3?user=123 HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if !test1 && !test2 && test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// non-router
	allTest = false
	test1 = false
	test2 = false
	test3 = false
	rw.rbuf.WriteString("GET /test3.user=123 HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if allTest && !test1 && !test2 && !test3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_ErrorMiddleWare1(t *testing.T) {
	app := New()
	var errorMiddleWare, middleware bool
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})

	// skip this middleware
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleware = true
	})

	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		errorMiddleWare = true
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// errorMiddleWare
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if !middleware && errorMiddleWare {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareOrder1(t *testing.T) {
	app := New()
	var errorMiddleWare bool
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})
	// skip this middleware
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		errorMiddleWare = true
	})

	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// errorMiddleWare
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if !errorMiddleWare {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareOrder2(t *testing.T) {
	app := New()
	var middleWare int = 0
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})
	// skip this middleware
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare += 1
		next(nil)
	})

	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		next(nil)
	})

	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare += 3
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// errorMiddleWare
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if middleWare == 3 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareOrder3(t *testing.T) {
	app := New()
	var middleWare int = 100
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare = 0
		next(nil)
	})

	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare += 1
		next(nil)
	})

	// skip this error-middleware
	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare += 3
		next(nil)
	})

	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		middleWare += 5
	})

	rw := &netConn{}
	s := &fasthttp.Server{
		Handler: app.Handler,
	}

	// errorMiddleWare
	rw.rbuf.WriteString("GET / HTTP/1.1\r\n\r\n")
	s.ServeConn(rw)
	if middleWare == 6 {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareResponese1(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {

	})
	// No response
	app.Use("/test/hello", func(ctx *fasthttp.RequestCtx, next func(error)) {

	})

	go app.Listen(":8080")

	// Non router match,
	// lu will response 404 Not Found
	code, body, _ := fasthttp.Get(nil, "http://localhost:8080/test1")
	if code == 404 && string(body) == "Not Found" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}

	// Middleware /test match , but no response
	// The 200 code is set by fasthttp default
	code, body, _ = fasthttp.Get(nil, "http://localhost:8080/test")
	if code == 200 && string(body) == "" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareResponese2(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(nil)
	})
	// No response
	app.Use("/test/hello", func(ctx *fasthttp.RequestCtx, next func(error)) {
		// next(nil)
	})

	go app.Listen(":8081")

	code, body, _ := fasthttp.Get(nil, "http://localhost:8081/test")
	if code == 404 && string(body) == "Not Found" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_MiddleWareResponese3(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(nil)
	})
	// Router don't match, skip
	app.Use("/test/world", func(ctx *fasthttp.RequestCtx, next func(error)) {
		// next(nil)
	})
	// Router match, response
	app.Use("/test/hello", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(302)
		ctx.SetBody([]byte("hello"))
	})

	go app.Listen(":8085")

	code, body, _ := fasthttp.Get(nil, "http://localhost:8085/test/hello")

	if code == 302 && string(body) == "hello" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_ErrorMiddleWareResponese1(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})
	// No response
	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {

	})

	go app.Listen(":8082")

	code, body, _ := fasthttp.Get(nil, "http://localhost:8082/test")
	if code == 200 && string(body) == "" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
		// t.Skip()
	}
}

func Test_ErrorMiddleWareResponese2(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})
	// No response
	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})

	go app.Listen(":8083")

	code, body, _ := fasthttp.Get(nil, "http://localhost:8083/test")
	if code == 404 && string(body) == "Not Found" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

func Test_ErrorMiddleWareResponese3(t *testing.T) {
	app := New()

	// No response
	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("error"))
	})
	// No response
	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(302)
		ctx.SetBody([]byte("world"))
	})

	go app.Listen(":8086")

	code, body, _ := fasthttp.Get(nil, "http://localhost:8086/test")
	if code == 302 && string(body) == "world" {
		t.Log("OK")
	} else {
		t.Error("ERROR")
	}
}

type netConn struct {
	net.Conn
	rbuf bytes.Buffer
	wbuf bytes.Buffer
}

func (rw *netConn) Close() error {
	return nil
}

func (rw *netConn) Read(b []byte) (int, error) {
	return rw.rbuf.Read(b)
}

func (rw *netConn) Write(b []byte) (int, error) {
	return rw.wbuf.Write(b)
}
