package repo

import (
	"errors"
	"github.com/obnahsgnaw/pbhttp"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/internal/service/queryutils"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type ColumnAble interface {
	GetFieldByName(fieldName string) (field.OrderExpr, bool)
}

type Repo struct {
	//
}

func (r *Repo) Db() *gorm.DB {
	return register.Provide.Database().Conn()
}

func (r *Repo) Rds() *pbhttp.CacheDriver {
	return register.Provide.CacheProvider()
}

func (r *Repo) ToDo(dao gen.Dao) gen.DO {
	return *dao.(*gen.DO)
}

func (r *Repo) InitPage(do gen.DO, page queryutils.Page) gen.DO {
	return InitPage(do, page)
}

func (r *Repo) DoSelect(do gen.DO, m ColumnAble, columns []string) (gen.DO, error) {
	if columns != nil && len(columns) > 0 {
		col, err := r.InitColumns(m, columns)
		if err != nil {
			return do, err
		}
		do = r.ToDo(do.Select(col...))
	}
	return do, nil
}

func (r *Repo) InitColumns(o ColumnAble, fields []string) (columns []field.Expr, err error) {
	if len(fields) == 0 {
		err = errors.New("repo parse column failed,err= not set select fields")
		return
	}
	for _, f := range fields {
		if c, ok := o.GetFieldByName(f); !ok {
			err = errors.New("repo parse column failed,err= column[" + f + "] not found")
			return
		} else {
			columns = append(columns, c)
		}
	}
	return
}

func ToDo(dao gen.Dao) gen.DO {
	return *dao.(*gen.DO)
}

func InitPage(do gen.DO, page queryutils.Page) gen.DO {
	if page.Id < 1 {
		page.Id = 1
	}
	do = ToDo(do.Offset(int((page.Id - 1) * page.Limit)))
	do = ToDo(do.Limit(int(page.Limit)))
	return do
}
