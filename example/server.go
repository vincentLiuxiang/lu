package main

import (
	"github.com/vincentLiuxiang/lu/example/controller"
	"github.com/vincentLiuxiang/lu/example/lib"

	"github.com/lugolang/static"

	"github.com/vincentLiuxiang/lu"
)

func main() {
	app := lu.New()
	redis := lib.NewRedis("tcp", "localhost:6379", 60*60)
	sessionMw := controller.NewSessionMiddleware("sid")
	loginMw := controller.NewLoginMiddleware("sid")
	fs := static.DefaultFS
	Static := static.New(fs)

	app.Use("/static/login", Static)
	app.Post("/login", loginMw.Check(redis))
	app.Use("/", sessionMw.Middleware(redis))
	app.Use("/", loginMw.Middleware())
	app.Use("/static/home", Static)

	app.Listen(":8080")
}
