package appconfig

import "github.com/obnahsgnaw/api/service/authedapp"

func Demo(app authedapp.App) bool {
	v, ok := app.Attr("demo")
	return ok && v == "1"
}
