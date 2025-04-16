package config

import (
	gatewayv1 "github.com/obnahsgnaw/socketgateway/service/proto/gen/gateway/v1"
	"github.com/obnahsgnaw/socketutil/codec"
)

func CodecProvider() codec.Provider {
	toData := func(p *codec.PKG) codec.DataPtr {
		if p == nil {
			return &gatewayv1.GatewayPackage{}
		}
		return &gatewayv1.GatewayPackage{
			Action: p.Action.Val(),
			Data:   p.Data,
		}
	}
	toPkg := func(d codec.DataPtr) *codec.PKG {
		d1 := d.(*gatewayv1.GatewayPackage)
		return &codec.PKG{
			Action: codec.ActionId(d1.Action),
			Data:   d1.Data,
		}
	}
	p := codec.NewTcpProvider(toData, toPkg)
	p.SetDelimiter([]byte("\\N\\B"))
	p.SetMagicNum(0xAB)
	p.SetBodyMax(1024)
	return p
}
