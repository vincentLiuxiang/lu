### A login example

`cookie and session`


```go
package main

import (
	"lu/example/controller"
	"lu/example/lib"

	"github.com/lugolang/static"

	"github.com/vincentLiuxiang/lu"
)

func main() {
	app := lu.New()
  	// get a redis conn
	redis := lib.NewRedis("tcp", "localhost:6379", 60*60)
  	// get a session middleware instance
  	// sid is the key-value pair's key stored in cookie
	sessionMw := controller.NewSessionMiddleware("sid")
  	// get a login middleware instance
	loginMw := controller.NewLoginMiddleware("sid")
  	// init a static file middleware
	fs := static.DefaultFS
	Static := static.New(fs)
	
    // login page /static/login/index.html
	app.Get("/static/login", Static)
  	// Post, pass password, username to the server, check the truth and sign cookie / session, this is a restful api
	app.Post("/login", loginMw.Check(redis))
  	/**
	* app.Get("/getName", ...)
	* app.Post("/register", ...)
	*/
  	// sign session to any request via ctx.UserValue["session"]
	app.Use("/", sessionMw.Middleware(redis))
  	// judge if a request really has session info via cookie
  	// if ctx.UserValue["session"] == nil, redirect to /static/login/
	app.Use("/", loginMw.Middleware())
  	// if ctx.UserValue["session"] != nil, visi home page
	app.Get("/static/home", Static)

	app.Listen(":8080")
}
```

