package errcode

import "github.com/obnahsgnaw/socketgwservice/internal/errcode/errfac"

// 自定义业务错误码

var (
	DbFailed = errfac.NewStdErrCode(100)
)
