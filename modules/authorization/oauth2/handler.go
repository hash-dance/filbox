// Package oauth2 provider oauth method
package oauth2

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"gitee.com/szxjyt/filbox-backend/types"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/conf"
	"gitee.com/szxjyt/filbox-backend/db/redis"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// RedirectOauthServer redirect oauth2 server to login
func RedirectOauthServer(w http.ResponseWriter, r *http.Request, config *conf.Config) {
	u, err := url.Parse(config.OauthServer)
	if err != nil {
		render.SendError(w, r, render.ServerError, err)
		return
	}
	u.Path = path.Join(u.Path, config.OauthAPI)
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		render.SendError(w, r, render.ServerError, err)
		return
	}
	q := request.URL.Query()
	q.Add("client_id", config.ClientID)
	q.Add("redirect_uri", config.CallbackURL)
	q.Add("response_type", "code")
	// 生成state
	state := uuid.NewV4().String()
	// 存入redis，过期时间20s
	if err := redis.GetClient().SetObj(state, types.Redirect{
		URL: r.URL.String(),
	}, time.Second*60); err != nil {
		logrus.Infof("generate state error %s", err.Error())
		render.SendError(w, r, render.ServerError, err)
		return
	}
	q.Add("state", state)
	// request.URL.RawQuery = q.Encode() + http.Redirect("&scope=api+read_user"
	request.URL.RawQuery = q.Encode()
	if apicontext.ReadAPIContext(r.Context()).Config.Redirect {
		http.Redirect(w, r, request.URL.String(), http.StatusFound)
	} else {
		// 告诉前端路由，重定向到单点登录服务器去登录
		w.Header().Add("location", request.URL.String())
		render.SendError(w, r, render.RedirectError, fmt.Errorf("redirect: %s", request.URL.String()))
	}
}

// RequestAccessToken request an access_token using the code
func RequestAccessToken(code string, config *conf.Config) (accessToken string, err error) {
	u, err := url.Parse(config.OauthServer)
	if err != nil {
		logrus.Errorf("RequestAccessToken get url error [%s]", err.Error())
		return "", err
	}
	u.Path = path.Join(u.Path, config.TokenAPI)

	request, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		logrus.Errorf("RequestAccessToken create request client error [%s]", err.Error())
		return "", err
	}
	q := request.URL.Query()
	q.Add("client_id", config.ClientID)
	q.Add("redirect_uri", config.CallbackURL)
	q.Add("client_secret", config.Secret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	request.URL.RawQuery = q.Encode()

	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	resp, err := client.Do(request)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Errorf("close resp body err [%s]", err.Error())
		}
	}()
	if err != nil {
		logrus.Errorf("RequestAccessToken client.Do error [%s]", err.Error())
		return "", err
	}

	access := types.AccessToken{}
	if err := render.DecodeJSON(resp.Body, &access); err != nil {
		logrus.Errorf("RequestAccessToken DecodeJSON error [%s]", err.Error())
		return "", err
	}
	return access.AccessToken, nil
}

// GetUserInfo get userinfo from oauth2 provider
// https://docs.gitlab.com/ee/api/oauth2.html#access-gitlab-api-with-access-token
func GetUserInfo(accessToken string, config *conf.Config) (map[string]interface{}, error) {
	u, err := url.Parse(config.OauthServer)
	if err != nil {
		logrus.Errorf("RequestAccessToken get url error [%s]", err.Error())
		return nil, err
	}
	u.Path = path.Join(u.Path, config.UserAPI)

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		logrus.Infof("newRequest err: [%s]", err.Error())
		return nil, err
	}
	request.Header.Add(util.AuthHeaderName, "Bearer "+accessToken)

	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	resp, err := client.Do(request)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Errorf("close resp body err [%s]", err.Error())
		}
	}()
	if err != nil {
		logrus.Errorf("GetUserInfo client.Do error [%s]", err.Error())
		return nil, err
	}

	userInfo := make(map[string]interface{})
	if err := render.DecodeJSON(resp.Body, &userInfo); err != nil {
		logrus.Errorf("DecodeJSON userInfo error [%s]", err.Error())
		return nil, err
	}

	return userInfo, nil
}

// ValidateCode 通过code获取accessToken和用户信息
func ValidateCode(code string, config *conf.Config) (principal *models.Principal, err error) {
	// 获取token
	accesstoken, err := RequestAccessToken(code, config)
	if err != nil {
		logrus.Errorf("get access token error [%s]", err.Error())
		return nil, err
	}
	logrus.Infof("get accessToken: [%s]", accesstoken)
	// 获取用户信息
	userInfo, err := GetUserInfo(accesstoken, config)
	if err != nil {
		logrus.Errorf("get user info error: [%s]", err.Error())
		return nil, err
	}
	logrus.Infof("get user info [%v]", userInfo)
	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf("recover: [%+v]", e)
			err = fmt.Errorf("recover: [%+v]", e)
		}
	}()
	// 构建用户信息

	externalID, err := util.GetValueFromFields(strings.Split(config.ExternalID, "."), userInfo)
	if err != nil {
		return nil, err
	}
	username, err := util.GetValueFromFields(strings.Split(config.UserName, "."), userInfo)
	if err != nil {
		return nil, err
	}

	principal = &models.Principal{
		ExternalID: externalID,
		Username:   username,
		Role:       models.ROLEUSER,
	}
	return principal, nil
}
