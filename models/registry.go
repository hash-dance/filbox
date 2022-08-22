// Package models 定义model
// 同步到数据库的表
//
package models

import "gitee.com/szxjyt/filbox-backend/db/mysql"

// SyncDB models sync to database
func SyncDB() {
	mysql.GetClient().AutoMigrate([]interface{}{
		new(Principal), new(Wallet), new(File), new(DealInfo),
	}...)
	applyDefaultDatabase()
}

func applyDefaultDatabase() {}
