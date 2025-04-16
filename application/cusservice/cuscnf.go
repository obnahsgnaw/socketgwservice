package cusservice

import (
	"github.com/obnahsgnaw/pbhttp/core/application"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/config"
)

var cusCnf *config.Config

func init() {
	registerService(initCusConfig)
}

func initCusConfig(p *application.Project) error {
	cusCnf = &config.Config{}
	if err := p.Config().ParseCustomConfig(cusCnf); err != nil {
		return err
	}
	register.AddOption(register.CusConfig(cusCnf))
	return nil
}
