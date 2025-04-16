package cusservice

import "github.com/obnahsgnaw/pbhttp/core/application"

var services []func(*application.Project) error

func registerService(f func(*application.Project) error) {
	services = append(services, f)
}

func Exec(p *application.Project) error {
	for _, f := range services {
		if err := f(p); err != nil {
			return err
		}
	}
	return nil
}
