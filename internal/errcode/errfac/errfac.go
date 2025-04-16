package errfac

import (
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
)

func NewError(msg string, err error) error {
	return pbhttp.NewError(msg, err)
}

func NewErrCode(code uint32, msgHandler apierr.ErrMsgHandler) apierr.ErrCode {
	return register.ErrFactory().NewErrCode(code, msgHandler)
}

// NewStdErrCode return a new ErrCode with no message handler
func NewStdErrCode(code uint32) apierr.ErrCode {
	return register.ErrFactory().NewStdErrCode(code)
}

// NewMsgErrCode return a ErrCode with string message
func NewMsgErrCode(code uint32, msg string) apierr.ErrCode {
	return register.ErrFactory().NewMsgErrCode(code, msg)
}

// NewCommonErrCode return a common ErrCode
func NewCommonErrCode(code uint32, msg string) apierr.ErrCode {
	return register.ErrFactory().NewCommonErrCode(code, msg)
}
