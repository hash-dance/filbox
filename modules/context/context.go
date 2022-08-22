// Package apicontext defined common context
package apicontext

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/modules/common/searcher"
	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"gitee.com/szxjyt/filbox-backend/conf"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
)

// APIContext struct
type APIContext struct {
	Context   context.Context
	Config    *conf.Config
	Req       *http.Request
	Writer    http.ResponseWriter
	Validate  *validator.Validate  // validator
	QueryInfo *mysql.PageQueryInfo // 翻页查询参数

	Principal *models.Principal // 用户信息
	Token     *types.Token      // token
}

func (ctx *APIContext) Redirect(url string, code int) {
	http.Redirect(ctx.Writer, ctx.Req, url, code)
}

func (ctx *APIContext) JSON(data interface{}) {
	render.SendJSON(ctx.Writer, ctx.Req, data)
}

func (ctx *APIContext) JSONPagination(data interface{}) {
	render.SendPaginationJSON(ctx.Writer, ctx.Req, ctx.QueryInfo.Pagination, data)
}

func (ctx *APIContext) Error(code render.ErrorCode, err error, message string) {
	e := errors.Wrap(err, message)
	logrus.Error(e.Error())
	render.SendError(ctx.Writer, ctx.Req, code, e)
}

func (ctx *APIContext) Errorf(code render.ErrorCode, err error, format string, args ...interface{}) {
	e := errors.Wrapf(err, format, args)
	logrus.Error(e.Error())
	render.SendError(ctx.Writer, ctx.Req, code, e)
}

type contextKey struct {
	name string
}

// Middleware load common context
func Middleware(apiContext *APIContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey{"api-context"}, &APIContext{
				Context:  apiContext.Context,
				Config:   apiContext.Config,
				Validate: NewValidator(),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// ReadAPIContext load apicontext from context
func ReadAPIContext(ctx context.Context) *APIContext {
	return ctx.Value(contextKey{"api-context"}).(*APIContext)
}

// Bind read input and valida input fields
func Bind(handler interface{}, input ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ReadAPIContext(r.Context())
		ctx.Writer = w
		ctx.Req = r
		defer func(ctx *APIContext) {
			if err := recover(); err != nil {
				ctx.Error(render.ServerError, fmt.Errorf("bind error: [%s]", err), "defer")
			}
		}(ctx)

		contentType := ctx.Req.Header.Get("Content-Type")
		if ctx.Req.Method == "POST" || ctx.Req.Method == "PUT" || len(contentType) > 0 {
			if len(input) == 0 { // no request body input
				fn := reflect.ValueOf(handler)
				fn.Call([]reflect.Value{reflect.ValueOf(ctx)})
				// ctx.Error(render.InvalidData, fmt.Errorf("must body"), "apiContext bind")
				return
			}
			obj := input[0]
			typ := reflect.TypeOf(obj)
			data := reflect.New(typ).Interface()
			// todo parse contentType
			// 		switch {
			// 		case strings.Contains(contentType, "form-urlencoded"):
			// 		case strings.Contains(contentType, "multipart/form-data"):
			// 		case strings.Contains(contentType, "json"):
			// 		default:
			// 			err := fmt.Errorf("bind parse error")
			// 			logrus.Error(err.Error())
			// 		}
			if err := render.DecodeJSON(r.Body, data); err != nil {
				ctx.Error(render.InvalidData, err, "decode json")
				return
			}
			// validate input data
			err := ctx.Validate.Struct(data)
			if err != nil {
				ctx.Error(render.InvalidData, err, "validate input")
				return
			}

			fn := reflect.ValueOf(handler)
			fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(data)})

		} else { // get request
			if queryInfo, err := searcher.Builder(ctx.Req); err != nil {
				logrus.Errorf("parse queryInfo error: [%s]", err.Error())
				ctx.QueryInfo = searcher.DefaultQueryInfo()
			} else {
				ctx.QueryInfo = queryInfo
			}

			fn := reflect.ValueOf(handler)
			fn.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
	}
}

// func (ctx *APIContext)invoke()  {
// 	contentType := ctx.Req.Header.Get("Content-Type")
// 	if ctx.Req.Method == "POST" || ctx.Req.Method == "PUT" || len(contentType) > 0 {
// 		switch {
// 		case strings.Contains(contentType, "form-urlencoded"):
// 		case strings.Contains(contentType, "multipart/form-data"):
// 		case strings.Contains(contentType, "json"):
// 		default:
// 			err := fmt.Errorf("bind parse error")
// 			logrus.Error(err.Error())
// 		}
// 	} else {
//
// 	}
// }

// RequireSystemRole check system role
func RequireSystemRole(mustRole int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := ReadAPIContext(r.Context())
			currentRole := ctx.Principal.Role
			if !models.IsRole(mustRole) || !models.IsRole(currentRole) {
				ctx.Error(render.ServerError, fmt.Errorf("system role check error, must [%d], role is [%d]", mustRole, currentRole), "role check")
				return
			}
			if currentRole > mustRole {
				ctx.Error(render.DenyAccessError, fmt.Errorf("system role [%d] access deny, mustrole [%d]", currentRole, mustRole), "role check")
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
