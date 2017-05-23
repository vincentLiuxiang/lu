# lu
[![Build Status](https://travis-ci.org/vincentLiuxiang/lu.svg?branch=master_travis)](https://travis-ci.org/vincentLiuxiang/lu) [![Coverage Status](https://coveralls.io/repos/github/vincentLiuxiang/lu/badge.svg)](https://coveralls.io/github/vincentLiuxiang/lu) [![Go Report Card](https://goreportcard.com/badge/github.com/vincentLiuxiang/lu)](https://goreportcard.com/report/github.com/vincentLiuxiang/lu)

      ___
     |  /      ..    )))   
     | |     .    . (((
     | |    .    ||~~~~||
     | |___ .    | \__/ |   
     \_____/      \____/      version: 0.0.1

 **A high performance and Light-weighted go middleware web framework which is based on [fasthttp](https://github.com/valyala/fasthttp)**

 The MIT License

If you are a [node.js](https://nodejs.org/en/) developer, you will find lu is quite similar in usage with [connect](https://github.com/senchalabs/connect) and [express](https://github.com/expressjs/express), but performance better.

## install

```
go get -u github.com/vincentLiuxiang/lu
```

## test

```
go test -v -cover
```
ref [go test](https://blog.golang.org/cover) for more about go test 
## [example](https://github.com/lugolang/lu-example)

```go
package main

import (
	"github.com/valyala/fasthttp"
	"github.com/vincentLiuxiang/lu"
)

func main() {
	app := lu.New()
	app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(200)
		next(nil)
	})

	app.Get("/api", func(ctx *fasthttp.RequestCtx, next func(error)) {
		next(errors.New("something error occour\n"))
	})

	app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetBody([]byte("hello world\n"))
	})

	app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetBody([]byte(err.Error()))
	})

	app.Listen(":8080")
}

```
**result**

```
http://localhost:8080/test
200
hello world

http://localhost:8080/api
200
something error occour
```

## custom your server

```go
app := lu.New()
app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
	ctx.SetBody([]byte("helloworld"))
})
app.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
	ctx.SetBody([]byte("my name is go"))
})
server := &fasthttp.Server{
	Handler:     app.Handler,
	Concurrency: 1024 * 1024,
}
server.ListenAndServe(":8080")
```
## Use middleware
The first parameter of ```app.Use``` we call it **router**, and the second parameter we call it **middleware**, all the middlewares will be pushed to a ```stack``` inner lu.

when a http request comes, lu will compare ctx.Path() with []byte(router), The compare rules are below:  

* if ctx.Path() equals to  []byte(router) , it matches.

* if ctx.Path() starts with  []byte(router), and len(ctx.Path()) > len(router) , and ctx.Path()[len(router)] is '/' or '?',  it matches.

* if the router is "/"，it means this router matches any http request

http request will execute each middleware one-by-one until a middleware does not call next() within it.

```go
app.Use("/hello",func(ctx *fasthttp.RequestCtx, next func(error)){
	// ctx.Path() is starts with []byte("/hello"),
	// if len(ctx.Path()) > len("/hello") , ctx.Path()[len("/hello")] must be '/' or '?'
	next(nil)
})
app.Use("/world",func(ctx *fasthttp.RequestCtx, next func(error)){
	// ctx.Path() is starts with []byte("/world"),
	// if len(ctx.Path()) > len("/world") , ctx.Path()[len("/world")] must be '/' or '?'
	next(nil)
})
```

* for example

```
app.Use("/", ...)
app.Use("/api", ...)
app.Use("/test", ...)
```

> ```http://xxxx:xxx/test```  match "/", "/test"

> ```http://xxxx:xxx/test?xxx=xxx```  match "/", "/test"

> ```http://xxxx:xxx/test/hello``` match "/", "/test"

> ```http://xxxx:xxx/api```  match "/", "/api"

> ```http://xxxx:xxx/api/hello?xxx=xxx``` match "/", "/api"
>

## error-middleware and non-error-middleware

The  second parameter of ```app.Use``` method, can be to two different type

* func(ctx *fasthttp.RequestCtx, next func(error)), we call it non-error-middleware

* func(err error, ctx *fasthttp.RequestCtx, next func(error)), we call it error-middleware, only execute by call next(error) within a middleware

In **lu**, there are two stack arrays to store the midllewares which are accepted by app.Use. One is used to store non-error-middleware and the other one  is used to store error-middleware


No matter in which type of middleware, when it calls next(nil), next non-error-middleware will be excute if router match. But, when calls next(errors.New('some error')) , the program will skip all the non-error-middleware and directly execute the left first error-middleware if router match.

* for example

```go
app.Use("/",func(ctx *fasthttp.RequestCtx, next func(error)){
	next(nil)
})
app.Use("/",func(ctx *fasthttp.RequestCtx, next func(error)){
	next(errors.New("skip next non-error-middleware"))
})
app.Use("/",func(ctx *fasthttp.RequestCtx, next func(error)){
	fmt.Println("skip this non-error-middleware")
})
app.Use("/",func(err error, ctx *fasthttp.RequestCtx, next func(error)){
	fmt.Println(err.Error())
})
```
result:

```
skip next non-error-middleware
```
## response
* If an incoming http request doesn't match any router, lu will response a 404 statusCode and a "Not Found" string body.

```go
app.Use("/foo",func(ctx *fasthttp.RequestCtx, next func(error)){
	// no response
})
app.Use("/bar",func(ctx *fasthttp.RequestCtx, next func(error)){

})
```
miss all of the middlewares

```
http://xxxx:xxx/go

404
Not Found
```

* If an incoming http request match some routers, but all the matched middleware don't response to the client, lu will response a 200 statusCode and a "" string body (fasthttp default mechanism)

```go
app.Use("/foo",func(ctx *fasthttp.RequestCtx, next func(error)){
	// no response
})
app.Use("/bar",func(ctx *fasthttp.RequestCtx, next func(error)){

})
```
match /foo

```
http://xxxx:xxx/foo

200


```
* In the last middleware, no matter what type of the middleware,  if you call ```next```, lu will response a 404 statusCode and a "Not Found" string body. Because, there is no middleware after the last middleware.

```go
app.Use("/foo",func(ctx *fasthttp.RequestCtx, next func(error)){
	// no response
})
app.Use("/bar",func(ctx *fasthttp.RequestCtx, next func(error)){
	ctx.SetStatusCode(200)
	ctx.SetBody([]byte("helloworld"))
	next(nil) 
	// or 
	// next(errors.New("..."))
})
```
lu will ResetBody() and SetStatusCode(404) , SetBody("Not Found")

```
http://xxxx:xxx/bar
// not 
// 200
// helloworld

404
Not Found
```
* ```app. Finally```.  However, if you call ```next(nil)``` or ```next(error)``` in the last middleware. lu providers a Finally function, allow user to custom the response

```
app := New()
app.Use("/test", func(ctx *fasthttp.RequestCtx, next func(error)) {
	next(errors.New("error"))
})
app.Use("/", func(err error, ctx *fasthttp.RequestCtx, next func(error)) {
	ctx.SetStatusCode(302)
	ctx.SetBody([]byte("hello world"))
	next(errors.New("finally handle"))
})
app.Finally = func(err error, ctx *fasthttp.RequestCtx) {
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody([]byte(err.Error()))
		return
	}
	ctx.SetStatusCode(200)
	ctx.SetBody([]byte("hello world"))
}
app.Listen(":3005")
```
lu will SetStatusCode(500) , ctx.SetBody([]byte(err.Error()))

```
http://xxxx:xxx/test

500
finally handle
```

## Useful Middleware
* [static](https://github.com/lugolang/static) lu static file serving middleware, based on fasthttp.FS.

## api

* app.Finally func(err error, ctx *fasthttp.RequestCtx)

* app.Handler fasthttp.RequestHandler

* app.Use(router string, func(ctx *fasthttp.RequestCtx, next func(error)) register non-error-middleware

* app.Use(router string, func(err error, ctx *fasthttp.RequestCtx, next func(error))) register error-middleware

* app.Listen(port string) listen a port. app.Listen(":8080")

* app.Get(router string, func(ctx *fasthttp.RequestCtx, next func(error))) quite similar with app.Use but only handle http GET method

* app.Post(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http POST method

* app.Put(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http PUT method

* app.Head(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http HEAD method

* app.Patch(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http PATCH method

* app.Delete(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http DELETE method

* app.Options(router string, func(ctx *fasthttp.RequestCtx, next func(error))) only handle http OPTIONS method


