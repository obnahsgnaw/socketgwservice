package application

import (
	"errors"
	"github.com/obnahsgnaw/api/pkg/errobj"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/servertype"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/pbhttp/core/application"
	config2 "github.com/obnahsgnaw/pbhttp/core/config"
	"github.com/obnahsgnaw/socketgateway/pkg/socket/engine/custom/mqtt"
	"github.com/obnahsgnaw/socketgateway/pkg/socket/sockettype"
	"github.com/obnahsgnaw/socketgwservice/application/cusservice"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/config"
	"github.com/obnahsgnaw/socketgwservice/config/middleware"
	_ "github.com/obnahsgnaw/socketgwservice/internal"
	"github.com/obnahsgnaw/socketgwservice/version"
	"github.com/obnahsgnaw/socketgwserviceapi/doc"
	bcommonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
	fcommonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_frontend_api/common/v1"
)

var name string

func SetName(n string) {
	name = n
}

var channels = map[string]map[servertype.ServerType]struct{}{
	"outer": {
		servertype.Tcp: struct{}{},
	},
	"inner": {
		servertype.Wss: struct{}{},
	},
}

func SetSupportChannel(ch string, sct servertype.ServerType) {
	if _, ok := channels[ch]; !ok {
		channels[ch] = make(map[servertype.ServerType]struct{})
	}
	channels[ch][sct] = struct{}{}
}

func SetOuterSupportServerType(scts ...servertype.ServerType) {
	if _, ok := channels["outer"]; !ok {
		channels["outer"] = make(map[servertype.ServerType]struct{})
	}

	for _, sct := range scts {
		channels["outer"][sct] = struct{}{}
	}
}

func SetInnerSupportServerType(scts ...servertype.ServerType) {
	if _, ok := channels["inner"]; !ok {
		channels["inner"] = make(map[servertype.ServerType]struct{})
	}
	for _, sct := range scts {
		channels["inner"][sct] = struct{}{}
	}
}

func NewProject() *application.Project {
	return application.NewProject(config.Project)
}

