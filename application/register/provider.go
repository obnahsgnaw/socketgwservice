package register

import (
	"context"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/api/service"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/pbhttp/core/application"
	config2 "github.com/obnahsgnaw/pbhttp/core/config"
	"github.com/obnahsgnaw/pbhttp/core/project"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/index/v1"
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/sockethandler/sockettype"
	"github.com/obnahsgnaw/socketutil/codec"
	"google.golang.org/grpc"
	"net/http"
)

type Provider struct {
	project     *application.Project
	cusServices map[string]interface{}
}

func newProvider(p *application.Project, o ...Option) *Provider {
	s := &Provider{project: p, cusServices: make(map[string]interface{})}
	for _, opt := range o {
		opt(s)
	}
	return s
}

func (s *Provider) addCusService(name string, val interface{}) {
	s.cusServices[name] = val
}

func (s *Provider) Context() context.Context {
	return s.project.Context()
}

func (s *Provider) Project() project.Service {
	return s.project.Project()
}

func (s *Provider) Config() *config2.Config {
	return s.project.Config()
}

func (s *Provider) ErrFactory() *apierr.Factory {
	return s.project.ErrFactory()
}

func (s *Provider) CommonConfig() *config2.CommonConfig {
	return s.Config().CommonConfig()
}

func (s *Provider) ServiceEnabled(p project.Service) bool {
	return s.CommonConfig().ServiceEnabled(p)
}

func (s *Provider) ProjectConfig() *indexv1.Config {
	return s.Config().ProjectConfig().Raw().(*indexv1.Config)
}

func (s *Provider) ProjectConfigUpdate(cc *indexv1.Config) error {
	return s.Config().ProjectConfig().Sync(cc)
}

func (s *Provider) App() *pbhttp.Application {
	return s.project.App()
}

func (s *Provider) Database() *pbhttp.Database {
	return s.App().Database()
}

func (s *Provider) CacheProvider() *pbhttp.CacheDriver {
	return s.App().Cache().CacheDriver()
}

func (s *Provider) Service() *pbhttp.Service {
	return s.App().Service()
}

func (s *Provider) Frontend() *pbhttp.EndServer {
	return s.App().Frontend()
}

func (s *Provider) Backend() *pbhttp.EndServer {
	return s.App().Backend()
}

func (s *Provider) FrontendRps() *pbhttp.Rps {
	return s.Frontend().ApiServer().RpcServer()
}

func (s *Provider) BackendRps() *pbhttp.Rps {
	return s.Backend().ApiServer().RpcServer()
}

func (s *Provider) RegisterTcpHandler(act codec.Action, structure action.DataStructure, handler action.Handler) {
	s.Frontend().SocketHandler().RegisterHandler(sockettype.TCP, act, structure, handler)
}

func (s *Provider) TcpHandler() *sockethandler.Handler {
	return s.Frontend().SocketHandler().TrustGet(sockettype.TCP)
}

func (s *Provider) WssHandler() *sockethandler.Handler {
	return s.Frontend().SocketHandler().TrustGet(sockettype.WSS)
}

func (s *Provider) RegisterWssHandler(act codec.Action, structure action.DataStructure, handler action.Handler) {
	s.Frontend().SocketHandler().RegisterHandler(sockettype.WSS, act, structure, handler)
}

func (s *Provider) RegisterBackendController(name string, desc grpc.ServiceDesc, impl func() interface{}) *pbhttp.ControllerConfig {
	return s.project.App().Backend().Controller().Register(name, desc, impl)
}

func (s *Provider) BackendController(name string) interface{} {
	return s.Backend().Controller().TrustGet(name).GetImpl()
}

func (s *Provider) RegisterFrontendController(name string, desc grpc.ServiceDesc, impl func() interface{}) *pbhttp.ControllerConfig {
	return s.project.App().Frontend().Controller().Register(name, desc, impl)
}

func (s *Provider) FrontendController(name string) interface{} {
	return s.Frontend().Controller().TrustGet(name).GetImpl()
}

func (s *Provider) RegisterBackendRouter(meth string, pathPattern string, h func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)) {
	s.project.App().Backend().Router().Register(meth, pathPattern, h)
}

func (s *Provider) RegisterBackendStaticRouter(meth string, pathPattern string, h func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)) {
	s.project.App().Backend().Router().RegisterStatic(meth, pathPattern, h)
}

func (s *Provider) RegisterBackendRawRouter(p service.RouteProvider) {
	s.project.App().Backend().Router().RegisterRaw(p)
}

func (s *Provider) RegisterFrontendRouter(meth string, pathPattern string, h func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)) {
	s.project.App().Frontend().Router().Register(meth, pathPattern, h)
}

func (s *Provider) RegisterFrontendStaticRouter(meth string, pathPattern string, h func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)) {
	s.project.App().Frontend().Router().RegisterStatic(meth, pathPattern, h)
}

func (s *Provider) RegisterFrontendRawRouter(p service.RouteProvider) {
	s.project.App().Frontend().Router().RegisterRaw(p)
}
