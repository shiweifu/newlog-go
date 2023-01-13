package main

import "github.com/kataras/iris/v12"

func setupLogger(app *iris.Application) {
	app.Logger().SetLevel("debug")
}
