package searcher

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/modules/common/param"
)

// ParseSearch 解析模糊搜索参数
func ParseSearch(r *http.Request) (*mysql.QueryCondition, error) {
	// 获取search参数
	search := param.QueryString(r, "search")
	if search == "" {
		logrus.Warn("no search args")
		return nil, nil
	}
	logrus.Infof("want search [%s]", search)
	// 读取到模糊搜索结构
	fuzzyQuery := mysql.FuzzyQuery{}
	err := json.Unmarshal([]byte(search), &fuzzyQuery)
	if err != nil {
		logrus.Errorf("get search param, json.Unmarshal error [%s]", err.Error())
		return nil, err
	}
	// 解析成条件
	condition := buildSearchCondition(fuzzyQuery.Fields, fuzzyQuery.Key)
	return condition, nil
}

func buildSearchCondition(fields []string, pattern string) *mysql.QueryCondition {
	condition := mysql.QueryCondition{}
	for _, f := range fields {
		condition.Condition = append(condition.Condition, &mysql.Condition{
			Name:  f,
			Op:    "like",
			Value: pattern,
		})
	}
	return &condition
}
