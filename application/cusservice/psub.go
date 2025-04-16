package cusservice

import (
	"github.com/obnahsgnaw/pbhttp/core/application"
	"github.com/obnahsgnaw/pbhttp/pkg/psub"
	driver2 "github.com/obnahsgnaw/pbhttp/pkg/psub/driver"
	"github.com/obnahsgnaw/socketgwservice/application/register"
)

func init() {
	//registerService(pubSub)
}

func pubSub(p *application.Project) (err error) {
	var driver driver2.Driver
	var psb *psub.PSub
	driver, err = driver2.NewDriver(p.Context(), p.Config().Register.Driver, p.Config().Register.Endpoints, p.Config().Register.Timeout)
	if err != nil {
		return
	}
	p.App().Application().AddRelease(driver2.Release)
	psb = psub.New(driver)
	register.AddOption(register.PubSub(psb))
	return
}
