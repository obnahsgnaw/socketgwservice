package controller

import (
	"context"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/errcode"
	"github.com/obnahsgnaw/socketgwservice/internal/service/attrconfig"
	"github.com/obnahsgnaw/socketgwservice/internal/service/conflict"
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	commonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	register.Register(func(p *register.Provider) {
		/*c := RegisterController("xxx-attrConfig", attrconfigv1.AttrConfigService_ServiceDesc, func() interface{} {
			return &AttrConfigController{
				s: attrconfig.NewConfig("", nil, nil), // TODO 注入repo
			}
		})
		c.RegisterApiService(func(ctx context.Context, mux *runtime.ServeMux, impl interface{}) error {
			return attrconfigv1.RegisterAttrConfigServiceHandlerServer(ctx, mux, impl.(*AttrConfigController))
		}, true)*/
	})
}

type AttrConfigController struct {
	Controller
	s *attrconfig.ConfigService
}

func (c *AttrConfigController) Paginate(ctx context.Context, q *commonv1.AttrConfigPaginateAllRequest) (resp *commonv1.AttrConfigPaginateResponse, err error) {
	var md *pbhttp.MetaData
	var configs []*attrconfig.Config
	var total int64
	page := queryutils.Page{}
	filter := attrconfig.ConfigFilter{}
	target := attrconfig.NewTarget(q.TargetId, "") // TODO 其他服务需要替换名称

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if !q.All {
		page = queryutils.ParsePage(q.Page, q.Limit)
		queryutils.ParseFilter(&filter, q.Keyword)
	}
	if configs, total, err = c.s.GetConfigs(target, q.All, page, filter, false); err != nil {
		return nil, c.Failed(apierr.NewInternalError(errcode.DbFailed, err))
	}
	resp = &commonv1.AttrConfigPaginateResponse{
		Meta: &commonv1.PaginateMeta{
			Total: uint32(total),
			Page:  page.Id,
			Limit: page.Limit,
		},
		List: nil,
	}

	for _, cc := range configs {
		resp.List = append(resp.List, toExtAttr(cc))
	}

	return resp, nil
}
func (c *AttrConfigController) Refresh(ctx context.Context, q *commonv1.AttrConfigRequest) (resp *commonv1.AttrConfigData, err error) {
	var md *pbhttp.MetaData
	var hit bool
	var conflicted bool
	target := attrconfig.NewTarget(q.TargetId, "") // TODO 其他服务需要替换名称

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	config := &attrconfig.Config{
		Kind:    attrconfig.AttrKind(q.Config.Kind),
		Attr:    q.Config.Attr,
		Title:   q.Config.Title,
		Value:   q.Config.Value,
		Options: q.Config.Options,
	}
	if hit, conflicted, err = c.s.RefreshConfig(md.User, target, config, q.Config.Conflict.Value); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}

	if hit {
		if conflicted {
			return nil, c.Failed(apierr.NewConflictError())
		}
		return toExtAttr(config), nil
	} else {
		return nil, c.Created(toExtAttr(config))
	}
}

func (c *AttrConfigController) Enable(ctx context.Context, q *commonv1.AttrConfigRequest) (resp *commonv1.AttrConfigData, err error) {
	var md *pbhttp.MetaData
	var hit bool
	var conflicted bool
	var config *attrconfig.Config
	target := attrconfig.NewTarget(q.TargetId, "") // TODO 其他服务需要替换名称

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if config, hit, conflicted, err = c.s.EnDisableConfig(md.User, target, q.Config.Attr, true, q.Config.Conflict.Value); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if !hit {
		return nil, c.Failed(apierr.NewValidateError("属性不存在"))
	}
	if conflicted {
		return nil, c.Failed(apierr.NewConflictError())
	}
	return toExtAttr(config), nil
}
func (c *AttrConfigController) Disable(ctx context.Context, q *commonv1.AttrConfigRequest) (resp *commonv1.AttrConfigData, err error) {
	var md *pbhttp.MetaData
	var hit bool
	var conflicted bool
	var config *attrconfig.Config
	target := attrconfig.NewTarget(q.TargetId, "") // TODO 其他服务需要替换名称

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if config, hit, conflicted, err = c.s.EnDisableConfig(md.User, target, q.Config.Attr, false, q.Config.Conflict.Value); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if !hit {
		return nil, c.Failed(apierr.NewValidateError("属性不存在"))
	}
	if conflicted {
		return nil, c.Failed(apierr.NewConflictError())
	}
	return toExtAttr(config), nil
}

func (c *AttrConfigController) Delete(ctx context.Context, q *commonv1.AttrConfigIdRequest) (_ *emptypb.Empty, err error) {
	var md *pbhttp.MetaData
	var hit bool
	target := attrconfig.NewTarget(q.TargetId, "") // TODO 其他服务需要替换名称

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if hit, err = c.s.RemoveConfig(md.User, target, q.Attr); err != nil {
		return nil, c.Failed(apierr.NewInternalError(errcode.DbFailed, err))
	}
	if !hit {
		return nil, c.Failed(apierr.NewValidateError("属性不存在"))
	}
	return nil, c.Deleted()
}
func (c *AttrConfigController) Service() *attrconfig.ConfigService {
	return c.s
}
func toExtAttr(cc *attrconfig.Config) *commonv1.AttrConfigData {
	item := &commonv1.AttrConfigData{
		Config: &commonv1.AttrConfig{
			Attr:     cc.Attr,
			Title:    cc.Title,
			Kind:     commonv1.AttrConfig_Kind(cc.Kind),
			Value:    cc.Value,
			Options:  cc.Options,
			Conflict: &commonv1.Conflict{Value: conflict.Gen(cc.UpdatedAt)},
		},
		Enabled: commonv1.EnableState(cc.Enabled.Value()),
		OperateInfo: &commonv1.OperateInfo{
			CreatedAt: timestamppb.New(cc.CreatedAt),
			UpdatedAt: timestamppb.New(cc.UpdatedAt),
		},
	}
	if cc.DeletedAt != nil {
		item.OperateInfo.DeletedAt = timestamppb.New(*cc.DeletedAt)
	}
	return item
}
