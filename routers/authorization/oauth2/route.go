// Package oauth2 handle oauth2 login
package oauth2

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/go-chi/chi"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/db/redis"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/authorization/oauth2"
	"gitee.com/szxjyt/filbox-backend/modules/authorization/token"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// Router handler for logincallback
func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", apicontext.Bind(handleCallback))
	return r
}

func handleCallback(ctx *apicontext.APIContext) {
	code := ctx.Req.URL.Query().Get("code")
	if code == "" {
		ctx.Error(render.OauthAuthError, fmt.Errorf("get code null"), "handleCallback")
		return
	}

	state := ctx.Req.URL.Query().Get("state")
	if state == "" {
		ctx.Error(render.OauthAuthError, fmt.Errorf("get state null"), "handleCallback")
		return
	}

	redirect := types.Redirect{}
	if err := redis.GetClient().GetObj(state, &redirect); err != nil {
		ctx.Error(render.OauthAuthError, fmt.Errorf("get state from redis null, maybe Man-in-the-middle attack [%s]", err.Error()), "handleCallback")
		return
	}
	// 使用code去oauth服务器获取accessToken和用户信息
	principal, err := oauth2.ValidateCode(code, ctx.Config)
	if err != nil {
		ctx.Error(render.OauthAuthError, err, "handleCallback")
		return
	}

	// 更新到数据库
	if principal, err = models.CreateOrUpdatePrincipalByExternalID(mysql.GetClient(), principal); err != nil {
		ctx.Error(render.OauthAuthError, err, "handleCallback.syncUser")
		return
	}
	// 颁发本地token
	tkHandler := token.NewHandler()
	tk, err := tkHandler.CreateLoginTokenForUser(ctx.Req, principal)
	if err != nil {
		ctx.Error(render.Unauthorized, err, "handleCallback.createToken")
		return
	}

	{ // 设置cookie
		secure := false
		if ctx.Req.TLS != nil {
			secure = true
		}
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     util.CookieName,
			Value:    tk.Value,
			Secure:   secure,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(ctx.Config.SessionTimeOut) * time.Second),
		})
	}
	// oauth login to where?
	// 这个会被单点登录通过浏览器地址回调过来，可以使用redirect控制浏览器跳转到首页
	ctx.Redirect("/", http.StatusFound)
}
