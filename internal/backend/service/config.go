package service

import (
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	commonv1 "github.com/obnahsgnaw/socketgwserviceapi/gen/socketgw_backend_api/common/v1"
)

type PaginateConfig struct {
	keywordKinds  map[string]*commonv1.StringOption
	sortColumns   map[string]string
	optionColumns map[string]string
}

func NewPaginateConfig() *PaginateConfig {
	return &PaginateConfig{
		keywordKinds: map[string]*commonv1.StringOption{
			"keyword": {
				Id:   "keyword",
				Name: "关键字",
			},
		},
		sortColumns:   make(map[string]string),
		optionColumns: make(map[string]string),
	}
}

func (c *PaginateConfig) AddKeywordKind(key, title string) {
	c.keywordKinds[key] = &commonv1.StringOption{
		Id:   key,
		Name: title,
	}
}

func (c *PaginateConfig) AddSort(col, alias string) {
	if alias == "" {
		alias = col
	}
	c.sortColumns[col] = alias
}

func (c *PaginateConfig) AddOptionCol(col, optionKey string) {
	if optionKey == "" {
		optionKey = col
	}
	c.optionColumns[col] = optionKey
}

func (c *PaginateConfig) Get(page queryutils.Page) *commonv1.PaginateConfig {
	if page.Id > 1 {
		return nil
	}
	if len(c.keywordKinds) == 0 && len(c.sortColumns) == 0 && len(c.optionColumns) == 0 {
		return nil
	}
	var ko []*commonv1.StringOption
	for _, item := range c.keywordKinds {
		ko = append(ko, item)
	}
	return &commonv1.PaginateConfig{
		KeywordOptions: ko,
		SortColumns:    c.sortColumns,
		OptionColumns:  c.optionColumns,
	}
}
