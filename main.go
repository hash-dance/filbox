// Package main filbox-backend
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
package main

import (
	"context"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitee.com/szxjyt/filbox-backend/conf"
	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/deal"
	"gitee.com/szxjyt/filbox-backend/modules/util"
	"gitee.com/szxjyt/filbox-backend/routers"
)

var (
	// VERSION current version
	VERSION = "v0.1.0"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var config = conf.GetConfig()
	app := cli.NewApp()
	app.Name = "filbox-backend"
	app.Version = VERSION
	app.Usage = "filbox-backend!"

	app.Flags = initFlag(config) // 初始化日志配置

	ctx := util.SigTermCancelContext(context.Background())
	// 构建apicontext
	apiContext := apicontext.APIContext{
		Context: ctx,
		Config:  config,
	}

	app.Action = func(c *cli.Context) error {
		conf.InitLogs(config)         // 初始化日志配置
		mysql.SetupConnection(config) // 初始化数据库
		models.SyncDB()
		return run(&apiContext)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

// 初始化参数
func initFlag(config *conf.Config) []cli.Flag {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logs",
			EnvVar:      "DEBUG",
			Destination: &config.Debug,
		}, cli.StringFlag{
			Name:        "logs-format",
			Usage:       "logs Format, can be 'json' or 'text'",
			EnvVar:      "LOGS_FORMAT",
			Value:       conf.LogsText,
			Destination: &config.LogFormat,
		}, cli.StringFlag{
			Name:        "logs-Path",
			Usage:       "logs output path",
			EnvVar:      "LOGS_PATH",
			Value:       "./log",
			Destination: &config.LogPath,
		}, cli.BoolFlag{
			Name:        "logDispatch",
			Usage:       "dispatch log to info,error log file by logLevel",
			EnvVar:      "LOG_DISPATCH",
			Destination: &config.LogDispatch,
		},
		// session timeout seconds, 默认修改为一小时
		cli.Int64Flag{
			Name:        "session-timeout",
			Usage:       "session will invalid when timeout seconds",
			EnvVar:      "SESSION_TIMEOUT",
			Value:       60 * 60 * 60,
			Destination: &config.SessionTimeOut,
		},
		cli.StringFlag{
			Name:        "sms-region",
			Usage:       "sms region",
			EnvVar:      "SMS_REGION",
			Value:       "cn-hangzhou",
			Destination: &config.SmsKeyID,
		}, cli.StringFlag{
			Name:        "sms-keyid",
			Usage:       "sms keyid",
			EnvVar:      "SMS_KEYID",
			Value:       "Om",
			Destination: &config.SmsKeyID,
		}, cli.StringFlag{
			Name:        "sms-keySecret",
			Usage:       "sms keySecret",
			EnvVar:      "SMS_KEYSECRET",
			Value:       "Bqn",
			Destination: &config.SmsKeySecret,
		}, cli.StringFlag{
			Name:   "lotus-token",
			EnvVar: "LOTUS_TOEKN",
			Value:  "eysMRZQw_8WNxF-jYHYjX0IMy2Vp3mWvpPtg",
			Destination: &config.Lotus.Token,
		}, cli.StringFlag{
			Name:   "lotus-addr",
			EnvVar: "LOTUS_ADDR",
			Usage:  "example http://172.18.7.180:1234",
			Value:  "http://172.18.7.180:1234",
			Destination: &config.Lotus.Address,
		}, cli.StringFlag{
			Name:        "ipfs-token",
			EnvVar:      "IPFS_TOEKN",
			Value:       "",
			Destination: &config.Ipfs.Token,
		}, cli.StringFlag{
			Name:        "ipfs-addr",
			EnvVar:      "IPFS_ADDR",
			Usage:       "172.18.7.180:5001",
			Value:       "172.18.7.180:5001",
			Destination: &config.Ipfs.Address,
		}, cli.StringFlag{
			Name:        "tmp-path",
			EnvVar:      "TMP_PATH",
			Usage:       "tmp path save uploaded file",
			Value:       "/var/tmp",
			Destination: &config.TmpPath,
		},
	}
	for _, r := range [][]cli.Flag{
		basicFlag(config), mysql.ArgInit(config),
	} {
		flags = append(flags, r...)
	}
	return flags
}

func run(ctx *apicontext.APIContext) error {
	server := routers.NewServerConfig(ctx)
	server.Build()

	go deal.SyncDeals(ctx)
	go deal.MakeDeals(ctx)

	<-ctx.Context.Done()
	mysql.Destroy()

	return nil
}

func basicFlag(config *conf.Config) []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "cors",
			Usage:       "Enable cors server",
			EnvVar:      "CORS",
			Destination: &config.Cors,
		}, cli.BoolFlag{
			Name:        "redirect",
			Usage:       "open redirect 302, don't open it option unless you know",
			EnvVar:      "REDIRECT",
			Destination: &config.Redirect,
		}, cli.IntFlag{
			Name:        "http-listen-port",
			Usage:       "Server Port",
			EnvVar:      "HTTP_PORT",
			Value:       80,
			Destination: &config.HTTPListenPort,
		}, cli.BoolFlag{
			Name:        "monitor",
			Usage:       "Enable Monitor Service",
			EnvVar:      "MONITOR",
			Destination: &config.Monitor,
		},
		// ssl
		cli.BoolFlag{
			Name:        "ssl",
			Usage:       "Enable https server",
			EnvVar:      "SSL",
			Destination: &config.SSL,
		}, cli.StringFlag{
			Name:        "ssl-crt",
			Usage:       "ssl crt file path",
			EnvVar:      "SSL_CRT",
			Value:       "./conf/ssl/ssl.crt",
			Destination: &config.SSLCrtFile,
		}, cli.StringFlag{
			Name:        "ssl-key",
			Usage:       "ssl key file path",
			EnvVar:      "SSL_KEY",
			Value:       "./conf/ssl/ssl.key",
			Destination: &config.SSLKeyFile,
		},
	}
}
