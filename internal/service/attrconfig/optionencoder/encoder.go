package optionencoder

import "strings"

func Encode(options []string) string {
	if len(options) == 0 {
		return ""
	}
	var m = make(map[string]struct{})
	for _, o := range options {
		m[o] = struct{}{}
	}
	var ss []string
	for k, _ := range m {
		ss = append(ss, k)
	}
	return strings.Join(ss, ",")
}

func Decode(options string) []string {
	if options == "" {
		return nil
	}
	return strings.Split(options, ",")
}
