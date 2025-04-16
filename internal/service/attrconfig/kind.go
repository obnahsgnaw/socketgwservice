package attrconfig

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type AttrKind uint32

func (v AttrKind) Val() uint32 {
	return uint32(v)
}

const (
	KindText     AttrKind = 1
	KindInt      AttrKind = 2
	KindFloat    AttrKind = 3
	KindBool     AttrKind = 4
	KindOption   AttrKind = 5
	KindDatetime AttrKind = 6
	KindDate     AttrKind = 7
	KindTime     AttrKind = 8
)

func validateValue(kind AttrKind, value string, options []string) (validVal string, err error) {
	if value == "" {
		return
	}
	switch kind {
	case KindText:
		validVal = value
		return
	case KindInt:
		if v, err1 := strconv.Atoi(value); err1 != nil {
			err = errors.New("invalid int value")
			return
		} else {
			validVal = strconv.Itoa(v)
		}
		return
	case KindFloat:
		if v, err1 := strconv.ParseFloat(value, 64); err1 != nil {
			err = errors.New("invalid float value")
			return
		} else {
			validVal = strconv.FormatFloat(v, 'f', 10, 64)
		}
		return
	case KindBool:
		if value == "1" || strings.ToLower(value) == "true" {
			validVal = "1"
		} else {
			validVal = "0"
		}
		return
	case KindOption:
		for _, o := range options {
			if value == o {
				validVal = value
				return
			}
		}
		err = errors.New("invalid option value")
		return
	case KindDatetime:
		if v, err1 := time.Parse("2006-01-02 15:04:05", value); err1 != nil {
			err = errors.New("invalid datetime value")
			return
		} else {
			validVal = v.Format("2006-01-02 15:04:05")
		}
		return
	case KindDate:
		if v, err1 := time.Parse("2006-01-02", value); err1 != nil {
			err = errors.New("invalid date value")
			return
		} else {
			validVal = v.Format("2006-01-02")
		}
		return
	case KindTime:
		if v, err1 := time.Parse("15:04:05", value); err1 != nil {
			err = errors.New("invalid time value")
			return
		} else {
			validVal = v.Format("15:04:05")
		}
		return

	default:
		err = errors.New("invalid kind")
		return
	}
}
