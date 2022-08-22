// Package logout handler logout
package logout

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/redis"
	token2 "gitee.com/szxjyt/filbox-backend/modules/authorization/token"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// Router handler for logout
func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", logout)
	return r
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := token2.GetTokenAuthFromRequest(r)
	if _, err := redis.GetClient().Cli().Del(token).Result(); err != nil {
		logrus.Errorf("delete token err: [%s]", err.Error())
	}

	{ // 设置cookie
		secure := false
		if r.TLS != nil {
			secure = true
		}
		http.SetCookie(w, &http.Cookie{
			Name:     util.CookieName,
			Value:    "",
			Secure:   secure,
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
		})
	}
	http.Redirect(w, r, "/logout", http.StatusFound)
}
