package register

import (
	"github.com/obnahsgnaw/pbhttp/core/application"
)

var processors []func(p *Provider)
var Provide *Provider
var options []Option

func Register(p func(p *Provider)) {
	processors = append(processors, p)
}

func AddOption(option Option) {
	options = append(options, option)
}

func Exec(p *application.Project, o ...Option) {
	o = append(o, options...)
	Provide = newProvider(p, o...)
	for _, pp := range processors {
		pp(Provide)
	}
	processors = nil
}
