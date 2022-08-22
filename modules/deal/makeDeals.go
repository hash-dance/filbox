package deal

import (
	"time"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	apicontext "gitee.com/szxjyt/filbox-backend/modules/context"
	"gitee.com/szxjyt/filbox-backend/modules/lotus"
	"github.com/sirupsen/logrus"
)

// MakeDeals 列出订单
func MakeDeals(apiCtx *apicontext.APIContext) {
	ctx := apiCtx.Context
	waittime := time.Minute * 10

	for {
		select {
		case <-ctx.Done():
			logrus.Info("exist make deals")
			return
		default:
			logrus.Info("start make deals")
			err := makedeal(apiCtx)
			if err != nil {
				logrus.Errorf("%s\n", err.Error())
			}
			time.Sleep(waittime)

		}
	}
}

func makedeal(ctx *apicontext.APIContext) error {
	deals := make([]*models.DealInfo, 0)
	if err := mysql.GetClient().Where("dealcid = ?", "").Find(&deals).Error; err != nil {
		return err
	}

	// 找出所有的, dealcid
	for _, d := range deals {
		file := models.File{}
		if err := mysql.GetClient().Model(&models.File{}).Where("filecid = ?", d.Filecid).First(&file).Error; err != nil {
			logrus.Errorf("find deal file by cid %s, %s", d.Filecid, err.Error())
			continue
		}
		if cid, err := lotus.Makedeal(ctx.Config.Lotus.Address, ctx.Config.Lotus.Token, &file, d); err != nil {
			logrus.Errorf("make deal err %s", err.Error())
			continue
		} else {
			logrus.Infof("cid is %s", cid)
			if err := updateDealcid(d.ID, cid); err != nil {
				logrus.Errorf("insert deal cid err %s", err.Error())
			}
		}
	}
	return nil
}

func updateDealcid(id uint, dealcid string) error {
	return models.UpdateDealInfos(mysql.GetClient(), id, map[string]interface{}{
		"dealcid": dealcid,
	}, "dealcid")
}
