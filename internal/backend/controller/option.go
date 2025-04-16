package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/option"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/index/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterBackendController("option", indexv1.OptionsService_ServiceDesc, func() interface{} {
			return &OptionController{}
		})
		// register api
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return indexv1.RegisterOptionsServiceHandlerServer(ctx, mux, impl.(*OptionController))
		}, true)
	})
}

type OptionController struct {
	Controller
}

func (c *OptionController) OptionList(context.Context, *indexv1.OptionListRequest) (*indexv1.OptionListResponse, error) {
	return &indexv1.OptionListResponse{
		EnableState: option.NewEnumOptionProvider(option.Default, option.Enabled, option.Disabled).Backend(),
	}, nil
}
