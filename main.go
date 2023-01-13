package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/adrg/frontmatter"
	"github.com/kataras/iris/v12"
)

var config *Config

func main() {
	// 加载配置文件 config.yml
	var err error
	config, err = readConf("./config.yml")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	LoadData(config.BlogPath)
	go watch(config.BlogPath)
	app := iris.New()
	setupRoutes(app)
	setupLogger(app)
	hostAndPort := fmt.Sprintf(":%d", config.GetPort())
	app.Listen(hostAndPort)
}
