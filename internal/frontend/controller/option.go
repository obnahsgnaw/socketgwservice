package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/option"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/index/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterFrontendController("option", indexv1.OptionsService_ServiceDesc, func() interface{} {
			return &optionController{}
		})
		// register api
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return indexv1.RegisterOptionsServiceHandlerServer(ctx, mux, impl.(*optionController))
		}, true)
	})
}

type optionController struct {
	Controller
}

func (c *optionController) OptionList(context.Context, *indexv1.OptionListRequest) (*indexv1.OptionListResponse, error) {
	return &indexv1.OptionListResponse{
		EnableState: option.NewEnumOptionProvider(option.Default, option.Enabled, option.Disabled).Frontend(),
	}, nil
}
