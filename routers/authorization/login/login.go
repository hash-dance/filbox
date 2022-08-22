// Package login handler local admin login
package login

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/authorization/token"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// Router handler for admin login
func Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/", apicontext.Bind(adminLogin, types.Input{}))
	r.Post("/code", apicontext.Bind(sendcode, types.PhoneNumber{}))
	return r
}

func sendcode(ctx *apicontext.APIContext, input *types.PhoneNumber) {
	if err := util.SendSms(input.Phone, ctx.Config.SmsRegion, ctx.Config.SmsKeyID, ctx.Config.SmsKeySecret); err != nil {
		ctx.Error(render.ServerError, err, "验证码发送失败")
		return
	}
	ctx.JSON(nil)
}

func adminLogin(ctx *apicontext.APIContext, input *types.Input) {
	// p, err := principal.ValidateLoginUser(input.Phone)
	// if err != nil {
	// 	logrus.Errorf("user [%s] use password [%s] login error", input.Username, input.Password)
	// 	ctx.Error(render.Unauthorized, fmt.Errorf("invalid username or password"), "login")
	// 	return
	// }
	// 检查code是否正确
	if !util.CodeIsEq(input.Phone, input.Code) {
		ctx.Error(render.Unauthorized, fmt.Errorf("验证码错误"), "error code")
		return
	}
	// 创建用户
	principal, err := models.GetPrincipalByPhone(input.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			principal = &models.Principal{
				Phone: input.Phone,
			}
			if err2 := models.CreatePrincipal(mysql.GetClient(), principal); err2 != nil {
				ctx.Error(render.Unauthorized, fmt.Errorf("创建用户失败"), "error create")
				return
			}
			goto TOKEN

		}
		ctx.Error(render.ServerError, err, "")
		return
	}

TOKEN:
	// 创建token
	tokestring, err := token.CreateToken(input.Phone)
	if err != nil {
		ctx.Error(render.Unauthorized, fmt.Errorf("创建token失败"), "error token")
		return
	}

	{ // 设置cookie
		secure := false
		if ctx.Req.TLS != nil {
			secure = true
		}
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     util.CookieName,
			Value:    tokestring,
			Secure:   secure,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(ctx.Config.SessionTimeOut) * time.Second),
		})
	}
	ctx.JSON(types.Token{tokestring})
}
