package main

import (
	"log"
	"os"

	_ "github.com/adrg/frontmatter"
	"github.com/kataras/iris/v12"
)

func main() {
	// 加载配置文件 config.yml
	config, err := readConf("./config.yml")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	LoadData(config.BlogPath)
	go watch(config.BlogPath)
	app := iris.New()
	setupRoutes(app)
	app.Listen(":8080")
}