func Init(p *application.Project) error {
	if name != "" {
		p.Rename(name)
	}
	p.SetRawProjectConfig(func(c *config2.Config) config2.RawConfig {
		return config.NewConfig(c)
	})
	p.SetVersionProvider(func() string {
		return version.Info().String()
	})
	p.SetErrObjProvider(endtype.Backend, func(param errobj.Param) interface{} {
		d := &bcommonv1.ErrorsObject{
			Code:    param.Code,
			Message: param.Message,
		}

		if len(param.Errors) > 0 {
			for _, er := range param.Errors {
				d.Errors = append(d.Errors, &bcommonv1.ErrorObject{
					Code:    er.Code,
					Message: er.Message,
				})
			}
		}
		return d
	})
	p.SetErrObjProvider(endtype.Frontend, func(param errobj.Param) interface{} {
		d := &fcommonv1.ErrorsObject{
			Code:    param.Code,
			Message: param.Message,
		}

		if len(param.Errors) > 0 {
			for _, er := range param.Errors {
				d.Errors = append(d.Errors, &fcommonv1.ErrorObject{
					Code:    er.Code,
					Message: er.Message,
				})
			}
		}
		return d
	})
	p.SetDocProvider(endtype.Frontend, servertype.Api, func() ([]byte, error) {
		return doc.FrontendDoc()
	})
	p.SetDocProvider(endtype.Backend, servertype.Api, func() ([]byte, error) {
		return doc.BackendDoc()
	})
	p.SetDocProvider(endtype.Frontend, servertype.Tcp, func() ([]byte, error) {
		return doc.TcpDoc()
	})
	p.SetDocProvider(endtype.Frontend, servertype.Wss, func() ([]byte, error) {
		return doc.WssDoc()
	})
	p.SetAppInitializer(func(pa *pbhttp.Application) {
		// ignorer
		pa.Backend().ApiServer().RegisterAppIgnorer(config.BackendAppIgnorer)
		pa.Backend().ApiServer().RegisterUserIgnorer(config.BackendAuthIgnorer)
		pa.Backend().ApiServer().RegisterPermIgnorer(config.BackendPermIgnorer)
		pa.Frontend().ApiServer().RegisterAppIgnorer(config.FrontendAppIgnorer)
		pa.Frontend().ApiServer().RegisterUserIgnorer(config.FrontendAuthIgnorer)
		pa.Frontend().ApiServer().RegisterPermIgnorer(config.FrontendPermIgnorer)
		for mn, hdl := range middleware.FrontendMiddleware {
			pa.Frontend().ApiServer().RegisterMiddleware(mn, hdl, true)
		}
		for mn, hdl := range middleware.FrontendMuxMiddleware {
			pa.Frontend().ApiServer().RegisterMuxMiddleware(mn, hdl, true)
		}
		for mn, hdl := range middleware.BackendMiddleware {
			pa.Backend().ApiServer().RegisterMiddleware(mn, hdl, true)
		}
		for mn, hdl := range middleware.BackendMuxMiddleware {
			pa.Backend().ApiServer().RegisterMuxMiddleware(mn, hdl, true)
		}
	})
	// config
	if err := p.NewConfig(true); err != nil {
		return err
	}
	if p.Config().Independent() {
		p.WithoutRoutePrefix()
	}
	p.Config().Http.Frontend.ApiDisable = config.FrontendApiDisable
	p.Config().Http.Frontend.RpcDisable = config.FrontendRpcDisable
	p.Config().Http.Frontend.DocDisable = config.FrontendDocDisable
	p.Config().Http.Frontend.SocketHandler = &config2.SocketHandler{
		Module:    config.HandlerModule,
		SubModule: config.HandlerSubModule,
		Desc:      config.HandlerDesc,
		DocPublic: config.SocketDocPublic,
		TcpHandler: &config2.TcpHandler{
			Disable:    config.TcpDisable,
			DocDisable: config.TcpDisable,
		},
		WssHandler: &config2.WssHandler{
			Disable:    config.WssDisable,
			DocDisable: config.WssDisable,
		},
	}
	p.Config().Http.Backend.ApiDisable = config.BackendApiDisable
	p.Config().Http.Backend.RpcDisable = config.BackendRpcDisable
	p.Config().Http.Backend.DocDisable = config.BackendDocDisable
	initDoc(p)

	if err := initSocketGateway(p); err != nil {
		return err
	}
	// database
	if !config.MustDb {
		p.Config().Database = nil
	}
	// redis
	if !config.MustCache {
		p.Config().Cache = nil
	}
	// app
	if err := p.NewApp(); err != nil {
		return err
	}
	// cus service
	if err := cusservice.Exec(p); err != nil {
		return err
	}
	// register
	register.Exec(p)

	return p.Prepare()
}

func initDoc(p *application.Project) {
	p.Config().Http.Backend.Doc = nil
}

func initSocketGateway(p *application.Project) error {
	p.SetSocketGatewayCodecProvider(sockettype.TCP, config.CodecProvider())
	p.SetSocketGatewayCodecProvider(sockettype.TCP4, config.CodecProvider())
	p.SetSocketGatewayCodecProvider(sockettype.TCP6, config.CodecProvider())
	if len(p.Config().Http.Frontend.SocketGateways) == 0 {
		return errors.New("gateway tcp and rpc port is required")
	}
	for ch, gws := range p.Config().Http.Frontend.SocketGateways {
		if _, ok := channels[ch]; !ok {
			return errors.New("unsupport channel: " + ch)
		}
		for st, gw := range gws {
			if _, ok := channels[ch][st]; !ok {
				return errors.New("socket type[" + st.String() + "] not support")
			}
			if gw.Port <= 0 {
				return errors.New("invalid gateway[" + st.String() + "] port")
			}
			if st == ("mqtt") && gw.Mqtt != nil {
				gw.Mqtt.ClientTopic = &mqtt.QosTopic{
					Topic: "efly/devices/{device_sn}/actions/{action}",
					Qos:   2,
				}
				gw.Mqtt.ServerTopic = &mqtt.QosTopic{
					Topic: "efly/server/{device_sn}/actions/{action}",
					Qos:   2,
				}
			}
		}
	}

	return nil
}
