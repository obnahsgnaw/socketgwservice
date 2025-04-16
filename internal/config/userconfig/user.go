package appconfig

import (
	"github.com/obnahsgnaw/api/service/autheduser"
)

func Demo(user autheduser.User) bool {
	v, ok := user.Attr("demo")
	return ok && v == "1"
}
