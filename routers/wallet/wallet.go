package wallet

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/common/render"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/lotus"
)

// Router handler for wallet
func Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/{phone}", func(r chi.Router) {
		r.Get("/", apicontext.Bind(getWallet))
	})
	return r
}

func getWallet(ctx *apicontext.APIContext) {
	phone := chi.URLParam(ctx.Req, "phone")
	if phone != ctx.Principal.Phone {
		ctx.Error(render.InvalidData, errors.New("非法请求"), "非法请求")
		return
	}
	// 是否有钱包, 有就返回, 没有就创建
	wallet, err := models.GetWallet(mysql.GetClient(), phone)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 调用lotus创建钱包
			walletAddr, err := lotus.WalletNew(ctx.Config.Lotus.Address, ctx.Config.Lotus.Token)
			if err != nil {
				ctx.Error(render.ServerError, err, "创建钱包")
				return
			}
			// 添加到数据库
			wallet = &models.Wallet{
				Phone:   phone,
				Address: walletAddr,
			}
			if err := models.CreateWallet(mysql.GetClient(), wallet); err != nil {
				ctx.Error(render.InvalidData, err, "添加钱包")
				return
			}
			logrus.Infof("创建钱包成功 %s %s", phone, walletAddr)
		} else {
			ctx.Error(render.ServerError, err, "查询数据库")
			return
		}
	}

	// 查询余额
	balance, err := lotus.WalletBalance(ctx.Config.Lotus.Address, ctx.Config.Lotus.Token, wallet.Address)
	if err != nil {
		ctx.Error(render.ServerError, err, "查询余额")
		return
	}

	ctx.JSON(struct {
		Address string `json:"address"`
		Balance string `json:"balance"`
	}{wallet.Address, balance})
}
