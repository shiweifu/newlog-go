package main

import (
	"os"
	"time"

	"github.com/kataras/iris/v12"
)

func setupLogger(app *iris.Application) {
	logsPath := "./logs"

	// check if logs directory exists
	_, err := os.Stat(logsPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logsPath, 0777)
		}
	}
	var logFilePath string = logsPath + "/log_" + time.Now().Format("20060102") + ".log"

	f, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	app.Logger().SetOutput(f)
	if Env != "production" {
		app.Logger().SetLevel("debug")
		app.Logger().Debugf(`debug mode\n`)
	}
	app.ConfigureHost(func(su *iris.Supervisor) {
		su.RegisterOnShutdown(func() {
			f.Close()
		})
	})

}
