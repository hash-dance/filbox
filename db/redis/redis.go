/*Package redis init redis connection
 */
package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitee.com/szxjyt/filbox-backend/conf"
)

// Interface operate golang struct
type Interface interface {
	SetObj(key string, obj interface{}, expiration time.Duration) error
	GetObj(key string, obj interface{}) error
	Cli() *redis.Client
}

type client struct {
	Client *redis.Client
}

// var client *redis.Client
var c client

// GetClient return redis client
func GetClient() Interface {
	return &c
}

// SetupConnection init redis connection
func SetupConnection(config *conf.Config) {
	logrus.Infof("start connect redis addr=[%s], password=[%s], dbNumber=[%d]",
		config.RedisAddress, config.RedisPassword, config.RedisDBNumber)

	doConnect := func() error {
		newClient := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			Password: config.RedisPassword,
			DB:       config.RedisDBNumber,
		})

		_, err := newClient.Ping().Result()
		if err != nil {
			return err
		}
		c.Client = newClient
		return nil
	}
	// ticker := time.NewTicker(time.Second * 5)
	// for range ticker.C {
	// 	if err := doConnect(); err != nil {
	// 		logrus.Infof("can't connect redis addr=[%s], password=[%s], dbNumber=[%d], retry",
	// 			config.RedisAddress, config.RedisPassword, config.RedisDBNumber)
	// 		continue
	// 	}
	// 	ticker.Stop()
	// 	break
	// }

	go func() {
		var count int
	connDB:
		if count > 10 {
			panic("can not connect redis, panic")
		}
		if err := doConnect(); err != nil {
			logrus.Errorf("can't connect redis addr=[%s], password=[%s], dbNumber=[%d], retry",
				config.RedisAddress, config.RedisPassword, config.RedisDBNumber)
			time.Sleep(time.Second * 1)
			count++
			goto connDB
		}
		logrus.Infof("Redis Connection established")
	}()
}

func (c *client) SetObj(key string, obj interface{}, expiration time.Duration) error {
	b, err := json.Marshal(&obj)
	if err != nil {
		return err
	}
	v, err := c.Client.Set(key, string(b), expiration).Result()
	if err != nil {
		return err
	}
	logrus.Debugf("redis save success %s: ", v)
	return nil
}

func (c *client) GetObj(key string, obj interface{}) error {
	v, err := c.Client.Get(key).Result()
	if err != nil {
		return err
	}
	logrus.Debugf("get value from redis %s", v)
	err = json.Unmarshal([]byte(v), obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Cli() *redis.Client {
	return c.Client
}

// ArgInit args redis needed
func ArgInit(config *conf.Config) []cli.Flag {
	return []cli.Flag{
		// redis
		cli.StringFlag{
			Name:        "redis-addr",
			Usage:       "redis address ",
			EnvVar:      "REDIS_ADDR",
			Value:       "localhost:6379",
			Destination: &config.RedisAddress,
		}, cli.StringFlag{
			Name:        "redis-password",
			Usage:       "redis password ",
			EnvVar:      "REDIS_PASSWORD",
			Value:       "",
			Destination: &config.RedisPassword,
		}, cli.IntFlag{
			Name:        "redis-database",
			Usage:       "redis database number",
			EnvVar:      "REDIS_DATABASE",
			Value:       0,
			Destination: &config.RedisDBNumber,
		},
	}
}
