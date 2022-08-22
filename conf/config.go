/*Package conf defined configuration
 */
package conf

// Config global config
type Config struct {
	Debug    bool // enable debug
	Cors     bool // enable cors http
	Redirect bool

	LogFormat   string // logsFormat
	LogPath     string // log output path
	LogDispatch bool   // dispatch to different file

	HTTPListenPort int  // server port
	Monitor        bool // 是否开启监控

	RedisAddress  string // redis address
	RedisPassword string // redis password
	RedisDBNumber int    // redis database number

	MysqlAddress  string
	MysqlUserName string
	MysqlPassword string
	MysqlDatabase string

	SessionTimeOut int64 // session timeout seconds

	SSL        bool
	SSLCrtFile string
	SSLKeyFile string

	// oauth2 服务配置参数
	ClientID    string
	Secret      string
	CallbackURL string
	OauthServer string
	OauthAPI    string // `/oauth/authorize`
	TokenAPI    string // `/oauth/token`
	UserAPI     string // `/api/v4/user`
	ExternalID  string // 单点登录系统的用户识别ID，UserAPI返回的字段
	UserName    string // 单点登录系统的用户名，UserAPI返回的字段

	SmsRegion    string // 短信服务
	SmsKeyID     string
	SmsKeySecret string

	Lotus struct {
		Token   string
		Address string
	}
	Ipfs struct {
		Token   string
		Address string
	}

	TmpPath string
}

var config *Config

// GetConfig return config
func GetConfig() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

// SetOauthConfig set oauth config
func (c *Config) SetOauthConfig(clientID, secret, callbackURL, oauthServer, oauthAPI, tokenAPI, userAPI, externalID, userName string) {
	c.ClientID = clientID
	c.Secret = secret
	c.CallbackURL = callbackURL
	c.OauthServer = oauthServer
	c.OauthAPI = oauthAPI
	c.TokenAPI = tokenAPI
	c.UserAPI = userAPI
	c.ExternalID = externalID
	c.UserName = userName
}
