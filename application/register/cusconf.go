package register

import "github.com/obnahsgnaw/socketgwservice/config"

func CusConfig(cusCnf *config.Config) Option {
	return func(p *Provider) {
		p.addCusService("cusCnf", cusCnf)
	}
}

func (s *Provider) CusConfig() *config.Config {
	return s.cusServices["cusCnf"].(*config.Config)
}
