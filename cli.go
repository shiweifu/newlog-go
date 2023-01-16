package main

import (
	"github.com/urfave/cli"
	_ "github.com/urfave/cli/v2"
	_ "github.com/urfave/cli/v2/altsrc"
)

func NewCliApp() *cli.App {
	cliApp := cli.NewApp()
	cliApp.Name = "markdown-blog"
	cliApp.Usage = "Markdown Blog App"
	return cliApp
}
