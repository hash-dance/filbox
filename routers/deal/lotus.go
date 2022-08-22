package deal

import (
	"net/http"

	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/lotus"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Router handler for miner
func MinerRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/{miner}", func(r chi.Router) {
		r.Get("/", apicontext.Bind(askMiner))
	})
	return r
}

func askMiner(ctx *apicontext.APIContext) {
	miner := chi.URLParam(ctx.Req, "miner")

	logrus.Infof("miner: %s", miner)
	info, err := lotus.ClientQueryAsk(ctx.Config.Lotus.Address, ctx.Config.Lotus.Token, miner)
	if err != nil {
		ctx.Error(render.ServerError, err, "查询矿工")
		return
	}
	ctx.JSON(info)

}
