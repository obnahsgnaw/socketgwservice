package fileasset

import (
	"errors"
	"github.com/obnahsgnaw/api/service/autheduser"
	"github.com/obnahsgnaw/application/pkg/utils"
	uploadv1 "github.com/obnahsgnaw/fileassetapi/gen/fileasset_backend_api/upload/v1"
	viewv1 "github.com/obnahsgnaw/fileassetapi/gen/fileasset_backend_api/view/v1"
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/pbhttp/pkg/cache"
	"github.com/obnahsgnaw/pbhttp/pkg/unique"
	"github.com/obnahsgnaw/socketgwservice/config"
	"github.com/obnahsgnaw/socketgwservice/internal/backend/rpc/fileasset"
	"strconv"
	"time"
)

type Server struct {
	cache         cache.Cache
	project       string
	module        string
	maxSize       uint32
	maxCount      uint32
	viewUrlTtl    uint32
	UploadUrlTtl  uint32
	types         map[string][]string
	typeMaps      map[string]map[string]struct{}
	contentTypes  []string
	extensions    []string
	multipart     bool
	sessionTtl    uint32
	enableChecker func() bool
}

type Config struct {
	SessionId    string
	MaxSize      uint32
	ContentTypes []string
	Extensions   []string
	Ttl          uint32
	Multipart    bool
	MaxCount     uint32
}

type FileForm struct {
	SessionId string
	Items     []*FileItem
}
type FileItem struct {
	UploadId string
	Name     string
}

type File struct {
	Name    string
	ViewUrl string
}

func New(c cache.Cache, project project.Service, module string, types map[string][]string, o ...Option) *Server {
	s := &Server{
		cache:        c,
		project:      project.Key(),
		module:       module,
		maxSize:      5 * 1024 * 1024,
		maxCount:     1,
		viewUrlTtl:   60 * 60,
		UploadUrlTtl: 60 * 10,
		types:        types,
		multipart:    false,
		sessionTtl:   60 * 10,
	}
	s.with(o...)
	s.init()
	return s
}

func (c *Server) with(o ...Option) {
	for _, oo := range o {
		if oo != nil {
			oo(c)
		}
	}
}

func (c *Server) Module() string {
	return c.module
}

func (c *Server) Config(user autheduser.User, rqId string) (resp *Config, err error) {
	if !c.enabled() {
		err = errors.New("file asset server not enabled")
		return
	}
	if err = validateQid(rqId); err != nil {
		return
	}
	resp = &Config{
		MaxSize:      c.maxSize,
		ContentTypes: c.contentTypes,
		Extensions:   c.extensions,
		Ttl:          c.sessionTtl,
		Multipart:    c.multipart,
		MaxCount:     c.maxCount,
	}
	if resp.SessionId, err = c.newSessionId(user, rqId); err != nil {
		return nil, errors.New("create session id failed")
	}

	return
}

func (c *Server) Url(user autheduser.User, sessionId, contentType, extension string, partNum uint32) (uploadId, name string, urls []string, err error) {
	if !c.enabled() {
		err = errors.New("file asset server not enabled")
		return
	}
	var urlResp *uploadv1.FetchUrlResponse
	if sessionId == "" {
		err = errors.New("session id not required")
		return
	}
	if contentType == "" {
		err = errors.New("content type required")
		return
	}
	if extension == "" {
		err = errors.New("extension required")
		return
	}
	if extensions, ok := c.typeMaps[contentType]; !ok {
		err = errors.New("content type not allowed")
		return
	} else {
		if _, ok = extensions[extension]; !ok {
			err = errors.New("extension not follow the content type")
			return
		}
	}
	if !c.multipart && partNum > 1 {
		err = errors.New("multipart not allowed")
		return
	}
	if err = c.ValidateSessionId(user, sessionId); err != nil {
		return
	}
	if urlResp, err = fileasset.GetUploadUrl(user, &uploadv1.FetchUrlRequest{
		Project:     c.project,
		Uid:         user.Id(),
		SessionId:   sessionId,
		Module:      c.module,
		MaxSize:     c.maxSize,
		ContentType: contentType,
		Extension:   extension,
		Ttl:         c.UploadUrlTtl,
		Part:        partNum,
		MaxCount:    c.maxCount,
	}); err != nil {
		err = errors.New("获取地址失败")
		return
	}
	uploadId = urlResp.UploadId
	name = urlResp.Name
	urls = urlResp.Url
	return
}

func (c *Server) ValidateOne(uploadFile *FileForm, savedFile string) (hit bool, file string, err error) {
	var files []string
	hit, files, err = c.Validate(uploadFile, c.ToSlice(savedFile))
	if len(files) > 0 {
		file = files[0]
	}
	return
}

