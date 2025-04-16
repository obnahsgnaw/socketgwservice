package upload

import (
	//"context"
	//"github.com/obnahsgnaw/socketgwservice/application"
	//"github.com/obnahsgnaw/socketgwservice/config"
	//"github.com/obnahsgnaw/socketgwservice/internal/backend/controller"
	//"github.com/obnahsgnaw/socketgwservice/internal/config/fileconfig"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/service/fileasset"
	//"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	//uploadv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/upload/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		/*
			c := controller.RegisterController(fileconfig.DemoService().Module()+"-upload", uploadv1.DemoUploadService_ServiceDesc, func() interface{} {
				return &DemoUploadController{
					Controller: Controller{
						cache: p.CacheProvider().Cache,
						s:     fileconfig.DemoService(),
					},
				}
			})
			c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
				return uploadv1.RegisterDemoUploadServiceHandlerServer(ctx, mux, impl.(*DemoUploadController))
			}, true)
		*/
	})
}

type DemoUploadController struct {
	Controller
}

func (c *DemoUploadController) Service() *fileasset.Server {
	return c.s
}
