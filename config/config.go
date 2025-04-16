package config

import (
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/api"
	"github.com/obnahsgnaw/application/pkg/utils"
	config2 "github.com/obnahsgnaw/pbhttp/core/config"
	"github.com/obnahsgnaw/pbhttp/core/project"
	indexv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/index/v1"
	"strings"
)

const (
	Project = project.TcpGatewayService

	Version api.Version = 1

	AdminRoute         = "/admin"
	MustDb             = false
	MustCache          = false
	PageLimit          = 10
	PageLimitMax       = 100
	SocketDocPublic    = false
	HandlerModule      = "" // 此三项开启handler
	HandlerSubModule   = ""
	HandlerDesc        = ""
	TcpDisable         = false // 单独禁用此类型
	WssDisable         = false // 单独禁用此类型
	FrontendApiDisable = false
	FrontendRpcDisable = false
	FrontendDocDisable = true
	BackendApiDisable  = false
	BackendRpcDisable  = false
	BackendDocDisable  = true
)

type Config struct {
	// 自定义 启动配置
}

func NewConfig(c *config2.Config) *indexv1.Config {
	return &indexv1.Config{Debug: c.Application.Debug}
}

func VersionRoute(route string) string {
	return utils.ToStr("/", Version.String(), "/", strings.Trim(route, "/"))
}

func VersionProjectRoute(route string) string {
	return utils.ToStr("/", Version.String(), "/", Project.Key(), "/", strings.Trim(route, "/"))
}

func FrontendAppIgnorer(c *gin.Context) bool {
	if len(frontendIgnorer) > 0 {
		for _, ig := range frontendIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	if len(frontendAppIgnorer) > 0 {
		for _, ig := range frontendAppIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	return false
}

func FrontendAuthIgnorer(c *gin.Context) bool {
	if len(frontendIgnorer) > 0 {
		for _, ig := range frontendIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	if len(frontendAuthIgnorer) > 0 {
		for _, ig := range frontendAuthIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	return false
}

func FrontendPermIgnorer(method, uriPattern string) bool {
	if len(frontendIgnorer) > 0 {
		for _, ig := range frontendIgnorer {
			if ig(method, uriPattern) {
				return true
			}
		}
	}
	if len(frontendPermIgnorer) > 0 {
		for _, ig := range frontendPermIgnorer {
			if ig(method, uriPattern) {
				return true
			}
		}
	}
	return false
}

func BackendAppIgnorer(c *gin.Context) bool {
	if len(backendIgnorer) > 0 {
		for _, ig := range backendIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	if len(backendAppIgnorer) > 0 {
		for _, ig := range backendAppIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	return false
}

func BackendAuthIgnorer(c *gin.Context) bool {
	if len(backendIgnorer) > 0 {
		for _, ig := range backendIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	if len(backendAuthIgnorer) > 0 {
		for _, ig := range backendAuthIgnorer {
			if ig(c.Request.Method, c.Request.URL.Path) {
				return true
			}
		}
	}
	return false
}

func BackendPermIgnorer(method, uriPattern string) bool {
	if len(backendIgnorer) > 0 {
		for _, ig := range backendIgnorer {
			if ig(method, uriPattern) {
				return true
			}
		}
	}
	if len(backendPermIgnorer) > 0 {
		for _, ig := range backendPermIgnorer {
			if ig(method, uriPattern) {
				return true
			}
		}
	}
	return false
}
