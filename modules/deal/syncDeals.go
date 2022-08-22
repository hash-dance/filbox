package deal

import (
	"time"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/lotus"
	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/big"

	// ltypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/sirupsen/logrus"
)

// SyncDeals 列出订单
func SyncDeals(apiCtx *apicontext.APIContext) {
	ctx := apiCtx.Context
	waittime := time.Minute * 10

	for {
		select {
		case <-ctx.Done():
			logrus.Info("exist sync deals")
			return
		default:
			logrus.Info("start sync deals")
			err := doSync(apiCtx)
			if err != nil {
				logrus.Errorf("%s\n", err.Error())
			}
			time.Sleep(waittime)

		}
	}
}

func doSync(ctx *apicontext.APIContext) error {
	// 获取dealtimeout小时以前的单子进行状态更新
	deals := make([]*models.DealInfo, 0)
	if err := mysql.GetClient().Where("created_at < ?", time.Now().Add(-time.Hour*10)).Where("dealcid != ?", "").Where("status != ? AND status != ?",
		storagemarket.StorageDealActive, storagemarket.StorageDealError).Find(&deals).Error; err != nil {
		return err
	}

	// 找出所有的, dealcid
	for _, d := range deals {
		dealOnline := types.DealOnline{}
		if err := lotus.GetDeal(ctx.Config.Lotus.Address, ctx.Config.Lotus.Token, d.Dealcid, &dealOnline); err != nil {
			logrus.Errorf("get deal to sync: %s", err.Error())
			continue
		}
		if err := updateOnlineDeal(d.ID, &dealOnline); err != nil {
			logrus.Errorf("update deal: %s", err.Error())
			continue
		}
	}
	return nil
}

func updateOnlineDeal(id uint, dealOnline *types.DealOnline) error {
	return models.UpdateDealInfos(mysql.GetClient(), id, map[string]interface{}{
		"dealid":      uint64(dealOnline.Dealid),
		"total_price": big.Mul(big.MustFromString(dealOnline.Priceperepoch), big.NewInt(int64(dealOnline.Duration))).String(),
		"status":      uint64(dealOnline.State),
		"statusmsg":   storagemarket.DealStates[storagemarket.StorageDealStatus(dealOnline.State)],
	}, "dealid", "status", "statusmsg", "total_price")
}
