package tcp

import (
	"context"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	demov1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_tcp_api/demo/v1"

	//"github.com/obnahsgnaw/socketgwservice/app/frontend/frontreg"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func init() {
	register.Register(func(p *register.Provider) {
		/*
			p.RegisterTcpHandler(demoConnectReq, func() codec.DataPtr {
				return &demov1.DemoConnectRequest{}
			}, demoConnectHandler)
		*/
	})
}

func demoConnectHandler(_ context.Context, rq *action.HandlerReq) (respAct codec.Action, response codec.DataPtr, err error) {
	rqData := rq.Data.(*demov1.DemoConnectRequest)
	respAct = demoConnectResp
	resp := &demov1.DemoConnectResponse{
		Status:   demov1.DemoConnectResponse_InvalidArgument,
		Datetime: timestamppb.New(time.Now()),
	}
	response = resp

	if err1 := rqData.Validate(); err1 != nil {
		return
	}

	// 逻辑
	resp.Status = demov1.DemoConnectResponse_Success

	return
}
