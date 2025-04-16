package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/api/service"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/version"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/index/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterBackendController("index", indexv1.IndexService_ServiceDesc, func() interface{} {
			return &IndexController{}
		})
		// register method md
		c.RegisterMetaData("", map[string]service.MdValParser{
			//
		})
		// register api
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return indexv1.RegisterIndexServiceHandlerServer(ctx, mux, impl.(*IndexController))
		}, true)
	})
}

type IndexController struct {
	Controller
	indexv1.UnimplementedIndexServiceServer
}

func (c *IndexController) Version(_ context.Context, _ *indexv1.VersionRequest) (*indexv1.VersionResponse, error) {
	return &indexv1.VersionResponse{
		Version: version.Version(),
	}, nil
}
