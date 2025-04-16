package service

import "github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"

type AnyPtr queryutils.AnyPtr

type Service struct {
	//
}

// input =  xxx|id=1|name=xx|keyword=xxx|xx=1,2,3
/*
type xx struct {
	// keyword filter
	Keyword string `query:"keyword"`

	// field filter: int8-64 uint8-64 float32-64 bool string []int32 time.Time
	Id     uint32 `query:"id"`
	Enable []int32  `query:"enable" options:"1,2,3"`
}
*/

func (s *Service) ParseFilter(input string, ptr AnyPtr) {
	queryutils.ParseFilter(ptr, input)
}

// input = sort = a,-b,+c
/*
	type xx struct {
		Id        SortVal `query:"id" default:"0"`
		CreatedAt SortVal `query:"cat"`
	}
*/

func (s *Service) ParseSort(input string, ptr AnyPtr) {
	queryutils.ParseSort(ptr, input)
}

func (s *Service) ParsePage(id, limit uint32) queryutils.Page {
	return queryutils.ParsePage(id, limit)
}