func (c *Server) Validate(uploadFile *FileForm, savedFile []string) (hit bool, files []string, err error) {
	if uploadFile != nil {
		if uploadFile.SessionId == "" {
			err = errors.New("invalid session id")
			return
		}
		if len(uploadFile.Items) > int(c.maxCount) {
			err = errors.New("invalid file count, too many")
			return
		}
		if c.maxCount == 1 {
			if len(uploadFile.Items) > 0 && (len(savedFile) == 0 || uploadFile.Items[0].Name != savedFile[0]) {
				hit = true
			}
		} else {
			if !sliceEq(uploadFile.Items, savedFile) {
				hit = true
			}
		}
	}
	if hit {
		for _, i := range uploadFile.Items {
			files = append(files, i.Name)
		}
	}
	return
}

func (c *Server) Confirm(user autheduser.User, target string, file *FileForm) error {
	if !c.enabled() {
		return errors.New("file asset server not enabled")
	}
	rq := &uploadv1.ConfirmRequest{
		Project:   c.project,
		Uid:       user.Id(),
		SessionId: file.SessionId,
		Module:    c.module,
		Target:    target,
		Items:     nil,
	}
	for _, i := range file.Items {
		rq.Items = append(rq.Items, &uploadv1.FileItem{
			UploadId: i.UploadId,
			Name:     i.Name,
		})
	}
	return fileasset.Confirm(user, rq)
}

func (c *Server) ViewUrl(user autheduser.User, name string) (url2 File, err error) {
	if !c.enabled() {
		err = errors.New("file asset server not enabled")
		return
	}
	var resp *viewv1.ViewUrlResponse
	if resp, err = fileasset.GetViewUrl(user, &viewv1.ViewUrlRequest{
		Project: c.project,
		Module:  c.module,
		Name:    name,
		Ttl:     c.viewUrlTtl,
	}); err != nil {
		return
	}
	url2 = File{
		Name:    name,
		ViewUrl: resp.Url,
	}
	return
}

func (c *Server) ViewUrls(user autheduser.User, names []string) (urls map[string]File, err error) {
	if !c.enabled() {
		err = errors.New("file asset server not enabled")
		return
	}
	var resp *viewv1.ViewUrlsResponse
	urls = map[string]File{}
	if resp, err = fileasset.GetViewUrls(user, &viewv1.ViewUrlsRequest{
		Project: c.project,
		Module:  c.module,
		Names:   names,
		Ttl:     c.viewUrlTtl,
	}); err != nil {
		return
	}
	for name, url := range resp.Urls {
		urls[name] = File{
			Name:    name,
			ViewUrl: url,
		}
	}
	return
}

func (c *Server) enabled() bool {
	if c.enableChecker != nil {
		return c.enableChecker()
	}
	return false
}

func (c *Server) init() {
	extensions := make(map[string]struct{})
	c.typeMaps = make(map[string]map[string]struct{})
	for k, item := range c.types {
		c.contentTypes = append(c.contentTypes, k)
		c.typeMaps[k] = make(map[string]struct{})
		for _, k1 := range item {
			c.typeMaps[k][k1] = struct{}{}
			extensions[k1] = struct{}{}
		}
	}
	for k := range extensions {
		c.extensions = append(c.extensions, k)
	}
}

func (c *Server) newSessionId(user autheduser.User, rqId string) (sid string, err error) {
	rqCacheKey := c.cacheKey(user, rqId)
	// rqId => sid
	if err = c.cache.Remove(rqCacheKey); err != nil {
		return "", err
	}
	sid = rqId + unique.Unique()
	// cached rqId => sid
	if err = c.cache.Cache(rqCacheKey, sid, time.Second*time.Duration(c.sessionTtl)); err != nil {
		return "", err
	}
	return sid, nil
}

func (c *Server) ValidateSessionId(user autheduser.User, sid string) error {
	if err := validateSid(sid); err != nil {
		return err
	}
	rqId := sid[:10]

	cachedSid, _, err := c.cache.Cached(c.cacheKey(user, rqId))
	if err != nil {
		return err
	}
	if cachedSid == "" || sid != cachedSid {
		return errors.New("invalid session id")
	}

	return nil
}

func (c *Server) cacheKey(user autheduser.User, rqId string) string {
	return utils.ToStr("uploadSids:", config.Project.Key(), ":", c.module, ":", strconv.Itoa(int(user.Id())), ":", rqId)
}

func sliceEq(upload []*FileItem, saved []string) bool {
	// all is 0
	if len(upload) == 0 && len(saved) == 0 {
		return true
	}
	// len diff
	if len(upload) != len(saved) {
		return false
	}
	// len same
	savedMap := make(map[string]struct{})
	for _, i := range saved {
		savedMap[i] = struct{}{}
	}
	for _, i := range upload {
		if _, ok := savedMap[i.Name]; !ok {
			return false
		}
	}
	return true
}

func (c *Server) ToSlice(f string) []string {
	if f != "" {
		return []string{f}
	}
	return nil
}

// rqId len = 10
func validateQid(rqId string) error {
	if rqId == "" || len(rqId) != 10 {
		return errors.New("rqId invalid")
	}
	return nil
}

// sid len = 30
func validateSid(sid string) error {
	if sid == "" || len(sid) != 30 {
		return errors.New("sid invalid")
	}
	return nil
}
