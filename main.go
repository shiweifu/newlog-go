package main

import (
	"log"
	"os"

	_ "github.com/adrg/frontmatter"
)

var config *Config
var Env string

func main() {
	app := NewCliApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
