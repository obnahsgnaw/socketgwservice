package fileasset

import (
	"context"
	"github.com/obnahsgnaw/api/service/autheduser"
	viewv1 "github.com/obnahsgnaw/fileassetapi/gen/fileasset_backend_api/view/v1"
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/errcode/errfac"
	"google.golang.org/grpc"
)

func GetViewUrls(user autheduser.User, q *viewv1.ViewUrlsRequest) (resp *viewv1.ViewUrlsResponse, err error) {
	err = register.Provide.BackendRps().AuthedCall(project.FileService.Key(), user, func(ctx context.Context, conn *grpc.ClientConn) error {
		c := viewv1.NewViewServiceClient(conn)
		resp, err = c.ViewUrls(ctx, q)
		return err
	})
	if err != nil {
		err = errfac.NewError("file asset rpc call get view urls failed", err)
	}

	return
}

func GetViewUrl(user autheduser.User, q *viewv1.ViewUrlRequest) (resp *viewv1.ViewUrlResponse, err error) {
	err = register.Provide.BackendRps().AuthedCall(project.FileService.Key(), user, func(ctx context.Context, conn *grpc.ClientConn) error {
		c := viewv1.NewViewServiceClient(conn)
		resp, err = c.ViewUrl(ctx, q)
		return err
	})
	if err != nil {
		err = errfac.NewError("file asset rpc call get view url failed", err)
	}

	return
}
