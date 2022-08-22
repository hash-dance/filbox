package searcher

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/modules/common/param"
)

// ParseCondition 解析自定义条件参数
func ParseCondition(r *http.Request) (*mysql.QueryCondition, error) {
	condition := mysql.QueryCondition{}
	// 获取query参数
	query := param.QueryString(r, "query")
	if query == "" {
		logrus.Warn("no query args")
		return &condition, nil
	}
	err := json.Unmarshal([]byte(query), &condition)
	if err != nil {
		logrus.Errorf("get query param, json.Unmarshal error [%s]", err.Error())
		return nil, err
	}

	return &condition, nil
}
