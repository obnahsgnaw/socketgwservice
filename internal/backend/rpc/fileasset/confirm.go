package fileasset

import (
	"context"
	"errors"
	"github.com/obnahsgnaw/api/service/autheduser"
	uploadv1 "github.com/obnahsgnaw/fileassetapi/gen/fileasset_backend_api/upload/v1"
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/errcode/errfac"
	"google.golang.org/grpc"
)

func Confirm(user autheduser.User, q *uploadv1.ConfirmRequest) error {
	err := register.Provide.BackendRps().AuthedCall(project.FileService.Key(), user, func(ctx context.Context, conn *grpc.ClientConn) error {
		c := uploadv1.NewUploadServiceClient(conn)
		resp, err := c.Confirm(ctx, q)
		if err != nil {
			return err
		}
		if !resp.Success {
			return errors.New(resp.Error)
		}
		return nil
	})
	if err != nil {
		return errfac.NewError("file asset rpc call confirm failed", err)
	}

	return nil
}
