package controller

import (
	"context"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/errcode"
	"github.com/obnahsgnaw/socketgwservice/internal/service/attrconfig"
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	v1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/common/v1"
)

func init() {
	register.Register(func(p *register.Provider) {
		/*
			c := RegisterController("xxx-attrConfig", attrconfigv1.AttrConfigService_ServiceDesc, func() interface{} {
				return &AttrConfigController{
					s: , // TODO 注入 backend attrConfigController 的 service
				}
			})
			c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
				return attrconfigv1.RegisterAttrConfigServiceHandlerServer(ctx, mux, impl.(*AttrConfigController))
			}, true)
		*/
	})
}

type AttrConfigController struct {
	Controller
	s *attrconfig.ConfigService
}

func (c *AttrConfigController) AttrConfigs(ctx context.Context, q *v1.AttrConfigRequest) (resp *v1.AttrConfigsResponse, err error) {
	var md *pbhttp.MetaData

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}

	list, _, err := c.s.GetConfigs(attrconfig.NewTarget(q.TargetId, ""), true, queryutils.Page{}, attrconfig.ConfigFilter{}, true)
	if err != nil {
		return nil, c.Failed(apierr.NewInternalError(errcode.DbFailed, err))
	}
	resp = &v1.AttrConfigsResponse{}
	for _, item := range list {
		resp.Attrs = append(resp.Attrs, &v1.AttrConfig{
			Attr:    item.Attr,
			Title:   item.Title,
			Kind:    v1.AttrConfig_Kind(item.Kind),
			Value:   item.Value,
			Options: item.Options,
		})
	}
	return
}
