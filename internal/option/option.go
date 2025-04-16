package option

import bcommonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
import fcommonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/common/v1"

type Option interface {
	String() string
	Is(int32) bool
	Value() int32
}

type EnumOptionProvider struct {
	data []Option
}

func NewEnumOptionProvider(o ...Option) *EnumOptionProvider {
	return &EnumOptionProvider{data: o}
}

func (s *EnumOptionProvider) AddItem(o Option) {
	s.data = append(s.data, o)
}

func (s *EnumOptionProvider) Backend() (resp []*bcommonv1.IntOption) {
	for _, o := range s.data {
		resp = append(resp, &bcommonv1.IntOption{
			Id:   o.Value(),
			Name: o.String(),
		})
	}
	return
}

func (s *EnumOptionProvider) Frontend() (resp []*fcommonv1.IntOption) {
	for _, o := range s.data {
		resp = append(resp, &fcommonv1.IntOption{
			Id:   o.Value(),
			Name: o.String(),
		})
	}
	return
}
