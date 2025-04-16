package controller

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obnahsgnaw/api/pkg/apierr"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"google.golang.org/grpc/metadata"
)

type Controller struct {
}

func (c *Controller) GetMdData(ctx context.Context, authed bool) *pbhttp.MetaData {
	return register.Provide.Frontend().Meta().Parse(ctx, authed)
}

func (c *Controller) GetRawMd(ctx context.Context) metadata.MD {
	if h, ok := metadata.FromIncomingContext(ctx); ok {
		return h
	}

	return metadata.New(map[string]string{})
}

func (c *Controller) GetMdString(md metadata.MD, key string) string {
	v := md.Get(key)
	if len(v) >= 1 {
		return v[0]
	}
	return ""
}

func (c *Controller) GetMdVal(md metadata.MD, key string) []string {
	return md.Get(key)
}

func (c *Controller) Failed(err error) *runtime.HTTPStatusError {
	return apierr.ToStatusError(err)
}

func (c *Controller) Created(data interface{}) *runtime.HTTPStatusError {
	err := apierr.NewCreated()
	err.Data = data
	return apierr.ToStatusError(err)
}

func (c *Controller) Deleted() *runtime.HTTPStatusError {
	return apierr.ToStatusError(apierr.NewDeleted())
}
