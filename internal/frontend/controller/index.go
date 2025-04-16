package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/api/service"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/version"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/index/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterFrontendController("index", indexv1.IndexService_ServiceDesc, func() interface{} {
			return &indexController{}
		})
		// register method md
		c.RegisterMetaData("", map[string]service.MdValParser{
			//
		})
		// register api
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return indexv1.RegisterIndexServiceHandlerServer(ctx, mux, impl.(*indexController))
		}, true)
	})
}

type indexController struct {
	Controller
}

func (c *indexController) Version(_ context.Context, _ *indexv1.VersionRequest) (*indexv1.VersionResponse, error) {
	return &indexv1.VersionResponse{
		Version: version.Version(),
	}, nil
}
