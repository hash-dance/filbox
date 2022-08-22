package searcher

import (
	"net/http"

	"gitee.com/szxjyt/filbox-backend/modules/util"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
)

// Builder 根据request构建一个查询条件结构体
func Builder(r *http.Request) (*mysql.PageQueryInfo, error) {
	conditionKey := mysql.Query
	// 解析查询搜索参数`query`
	condition, err := ParseCondition(r)
	if err != nil {
		return nil, err
	}
	logrus.Infof("%+v", condition)

	// 解析模糊搜索字段`search`
	c, err := ParseSearch(r)
	if err != nil {
		return nil, err
	}
	if c != nil { // 模糊搜索生效
		condition = c
		conditionKey = mysql.Search
	}

	// 解析翻页参数
	pagination, err := ParsePagination(r)
	if err != nil {
		return nil, err
	}

	return &mysql.PageQueryInfo{
		Pagination:     pagination,
		QueryCondition: condition,
		ConditionKey:   conditionKey,
	}, nil
}

// DefaultQueryInfo default paging
func DefaultQueryInfo() *mysql.PageQueryInfo {
	return &mysql.PageQueryInfo{
		Pagination: &mysql.Pagination{
			Page:     util.DefaultPageSize,
			PageSize: util.DefaultPageSize,
			OrderBy:  util.DefaultOrder,
			Total:    0,
		},
		QueryCondition: nil,
		ConditionKey:   "",
	}
}
