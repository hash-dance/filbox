/*Package token create tokens
 */
package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/redis"
	"gitee.com/szxjyt/filbox-backend/models"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

const (
	tokenCharacters = "bcdfghjklmnpqrstvwxza245678901-"
	tokenLength     = 64
)

var tokenCharsLength = big.NewInt(int64(len(tokenCharacters)))

// Interface token handler
type Interface interface {
	CreateLoginToken(req *http.Request, second int64) (token *types.Token, err error)
	CreateLoginTokenForUser(req *http.Request, principal *models.Principal) (token *types.Token, err error)
	Authenticate(req *http.Request) (valid bool, principal *models.Principal, token *types.Token, err error)
}

type handler struct{}

// NewHandler return handler
func NewHandler() Interface {
	return &handler{}
}

func (tk *handler) CreateLoginTokenForUser(req *http.Request, principal *models.Principal) (*types.Token, error) {
	context := apicontext.ReadAPIContext(req.Context())
	expired := context.Config.SessionTimeOut
	token, err := tk.CreateLoginToken(req, expired)
	if err != nil {
		return nil, err
	}
	// token存入redis
	err = tk.save(token, principal, time.Second*time.Duration(expired))
	if err != nil {
		return nil, err
	}
	return token, nil
}

// CreateLoginToken return token
func (tk *handler) CreateLoginToken(req *http.Request, second int64) (*types.Token, error) {
	// create token
	token, err := tk.Generate()
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Authenticate 本地验证，验证Token是否有效
func (tk *handler) Authenticate(req *http.Request) (bool, *models.Principal, *types.Token, error) {
	// check sessionToken
	tokenString := GetTokenAuthFromRequest(req)
	if tokenString == "" {
		return false, nil, nil, fmt.Errorf("must authenticate")
	}

	token := &types.Token{Value: tokenString}
	principal, err := tk.get(token)
	if err != nil {
		return false, nil, nil, err
	}
	// token验证成功，刷新token时间
	context := apicontext.ReadAPIContext(req.Context())
	expired := context.Config.SessionTimeOut
	tk.updateExpired(token, time.Second*time.Duration(expired))
	return true, principal, token, nil
}

func (tk *handler) Generate() (*types.Token, error) {
	token := make([]byte, tokenLength)
	for i := range token {
		r, err := rand.Int(rand.Reader, tokenCharsLength)
		if err != nil {
			logrus.Debugf("generate token error: %s", err.Error())
			return nil, err
		}
		token[i] = tokenCharacters[r.Int64()]
	}
	return &types.Token{Value: string(token)}, nil
}

func (tk *handler) save(token *types.Token, principal *models.Principal, expiration time.Duration) error {
	client := redis.GetClient()
	// remove password before save to redis
	principal.Password = ""
	if err := client.SetObj(token.Value, principal, expiration); err != nil {
		return err
	}
	return nil
}

func (tk *handler) get(token *types.Token) (*models.Principal, error) {
	client := redis.GetClient()
	var read = models.Principal{}
	if err := client.GetObj(token.Value, &read); err != nil {
		return nil, err
	}
	principal, err := models.GetPrincipalByID(read.ID)
	if err != nil {
		logrus.Errorf("can not find user [%v], [%s]", read.ID, err.Error())
		return nil, err
	}
	return principal, nil
}

func (tk *handler) updateExpired(token *types.Token, expiration time.Duration) {
	if _, err := redis.GetClient().Cli().Expire(token.Value, expiration).Result(); err != nil {
		logrus.Errorf("update expired for [%s] error: [%s]", token.Value, err.Error())
	}
}

// GetTokenAuthFromRequest parse request, return token
func GetTokenAuthFromRequest(req *http.Request) string {
	var tokenAuthValue string
	authHeader := req.Header.Get(util.AuthHeaderName)
	authHeader = strings.TrimSpace(authHeader)

	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if strings.EqualFold(parts[0], util.AuthValuePrefix) {
			if len(parts) > 1 {
				tokenAuthValue = strings.TrimSpace(parts[1])
			}
		} else if strings.EqualFold(parts[0], util.BasicAuthPrefix) {
			if len(parts) > 1 {
				base64Value := strings.TrimSpace(parts[1])
				data, err := base64.URLEncoding.DecodeString(base64Value)
				if err != nil {
					logrus.Errorf("Error %v parsing %v header", err, util.AuthHeaderName)
				} else {
					tokenAuthValue = string(data)
				}
			}
		}
	} else {
		cookie, err := req.Cookie(util.CookieName)
		if err == nil {
			tokenAuthValue = cookie.Value
		}
	}
	return tokenAuthValue
}
