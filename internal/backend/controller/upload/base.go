package upload

import (
	"context"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/pbhttp/pkg/cache"
	"github.com/obnahsgnaw/socketgwservice/internal/backend/controller"
	"github.com/obnahsgnaw/socketgwservice/internal/service/fileasset"
	uploadv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/upload/v1"
)

type Controller struct {
	controller.Controller
	cache cache.Cache
	s     *fileasset.Server
}

func (c *Controller) Config(ctx context.Context, q *uploadv1.ConfigRequest) (resp *uploadv1.ConfigResponse, err error) {
	var md *pbhttp.MetaData
	var cf *fileasset.Config

	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if cf, err = c.s.Config(md.User, q.ReqId); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	resp = &uploadv1.ConfigResponse{
		SessionId:    cf.SessionId,
		MaxSize:      cf.MaxSize,
		ContentTypes: cf.ContentTypes,
		Extensions:   cf.Extensions,
		Ttl:          cf.Ttl,
		Multipart:    cf.Multipart,
		MaxCount:     cf.MaxCount,
	}
	return
}

func (c *Controller) Url(ctx context.Context, q *uploadv1.UrlRequest) (resp *uploadv1.UrlResponse, err error) {
	var md *pbhttp.MetaData
	var uploadId string
	var name string
	var urls []string
	if md = c.GetMdData(ctx, true); md.Err != nil {
		return nil, c.Failed(md.Err)
	}
	if err = q.Validate(); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	if uploadId, name, urls, err = c.s.Url(md.User, q.SessionId, q.ContentType, q.Extension, q.PartNum); err != nil {
		return nil, c.Failed(apierr.NewValidateError(err.Error()))
	}
	resp = &uploadv1.UrlResponse{
		UploadId: uploadId,
		Name:     name,
		Urls:     urls,
	}
	return
}
