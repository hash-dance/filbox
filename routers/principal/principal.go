// Package principal handler principal logic
package principal

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/szxjyt/filbox-backend/modules/util"
	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/go-chi/chi"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/principal"
)

// Router handler for user
func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", apicontext.Bind(list))
	r.Get("/me", apicontext.Bind(getMe))
	r.Post("/", apicontext.Bind(create, types.CreatePrincipalOptions{}))
	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", apicontext.Bind(update, types.UpdateRoleArg{}))
	})
	return r
}

func create(ctx *apicontext.APIContext, form *types.CreatePrincipalOptions) {
	if err := models.CreatePrincipal(mysql.GetClient(), &models.Principal{
		ExternalID: "1",
		Username:   form.UserName,
		Password:   util.RandString(32),
		Role:       *form.Role,
	}); err != nil {
		ctx.Error(render.ServerError, err, "create principal")
		return
	}
	ctx.JSON(nil)
}

func getMe(ctx *apicontext.APIContext) {
	ctx.JSON(ctx.Principal)
}

// change principal's role
// only admin can access
func update(ctx *apicontext.APIContext, data *types.UpdateRoleArg) {
	current := ctx.Principal
	if !models.IsAdminRole(current.Role) {
		ctx.Error(render.DenyAccessError, fmt.Errorf("user [%s] can't access API", current.Username), "update")
		return
	}

	vl := chi.URLParam(ctx.Req, "id")
	id, err := strconv.Atoi(vl)
	if err != nil {
		ctx.Error(render.InvalidData, err, "get url id")
		return
	}

	if err := models.ModifyPrincipalRoleByID(mysql.GetClient(), data.Role, id); err != nil {
		ctx.Error(render.ServerError, err, "update")
		return
	}
	ctx.JSON(nil)
}

func list(ctx *apicontext.APIContext) {
	principals, err := principal.ListPrincipal(ctx.QueryInfo, ctx.Principal)
	if err != nil {
		ctx.Error(render.ServerError, err, "list principals")
		return
	}
	ctx.JSONPagination(principals)
}
