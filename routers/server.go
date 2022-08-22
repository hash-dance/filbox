/*Package routers registry all routers
 */
package routers

import (
	"net/http"
	"time"

	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/server"
	"gitee.com/szxjyt/filbox-backend/routers/authorization/login"
	"gitee.com/szxjyt/filbox-backend/routers/authorization/logout"
	"gitee.com/szxjyt/filbox-backend/routers/deal"
	"gitee.com/szxjyt/filbox-backend/routers/principal"
	"gitee.com/szxjyt/filbox-backend/routers/wallet"
)

// NewServerConfig return a server config
func NewServerConfig(ctx *apicontext.APIContext) *server.Config {
	return &server.Config{
		Context: ctx,
		Timeout: 60 * time.Second,
		AuthRouter: map[string]http.Handler{
			"/logout":    logout.Router(),
			"/principal": principal.Router(),
			"/wallet":    wallet.Router(),
			"/miner":     deal.MinerRouter(),
			"/file":      deal.FileRouter(),
		},
		PublicRouter: map[string]http.Handler{
			"/login": login.Router(),
		},
		CustomRouter: map[string]http.Handler{
			// "/": host.Router(),
		},
	}
}
