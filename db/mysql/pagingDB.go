// nolint
package mysql

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

type PagingDB struct {
	*gorm.DB
}

func NewPagingDB() *PagingDB {
	return &PagingDB{ormDB.New()}
}

func (db *PagingDB) Paging(pagination *Pagination) *PagingDB {
	if pagination == nil {
		return db
	}
	return &PagingDB{db.Count(&pagination.Total).Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).Order(pagination.OrderBy)}
}

func (db *PagingDB) ConditionAND(condition []*Condition) *PagingDB {
	dbi := db.DB
	for _, v := range condition {
		switch v.Op {
		case "=":
			dbi = dbi.Where(v.Name+" = ? ", v.Value)
		case "!=":
			dbi = dbi.Where(v.Name+" != ? ", v.Value)
		case "like":
			dbi = dbi.Where(v.Name+" like ? ", fixLikeValue(v.Value))
		case "not like":
			dbi = dbi.Where(v.Name+" not like ? ", fixLikeValue(v.Value))
		}
	}
	return &PagingDB{dbi}
}

func (db *PagingDB) ConditionOR(condition []*Condition) *PagingDB {
	query := make([]string, 0)
	val := make([]interface{}, 0)
	for _, v := range condition {
		switch v.Op {
		case "=", "!=":
			query = append(query, strings.Join([]string{v.Name, v.Op, "?"}, " "))
			val = append(val, v.Value)
		case "like", "not like":
			query = append(query, strings.Join([]string{v.Name, v.Op, "?"}, " "))
			val = append(val, fixLikeValue(v.Value))
		}
	}
	return &PagingDB{db.DB.Where(strings.Join(query, " OR "), val...)}
}

func (db *PagingDB) Where(query interface{}, args ...interface{}) *PagingDB {
	return &PagingDB{db.DB.Where(query, args...)}
}

func (db *PagingDB) Or(query interface{}, args ...interface{}) *PagingDB {
	return &PagingDB{db.DB.Or(query, args...)}
}

func (db *PagingDB) Model(value interface{}) *PagingDB {
	return &PagingDB{db.DB.Model(value)}
}

func (db *PagingDB) Table(name string) *PagingDB {
	return &PagingDB{db.DB.Table(name)}
}

func (db *PagingDB) Select(query interface{}, args ...interface{}) *PagingDB {
	return &PagingDB{db.DB.Select(query, args...)}
}

func (db *PagingDB) Group(query string) *PagingDB {
	return &PagingDB{db.DB.Group(query)}
}

func (db *PagingDB) Joins(query string, args ...interface{}) *PagingDB {
	return &PagingDB{db.DB.Joins(query, args...)}
}

func fixLikeValue(val interface{}) string {
	v, ok := val.(string)
	if ok {
		re := regexp.MustCompile(`^(%)*|(%)*$`)
		return string(re.ReplaceAll([]byte(v), []byte("%")))
	}
	return "%"
}
