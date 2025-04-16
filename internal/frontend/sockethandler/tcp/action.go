package tcp

import (
	demov1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_tcp_api/demo/v1"
	"github.com/obnahsgnaw/socketutil/codec"
)

var (
	demoConnectReq = codec.Action{
		Id:   codec.ActionId(demov1.ActionId_DemoConnectReq),
		Name: demov1.ActionId_DemoConnectReq.String(),
	}
	demoConnectResp = codec.Action{
		Id:   codec.ActionId(demov1.ActionId_DemoConnectResp),
		Name: demov1.ActionId_DemoConnectResp.String(),
	}
)
