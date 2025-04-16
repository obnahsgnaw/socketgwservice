package attrconfig

import (
	"errors"
	"github.com/obnahsgnaw/api/service/autheduser"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/socketgwservice/internal/backend/event"
	"github.com/obnahsgnaw/socketgwservice/internal/option"
	"github.com/obnahsgnaw/socketgwservice/internal/service"
	"github.com/obnahsgnaw/socketgwservice/internal/service/conflict"
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	"github.com/obnahsgnaw/socketgwservice/internal/service/transaction"
	"time"
)

type Config struct {
	Kind      AttrKind
	Attr      string
	Title     string
	Value     string
	Options   []string
	Enabled   option.EnableState
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Target struct {
	Id   uint32
	Name string
}

func NewTarget(id uint32, name string) Target {
	return Target{id, name}
}

type ConfigService struct {
	name  string
	trans TransactionRepo
	repo  ConfigRepo
}

type ConfigRepo interface {
	// ResourceTarget 资源名称 用于日志记录
	ResourceTarget() string
	// All 查询返回所有的配置
	All(target uint32, filter ConfigFilter, onlyEnabled bool) ([]*Config, int64, error)
	// Paginate 分页查询配置
	Paginate(target uint32, page queryutils.Page, filter ConfigFilter, onlyEnabled bool) ([]*Config, int64, error)
	// Exist 配置是否存在
	Exist(target uint32, key string) (bool, error)
	// GetOne 返回一个配置
	GetOne(target uint32, key string) (*Config, bool, error)
	// Create 添加一个配置
	Create(tx transaction.Tx, target uint32, config ...*Config) (uint32, error)
	// Update 更新一个配置
	Update(tx transaction.Tx, target uint32, config *Config) (uint32, error)
	// Remove 移除一个配置
	Remove(tx transaction.Tx, target uint32, key string) (uint32, error)
	// Clear 清空配置
	Clear(tx transaction.Tx, target uint32) error
}

type TransactionRepo interface {
	Transaction(transaction.Tx, func(tx transaction.Tx) error) error
}

type ConfigFilter struct {
	Keyword     string    `query:"keyword"`
	Attr        string    `query:"attr"`
	StartAt     time.Time `query:"start-at"`
	EndAt       time.Time `query:"end-at"`
	EnableState []int32   `query:"enable" options:"-1,0,1"`
	WithDeleted bool      `query:"-"`
}

func NewConfig(name string, trans TransactionRepo, repo ConfigRepo) *ConfigService {
	return &ConfigService{name: name, trans: trans, repo: repo}
}

// GetConfig 返回单个config
func (s *ConfigService) GetConfig(target Target, key string) (*Config, bool, error) {
	return s.repo.GetOne(target.Id, key)
}

// GetConfigs 返回多个config
func (s *ConfigService) GetConfigs(target Target, all bool, page queryutils.Page, filter ConfigFilter, onlyEnabled bool) ([]*Config, int64, error) {
	if all {
		return s.repo.All(target.Id, filter, onlyEnabled)
	}
	return s.repo.Paginate(target.Id, page, filter, onlyEnabled)
}

// RefreshConfig 刷新单个配置
func (s *ConfigService) RefreshConfig(operator autheduser.User, target Target, config *Config, conflictData string) (hit bool, conflicted bool, err error) {
	var id uint32
	var action string
	var desc string

	if err = s.validateConfig(config); err != nil {
		return
	}
	if exist, err1 := s.repo.Exist(target.Id, config.Attr); err1 != nil {
		err = err1
	} else {
		if exist {
			desc = "更新"
			action = "update"
			hit = true
			if !conflict.Validate(config.UpdatedAt, conflictData) {
				conflicted = true
			} else {
				id, err = s.repo.Update(transaction.Default(), target.Id, config)
			}
		} else {
			desc = "创建"
			action = "create"
			id, err = s.repo.Create(transaction.Default(), target.Id, config)
		}
	}
	// 触发操作事件
	event.NewOperateEvent(&event.OperateData{
		Operator:   operator,
		Target:     s.repo.ResourceTarget(),
		TargetId:   id,
		Action:     action,
		Content:    utils.ToStr(service.OperatorString(operator), desc, "了", s.name, target.Name, "的属性 ", config.Title, "[", config.Attr, "]"),
		DataBefore: nil,
		DataAfter:  config,
	}).Fire()
	return
}

// EnDisableConfig 启禁用单个配置
func (s *ConfigService) EnDisableConfig(operator autheduser.User, target Target, key string, enable bool, conflictData string) (config *Config, hit, conflicted bool, err error) {
	var dbf Config
	var action string
	var desc string
	var id uint32

	if config, hit, err = s.repo.GetOne(target.Id, key); err != nil {
		return
	} else {
		if !hit {
			return
		}
		if !conflict.Validate(config.UpdatedAt, conflictData) {
			conflicted = true
			return
		}
		dbf = *config
		if enable {
			action = "enable"
			desc = "启用"
			config.Enabled = option.Enabled
		} else {
			action = "disable"
			desc = "禁用"
			config.Enabled = option.Disabled
		}
		if id, err = s.repo.Update(transaction.Default(), target.Id, config); err != nil {
			return
		}
	}
	// 触发操作事件
	event.NewOperateEvent(&event.OperateData{
		Operator:   operator,
		Target:     s.repo.ResourceTarget(),
		TargetId:   id,
		Action:     action,
		Content:    utils.ToStr(service.OperatorString(operator), desc, "了", s.name, target.Name, "的属性 ", config.Title, "[", config.Attr, "]"),
		DataBefore: &dbf,
		DataAfter:  config,
	}).Fire()
	return
}

// RemoveConfig 删除单个配置
func (s *ConfigService) RemoveConfig(operator autheduser.User, target Target, key string) (hit bool, err error) {
	var id uint32
	if id, err = s.repo.Remove(transaction.Default(), target.Id, key); err != nil {
		return
	}
	hit = id > 0
	if hit {
		event.NewOperateEvent(&event.OperateData{
			Operator:   operator,
			Target:     s.repo.ResourceTarget(),
			TargetId:   id,
			Action:     "delete",
			Content:    utils.ToStr(service.OperatorString(operator), "删除了", s.name, target.Name, "的属性 ", "[", key, "]"),
			DataBefore: nil,
			DataAfter:  nil,
		}).Fire()
	}

	return
}

// RefreshConfigs 刷新所有配置
func (s *ConfigService) RefreshConfigs(operator autheduser.User, target Target, configs []*Config) (err error) {
	if len(configs) == 0 {
		return
	}
	for _, c := range configs {
		if err = s.validateConfig(c); err != nil {
			return
		}
	}
	if err = s.trans.Transaction(transaction.Default(), func(tx transaction.Tx) error {
		if err1 := s.repo.Clear(transaction.Default(), target.Id); err != nil {
			return err1
		}
		_, err1 := s.repo.Create(transaction.Default(), target.Id, configs...)
		return err1
	}); err != nil {
		return err
	}
	event.NewOperateEvent(&event.OperateData{
		Operator:   operator,
		Target:     s.repo.ResourceTarget(),
		TargetId:   target.Id,
		Action:     "refresh",
		Content:    utils.ToStr(service.OperatorString(operator), "批量替换了", s.name, target.Name, "的属性"),
		DataBefore: nil,
		DataAfter:  nil,
	}).Fire()

	return nil
}

// ClearConfigs 清空配置
func (s *ConfigService) ClearConfigs(operator autheduser.User, target Target) error {
	if err := s.repo.Clear(transaction.Default(), target.Id); err != nil {
		return err
	}
	event.NewOperateEvent(&event.OperateData{
		Operator:   operator,
		Target:     s.repo.ResourceTarget(),
		TargetId:   target.Id,
		Action:     "clear",
		Content:    utils.ToStr(service.OperatorString(operator), "清空了", s.name, target.Name, "的属性"),
		DataBefore: nil,
		DataAfter:  nil,
	}).Fire()

	return nil
}

func (s *ConfigService) validateConfig(config *Config) error {
	if config.Attr == "" || config.Title == "" {
		return errors.New("invalid arguments")
	}
	if v, err := validateValue(config.Kind, config.Value, config.Options); err != nil {
		return err
	} else {
		config.Value = v
		return nil
	}
}
