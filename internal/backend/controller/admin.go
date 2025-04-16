package controller

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/pbhttp/pkg/httpproxy"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/asset"
	"github.com/obnahsgnaw/socketgwservice/config"
	"net/http"
	"net/http/httputil"
	"os"
)

func init() {
	register.Register(func(p *register.Provider) {
		if config.AdminRoute == "" {
			return
		}
		adminRawPath := utils.ToStr("/", p.Project().Key(), "-admin")

		registerAdminRawRoute(p, adminRawPath)

		registerAdminRoute(p, adminRawPath)
	})
}

func registerAdminRawRoute(p *register.Provider, adminRawPath string) {
	p.RegisterBackendRawRouter(func(engine *gin.Engine) {
		engine.StaticFS(adminRawPath, &assetfs.AssetFS{
			Asset:    asset.Asset,
			AssetDir: asset.AssetDir,
			AssetInfo: func(path string) (os.FileInfo, error) {
				return os.Stat(path)
			},
			Prefix:   "html",
			Fallback: "",
		})
	})
}

func registerAdminRoute(p *register.Provider, adminRawPath string) {
	path := config.VersionRoute(config.AdminRoute)
	proxy, err := httpproxy.New(url.Origin{
		Protocol: url.HTTP,
		Host: url.Host{
			Ip:   p.Config().Application.InternalIp,
			Port: p.Config().Http.Backend.Port,
		},
	}.String(), path, adminRawPath)
	if err != nil {
		panic(err)
	}

	p.RegisterBackendStaticRouter("GET", path, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		file, _ := pathParams["path"]
		handler(proxy, w, r, path+"/"+file)
	})
	p.RegisterBackendStaticRouter("HEAD", path, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		file, _ := pathParams["path"]
		handler(proxy, w, r, path+"/"+file)
	})
}

func handler(proxy *httputil.ReverseProxy, w http.ResponseWriter, r *http.Request, uri string) {
	r.RequestURI = uri
	r.URL.Path = uri
	if r.URL.RawPath != "" {
		r.URL.RawPath = uri
	}
	proxy.ServeHTTP(w, r)
}
