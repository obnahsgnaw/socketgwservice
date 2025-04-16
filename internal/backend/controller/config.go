package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/index/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterBackendController("config", indexv1.ConfigService_ServiceDesc, func() interface{} {
			return &ConfigController{p: p}
		})
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return indexv1.RegisterConfigServiceHandlerServer(ctx, mux, impl.(*ConfigController))
		}, true)
	})
}

type ConfigController struct {
	Controller
	p *register.Provider
}

func (c *ConfigController) Detail(ctx context.Context, q *indexv1.ConfigRequest) (resp *indexv1.Config, err error) {
	var md *pbhttp.MetaData
	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	resp = c.p.ProjectConfig()
	return
}

func (c *ConfigController) Update(ctx context.Context, q *indexv1.Config) (resp *indexv1.Config, err error) {
	var md *pbhttp.MetaData
	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if err = c.p.ProjectConfigUpdate(q); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	resp = c.p.ProjectConfig()

	return
}
