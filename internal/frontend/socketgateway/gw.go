package socketgateway

import (
	"github.com/obnahsgnaw/socketgateway"
	"github.com/obnahsgnaw/socketgateway/pkg/socket"
	"github.com/obnahsgnaw/socketgateway/pkg/socket/sockettype"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	register.Register(func(p *register.Provider) {
		p.Frontend().SocketGateway().TrustGet(sockettype.TCP).With(socketgateway.Watcher(func(c socket.Conn, msg string, l zapcore.Level, data ...zap.Field) {
			logWatcher(p.Frontend().SocketGateway().NewAccessLogger(), c, msg, l, data...)
		}))
	})
}

func logWatcher(lg *zap.Logger, c socket.Conn, msg string, l zapcore.Level, data ...zap.Field) {
	data = append([]zap.Field{zap.Int("fd", c.Fd())}, data...)
	lg.Log(l, msg, data...)
}
