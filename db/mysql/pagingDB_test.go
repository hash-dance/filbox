// nolint
package mysql

import (
	"testing"
	"time"

	"gitee.com/szxjyt/filbox-backend/conf"
)

func Test_count(t *testing.T) {
	SetupConnection(&conf.Config{
		MysqlDatabase: "raging_server",
		MysqlAddress:  "127.0.0.1:33306",
		MysqlPassword: "password",
		MysqlUserName: "root",
	})
	time.Sleep(time.Second * 2)
	var count = 0
	users := make([]interface{}, 0)
	NewPagingDB().Table("users").
		Where("role = ?", "0").
		Count(&count).
		Limit(5).
		Offset(-1).
		Find(&users).Debug()
	t.Logf("count: %d", count)
}

func Test_switch(t *testing.T) {
	arr := []int{1, 2, 3}
	for _, v := range arr {
		switch v {
		case 1, 2, 3:
			t.Log(1)
		case 4:
			t.Log(4)
		default:
			t.Log("default")
		}
	}

}
