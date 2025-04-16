package main

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/pbhttp/core/application"
	application2 "github.com/obnahsgnaw/socketgwservice/application"
	"github.com/obnahsgnaw/socketgwservice/config"
	"time"
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600) // 东8区
}

var app *application.Project

func main() {
	utils.RecoverHandler(config.Project.Name(), func(err, stack string) {
		app.Log(err + ",stack=" + stack)
	})

	app = application2.NewProject()
	defer app.Release()

	if err := application2.Init(app); err != nil {
		app.Exit(1, err)
	}

	app.Run(func(err error) {
		app.Exit(2, err)
	})

	app.Wait()
}
