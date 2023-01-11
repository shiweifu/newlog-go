package main

import (
	"fmt"
	"os"

	_ "github.com/adrg/frontmatter"
	"github.com/kataras/iris/v12"
)

func main() {
	// 加载配置文件 config.yml
	config, err := readConf("./config.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	loadData(config.BlogPath)
	go watch(config.BlogPath)
	app := iris.New()
	setupRoutes(app)
	app.Listen(":8080")
}
