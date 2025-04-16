package register

import (
	"github.com/obnahsgnaw/pbhttp/pkg/psub"
)

func PubSub(psb *psub.PSub) Option {
	return func(p *Provider) {
		p.addCusService("psb", psb)
	}
}

func (s *Provider) PubSub() *psub.PSub {
	return s.cusServices["psb"].(*psub.PSub)
}
