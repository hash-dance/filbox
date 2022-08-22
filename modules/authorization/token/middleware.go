// Package token Provide a token middleware
package token

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// Middleware check login status
// if not login, redirect to oauth2 server
func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := apicontext.ReadAPIContext(r.Context())
			tokenString := GetTokenAuthFromRequest(r)
			if tokenString == "" {
				render.SendError(w, r, render.Unauthorized, fmt.Errorf("未认证"))
				return
			}
			token, claims, err := ParseToken(tokenString)
			if err != nil || !token.Valid {
				render.SendError(w, r, render.Unauthorized, err)
				return
			}

			principal, err := models.GetPrincipalByPhone(claims.Phone)
			if err != nil {
				render.SendError(w, r, render.Unauthorized, err)
				return
			}
			// 认证成功
			ctx.Principal = principal
			// ctx.Token = token

			{ // 设置cookie
				secure := false
				if r.TLS != nil {
					secure = true
				}
				http.SetCookie(w, &http.Cookie{
					Name:     util.CookieName,
					Value:    tokenString,
					Secure:   secure,
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Duration(ctx.Config.SessionTimeOut) * time.Second),
				})
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
