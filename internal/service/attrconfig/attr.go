package attrconfig

import (
	"errors"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	"github.com/obnahsgnaw/socketgwservice/internal/service/transaction"
	"strconv"
	"time"
)

type AttrService struct {
	name string
	ttl  time.Duration
	c    []*ConfigService
	repo AttrRepo
}

type AttrRepo interface {
	All(target uint32) (map[string]string, error)
	Clear(tx transaction.Tx, target uint32) error
	Save(tx transaction.Tx, target uint32, attrs map[string]string) error
	Cache(key string, data string, ttl time.Duration) error
	Cached(key string) (string, bool, error)
}

func NewAttrService(name string, cacheTtl time.Duration, repo AttrRepo, c ...*ConfigService) *AttrService {
	if cacheTtl < time.Second {
		cacheTtl = time.Second * 10
	}
	return &AttrService{name: name, ttl: cacheTtl, c: c, repo: repo}
}

// GetAll 查询返回所有的属性
func (s *AttrService) GetAll(target Target) (map[string]string, error) {
	return s.repo.All(target.Id)
}

// Refresh 刷新属性, 即删除就属性，添加新属性
func (s *AttrService) Refresh(tx transaction.Tx, target Target, attrs map[string]string) (err error) {
	if err = s.repo.Clear(tx, target.Id); err != nil {
		return
	}
	if len(attrs) > 0 {
		return s.repo.Save(tx, target.Id, attrs)
	}
	return nil
}

// Validate 验证数据
func (s *AttrService) Validate(configTarget []Target, attrs map[string]string) (map[string]string, error) {
	if attrs1, err := s.validateAttr(configTarget, attrs); err != nil {
		return nil, err
	} else {
		return attrs1, nil
	}
}

func (s *AttrService) validateAttr(configTarget []Target, attrs map[string]string) (map[string]string, error) {
	if attrs == nil {
		attrs = make(map[string]string)
	}
	attrs1 := make(map[string]string)
	for i1, c1 := range s.c {
		if len(configTarget) <= i1 {
			continue
		}
		target := configTarget[i1]
		if c2, err := s.getConfigAttrs(c1, target); err != nil {
			return nil, err
		} else {
			if len(c2) > 0 {
				for _, c3 := range c2 {
					if v, ok := attrs[c3.Attr]; ok {
						if v != "" {
							attrs1[c3.Attr] = v
							continue
						}
						if c3.Value != "" {
							attrs1[c3.Attr] = c3.Value
							continue
						}
					} else {
						if c3.Value != "" {
							attrs1[c3.Attr] = c3.Value
							continue
						}
					}
					return nil, errors.New(c3.Title + "必须")
				}
			}
		}
	}
	return attrs1, nil
}

func (s *AttrService) getConfigAttrs(service *ConfigService, configTarget Target) ([]*Config, error) {
	key := utils.ToStr(s.name, ":attr-configs:", strconv.Itoa(int(configTarget.Id)))
	if rs, ok, err := s.repo.Cached(key); err != nil {
		return nil, err
	} else {
		data := make(map[string][]*Config)
		if ok {
			if ok = utils.ParseJson([]byte(rs), &data); ok {
				if v, ok1 := data["data"]; ok1 {
					return v, nil
				}
			}
		}
		if list, _, err1 := service.GetConfigs(configTarget, true, queryutils.Page{}, ConfigFilter{}, true); err1 != nil {
			return nil, err
		} else {
			data["data"] = list
			val := utils.ToJson(data)
			if val == "" {
				return nil, errors.New("attr.getConfigAttrs: encode to json failed")
			}
			if err = s.repo.Cache(key, val, s.ttl); err != nil {
				return nil, err
			}
			return list, nil
		}
	}
}
