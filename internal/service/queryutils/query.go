package queryutils

import (
	"github.com/obnahsgnaw/socketgwservice/config"
	"github.com/obnahsgnaw/socketgwservice/internal/service"
	"reflect"
	"strconv"
	"strings"
)

/*
type xx struct {
	// keyword filter
	Keyword string `query:"keyword"`

	// field filter
	Id     uint32 `query:"id"`
	Enable []int  `query:"enable" options:"1,2,3"`
}

type LogSort struct {
	Id        queryutils.SortVal `query:"id" default:"0"`
	CreatedAt queryutils.SortVal `query:"cat"`
}

*/

const (
	None SortVal = 0 // the default if not set
	Desc SortVal = 2 // sort desc
	Asc  SortVal = 1 // sort asc
)

// Page paginate params
type Page struct {
	Id    uint32
	Limit uint32
}

// AnyPtr any data of ptr
type AnyPtr interface{}

// SortVal the value of sort
type SortVal int

func (v SortVal) Asc() bool {
	return v == Asc
}

func (v SortVal) Desc() bool {
	return v == Desc
}

func (v SortVal) Valid() bool {
	return v != None
}

// ParseFilter string bool int uint float []int32, input =  xxx|id=1|name=xx|keyword=xxx|xx=1,2,3
func ParseFilter(des AnyPtr, input string) {
	fq := parseKeyword(input)
	if len(fq) == 0 {
		return
	}
	v := reflect.ValueOf(des)
	t := v.Type()
	if t.Kind() != reflect.Pointer {
		return
	}
	pv := v.Elem()
	pt := t.Elem()
	for i := 0; i < pt.NumField(); i++ {
		if pt.Field(i).Name == "Columns" {
			continue
		}
		key := pt.Field(i).Tag.Get("query")
		if key == "" {
			key = strings.ToLower(pt.Field(i).Name)
		} else {
			if key == "-" {
				continue
			}
		}
		val, ok := fq[key]
		if !ok {
			continue
		}
		switch pt.Field(i).Type.Kind() {
		case reflect.String:
			pv.Field(i).SetString(val)
			break
		case reflect.Bool:
			pv.Field(i).SetBool(val == "1" || val == "true" || val == "TRUE" || val == "True")
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vv, err := strconv.Atoi(val); err == nil {
				pv.Field(i).SetInt(int64(vv))
			}
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if vv, err := strconv.Atoi(val); err == nil {
				pv.Field(i).SetUint(uint64(vv))
			}
			break
		case reflect.Float32, reflect.Float64:
			if vv, err := strconv.ParseFloat(val, 64); err == nil {
				pv.Field(i).SetFloat(vv)
			}
			break
		case reflect.Slice:
			var options []string
			var valOptions []string
			option := pt.Field(i).Tag.Get("options")
			if option != "" {
				options = strings.Split(option, ",")
			}
			if strings.Contains(val, ",") {
				valOptions = strings.Split(val, ",")
			} else {
				valOptions = []string{val}
			}
			if len(options) > 0 {
				valOptions = filterSlice(valOptions, options)
			}
			if len(valOptions) > 0 {
				var val11 []int32
				for _, v1 := range valOptions {
					if v2, err2 := strconv.Atoi(v1); err2 == nil {
						val11 = append(val11, int32(v2))
					}
				}
				if len(val11) > 0 {
					pv.Field(i).Set(reflect.ValueOf(val11))
				}
			}
			break
		default:
			switch pt.Field(i).Type.String() {
			case "time.Time":
				pv.Field(i).Set(reflect.ValueOf(service.Str2Time(val)))
				break
			default:

			}
		}
	}
}

// ParseSort parse sort
func ParseSort(des AnyPtr, input string) {
	sq := parseSort(input)
	v := reflect.ValueOf(des)
	t := v.Type()
	if t.Kind() != reflect.Pointer {
		return
	}
	pv := v.Elem()
	pt := t.Elem()
	for i := 0; i < pt.NumField(); i++ {
		key := pt.Field(i).Tag.Get("query")
		if key == "" {
			key = strings.ToLower(pt.Field(i).Name)
		}
		var vv SortVal
		if val, ok := sq[key]; ok {
			if val {
				vv = Asc
			} else {
				vv = Desc
			}
		} else {
			defVal := pt.Field(i).Tag.Get("default")
			if defVal != "" {
				if defVal == "1" || defVal == "true" || defVal == "True" || defVal == "TRUE" {
					vv = Asc
				} else {
					vv = Desc
				}
			}
		}
		if vv != None {
			pv.Field(i).Set(reflect.ValueOf(vv))
		}
	}
}

func ParsePage(page, limit uint32) Page {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = config.PageLimit
	}
	if limit > config.PageLimitMax {
		limit = config.PageLimitMax
	}

	return Page{
		Id:    page,
		Limit: limit,
	}
}

func filterSlice(s, options []string) (des []string) {
	o1 := make(map[string]struct{})
	s1 := make(map[string]struct{})
	for _, o := range options {
		o1[o] = struct{}{}
	}
	for _, ss := range s {
		if _, ok := o1[ss]; ok {
			s1[ss] = struct{}{}
		}
	}
	if len(s1) > 0 {
		for ss := range s1 {
			des = append(des, ss)
		}
	}
	return
}

// sort = a,-b,+c
func parseSort(sort string) map[string]bool {
	var querySorts = make(map[string]bool)

	if sort != "" {
		var sorts []string
		if strings.Contains(sort, ",") {
			// 字符串解析
			sorts = strings.Split(sort, ",")
		} else {
			sorts = []string{sort}
		}
		// 解析为字段:bool map
		for _, fd := range sorts {
			asc := !strings.HasPrefix(fd, "-")
			fdName := strings.TrimPrefix(fd, "-")
			fdName = strings.TrimPrefix(fdName, "+")
			querySorts[fdName] = asc
		}
	}

	return querySorts
}

// keyword = xxx or keyword= xxx|id=1|name=xx|keyword=xxx|xx=1,2,3
func parseKeyword(keyword string) (data map[string]string) {
	data = make(map[string]string)

	if keyword != "" {
		var keywords []string
		if strings.Contains(keyword, "|") {
			keywords = strings.Split(keyword, "|")
		} else {
			keywords = []string{keyword}
		}
		for _, kwd := range keywords {
			if kwd != "" {
				if strings.Contains(kwd, "=") {
					kvs := strings.Split(kwd, "=")
					if kvs[0] != "" && kvs[1] != "" {
						data[kvs[0]] = kvs[1]
					}
				} else {
					data["keyword"] = kwd
				}
			}

		}
	}

	return
}
