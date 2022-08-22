/*Package mysql init mysql connection
 */
package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // init mysql
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitee.com/szxjyt/filbox-backend/conf"
)

var ormDB *gorm.DB

// GetClient return db
func GetClient() *gorm.DB {
	return ormDB.New()
}

// SetupConnection setup connection
func SetupConnection(config *conf.Config) {
	mysqlURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.MysqlUserName, config.MysqlPassword, config.MysqlAddress, config.MysqlDatabase)
	logrus.Infof("start connect mysql %s", mysqlURL)

	var count int
connDB:
	if count > 10 {
		panic("can not connect mysql, panic")
	}

	db, err := gorm.Open("mysql", mysqlURL)
	// db.SingularTable(true) // table name in singular
	db.DB().SetMaxOpenConns(40)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetConnMaxLifetime(time.Duration(10*60) * time.Second)
	ormDB = db
	if err != nil {
		logrus.Errorf("connect mysql error: %s", err.Error())
		time.Sleep(time.Second * 1)
		count++
		goto connDB
	}
	if config.Debug {
		ormDB.SetLogger(&gormLogger{})
		ormDB = ormDB.Debug()
	}
	logrus.Infof("Mysql Connection established")
}

// Destroy close db connection
func Destroy() {
	if err := ormDB.Close(); err != nil {
		logrus.Errorf("close mysql db error: %s", err.Error())
	} else {
		logrus.Infof("disconnect mysql success")
	}
}

// GormLogger struct
type gormLogger struct{}

// Print - Log Formatter
func (*gormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(
			logrus.Fields{
				"module":        "gorm",
				"type":          "sql",
				"rows_returned": v[5],
				"src":           v[1],
				"values":        v[4],
				"duration":      v[2],
			},
		).Infof("%s => %v", v[3], v[4])
	case "log":
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

// ArgInit args redis needed
func ArgInit(config *conf.Config) []cli.Flag {
	return []cli.Flag{
		// mysql
		cli.StringFlag{
			Name:        "mysql-addr",
			Usage:       "mysql address",
			EnvVar:      "MYSQL_ADDR",
			Value:       "127.0.0.1:33306",
			Destination: &config.MysqlAddress,
		}, cli.StringFlag{
			Name:        "mysql-username",
			Usage:       "mysql-username ",
			EnvVar:      "MYSQL_USERNAME",
			Value:       "root",
			Destination: &config.MysqlUserName,
		}, cli.StringFlag{
			Name:        "mysql-password",
			Usage:       "mysql-password ",
			EnvVar:      "MYSQL_PASSWORD",
			Value:       "password",
			Destination: &config.MysqlPassword,
		}, cli.StringFlag{
			Name:        "mysql-database",
			Usage:       "mysql database name",
			EnvVar:      "MYSQL_DATABASE",
			Value:       "filbox",
			Destination: &config.MysqlDatabase,
		},
	}
}
