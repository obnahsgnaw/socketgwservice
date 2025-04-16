package controller

import (
	"context"
	"github.com/obnahsgnaw/socketgateway"
	"github.com/obnahsgnaw/socketgateway/pkg/socket/sockettype"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	v1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
	docv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/doc/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func init() {
	register.Register(func(p *register.Provider) {
		c := p.RegisterBackendController("doc", docv1.DocService_ServiceDesc, func() interface{} {
			return &docController{
				s: func() *socketgateway.DocServer {
					return p.Frontend().SocketGateway().TrustGet(sockettype.TCP).DocServer()
				},
			}
		})
		c.RegisterRpc()
	})
}

type docController struct {
	Controller
	s func() *socketgateway.DocServer
}

func (c *docController) List(context.Context, *emptypb.Empty) (*v1.StringOptionResponse, error) {
	resp := &v1.StringOptionResponse{List: nil}

	if c.s != nil {
		for k, v := range c.s().Manager.GetModules(true) {
			resp.List = append(resp.List, &v1.StringOption{
				Id:   strings.Replace(c.s().IndexDocUrl(), ":md", k, 1),
				Name: v,
			})
		}
	}
	return resp, nil
}
