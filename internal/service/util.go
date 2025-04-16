package service

import (
	"github.com/obnahsgnaw/api/service/autheduser"
	"github.com/obnahsgnaw/application/pkg/utils"
	"time"
)

func OperatorString(operator autheduser.User) string {
	return utils.ToStr("[管理员(", operator.Name(), ")]")
}

func parseStr2Time(layout []string, input string) time.Time {
	for _, tt := range layout {
		t, err := time.ParseInLocation(tt, input, time.Local)
		if err == nil && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

func Str2Time(input string) time.Time {
	return parseStr2Time([]string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
	}, input)
}

func Str2Date(input string) time.Time {
	return parseStr2Time([]string{
		"2006-01-02",
		"2006/01/02",
	}, input)
}
