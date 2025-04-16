package fileconfig

import (
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/config"
	"github.com/obnahsgnaw/socketgwservice/internal/service/fileasset"
)

func init() {
	register.Register(func(p *register.Provider) {
		p.Service().Register("upload-demo", func() interface{} {
			return fileasset.New(
				p.CacheProvider().Cache,
				config.Project,
				"demo",
				map[string][]string{
					"image/jpg":  {"jpg", "jpeg"},
					"image/jpeg": {"jpg", "jpeg"},
					"image/png":  {"png"},
					"image/gif":  {"gif"},
				},
				fileasset.MaxSize(5*1024*1024),
				fileasset.MaxCount(1),
				fileasset.UploadUrlTtl(5*60),
				fileasset.ViewUrlTtl(5*60),
				fileasset.SessionIdTtl(5*60),
				fileasset.Enable(func() bool {
					return p.CommonConfig().ServiceEnabled(project.FileService)
				}),
			)
		})
	})
}

func DemoService() *fileasset.Server {
	return register.Provide.Service().TrustGet("upload-demo").(*fileasset.Server)
}
