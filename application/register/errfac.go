package register

import (
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/socketgwservice/config"
)

var errfac *apierr.Factory

func init() {
	errfac = apierr.New(config.Project.Id())
}

func ErrFactory() *apierr.Factory {
	return errfac
}
