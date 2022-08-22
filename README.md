# filbox-backend
- 单点登录对接
- 用户认证
- 自定义路由
- 完整的日志模块
- 路由日志追踪
- 支持翻页搜索
- 支持请求参数校验

# 环境依赖
## golang
```
$ go version
go version go1.13 linux/amd64
```

## mariadb
`mariadb:10.1`

## redis 4.0.9

```text
$ redis-cli --version
redis-cli 4.0.9

$ redis-server --version
Redis server v=4.0.9 sha=00000000:0 malloc=jemalloc-3.6.0 bits=64 build=76095d16786fbcba
```

## 代码静态检测

代码检测使用[golangci-lint](https://github.com/golangci/golangci-lint)

```
$ golangci-lint version
golangci-lint has version 1.20.1 built from 849044b on 2019-10-15T19:12:01Z
```

在项目的配置文件`.golangci.yml`定义了检测规则


**有多种使用方法**

- [点击这里](https://github.com/golangci/golangci-lint#editor-integration)，将静态检测集成到编辑器，配置自动检测


- 如果编辑器不支持配置，在项目根路径执行`golangci-lint run`即可

- 调试`GL_DEBUG=linters_output GOPACKAGESPRINTGOLISTERRORS=1 golangci-lint run`


# 日志输出
默认输出目录是`./log`，日志名称`log`,日志保存七天,每小时分割

通过环境变量`$LOG_DISPATCH`启用日志分级
 - `info|debug`级别的日志存储在`info`文件
 - `panic|fatal|error|warn`级别的日志存储在`error`文件
 
通过环境变量`$LOGS_PATH`修改日志输出路径


# 二进制启动

## 初始化数据库
程序启动，会在配置的mysql数据库初始化表
程序启动完成后，手动执行数据库初始化脚本，数据库脚本在源代码的`resources/database/`目录


## 编译运行

运行前先配置参数

```text
$  go build -mod=vendor -o main main.go
$ ./main
```

## 执行
```text
NAME:
   app-server - server!

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug                   Enable debug logs [$DEBUG]
   --logs-format value       logs Format, can be 'json' or 'text' (default: "text") [$LOGS_FORMAT]
   --logs-Path value         logs output path (default: "./log") [$LOGS_PATH]
   --logDispatch             dispatch log to info,error log file by logLevel [$LOG_DISPATCH]
   --session-timeout value   session will invalid when timeout seconds (default: 3600) [$SESSION_TIMEOUT]
   --cors                    Enable cors server [$CORS]
   --redirect                open redirect 302, don't open it option unless you know [$REDIRECT]
   --http-listen-port value  Server Port (default: 80) [$HTTP_PORT]
   --monitor                 Enable Monitor Service [$MONITOR]
   --ssl                     Enable https server [$SSL]
   --ssl-crt value           ssl crt file path (default: "./conf/ssl/ssl.crt") [$SSL_CRT]
   --ssl-key value           ssl key file path (default: "./conf/ssl/ssl.key") [$SSL_KEY]
   --redis-addr value        redis address  (default: "localhost:6379") [$REDIS_ADDR]
   --redis-password value    redis password [$REDIS_PASSWORD]
   --redis-database value    redis database number (default: 0) [$REDIS_DATABASE]
   --mysql-addr value        mysql address (default: "127.0.0.1:33306") [$MYSQL_ADDR]
   --mysql-username value    mysql-username  (default: "root") [$MYSQL_USERNAME]
   --mysql-password value    mysql-password  (default: "password") [$MYSQL_PASSWORD]
   --mysql-database value    mysql database name (default: "raging_server") [$MYSQL_DATABASE]
   --help, -h                show help
   --version, -v             print the version

```

## docker 启动
启动参数通过`--env-file或者--env`配置

```text
$ docker run -d --name uaa -p 8001:80 --env-file env.file server:${TAG}
```

