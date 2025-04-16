package fileasset

import (
	"context"
	"github.com/obnahsgnaw/api/service/autheduser"
	uploadv1 "github.com/obnahsgnaw/fileassetapi/gen/fileasset_backend_api/upload/v1"
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/errcode/errfac"
	"google.golang.org/grpc"
)

func GetUploadUrl(user autheduser.User, q *uploadv1.FetchUrlRequest) (resp *uploadv1.FetchUrlResponse, err error) {
	err = register.Provide.BackendRps().AuthedCall(project.FileService.Key(), user, func(ctx context.Context, conn *grpc.ClientConn) error {
		c := uploadv1.NewUploadServiceClient(conn)
		resp, err = c.FetchUrl(ctx, q)
		return err
	})
	if err != nil {
		err = errfac.NewError("file asset rpc call get upload url failed", err)
	}

	return
}
