// Package searcher parse conditions and pagination
package searcher

import (
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/modules/common/param"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// ParsePagination 解析翻页参数
func ParsePagination(r *http.Request) (*mysql.Pagination, error) {
	page, err := param.QueryStringInt(r, "page")
	if err != nil {
		logrus.Warn("no page args, set default page 1")
		page = 1
	}
	pageSize, err := param.QueryStringInt(r, "pageSize")
	if err != nil {
		logrus.Warn("no pageSize args, set default pageSize 0")
		pageSize = 0
	}
	pagination := mysql.Pagination{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  param.QueryString(r, "order"),
		Total:    0,
	}
	fixPagination(&pagination)
	return &pagination, nil
}

func fixPagination(option *mysql.Pagination) {
	if option.PageSize <= 0 { // 默认每页条数
		option.PageSize = util.DefaultPageSize
	}

	if option.Page == -1 { // page为-1，获取全部数据
		option.Page = -1
		option.PageSize = -1
	} else if option.Page < 1 { // page < 1, 设置为1
		option.Page = 1
	}
	orderStr := option.OrderBy
	if orderStr != "" {
		option.OrderBy = parseOrder(orderStr)
	} else {
		option.OrderBy = util.DefaultOrder
	}
}

var orderReg = regexp.MustCompile(`(\+|-)([^,]+)`)

func parseOrder(order string) string {
	return orderReg.ReplaceAllStringFunc(order, func(a string) string {
		direction := ""
		if a[0] == '-' {
			direction = " desc"
		}
		return a[1:] + direction
	})
}
