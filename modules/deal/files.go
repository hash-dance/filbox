package deal

import (
	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/sirupsen/logrus"
)

func ListFiles(queryInfo *mysql.PageQueryInfo, current *models.Principal) ([]*models.DealResponse, error) {

	scanRes := []struct {
		Filecid string
	}{}
	mysql.GetClient().Raw("select distinct filecid FROM deal_infos where phone = ?", current.Phone).Scan(&scanRes)
	logrus.Infof("filecids %+v", scanRes)

	if len(scanRes) == 0 {
		return nil, nil
	}

	// 获取有交易的cid
	filecids := make([]string, 0)
	for _, v := range scanRes {
		filecids = append(filecids, v.Filecid)
	}

	files := make([]*models.File, 0)
	db := mysql.NewPagingDB().Model(models.File{})

	db = db.Where("phone = ? AND filecid in (?)", current.Phone, filecids)

	switch queryInfo.ConditionKey {
	case mysql.Query:
		db = db.ConditionAND(queryInfo.Condition)
	case mysql.Search:
		db = db.ConditionOR(queryInfo.Condition)
	}
	if err := db.Paging(queryInfo.Pagination).Find(&files).Error; err != nil {
		return nil, err
	}

	fileResponse := make([]*models.DealResponse, 0)
	for i, f := range files {
		totalPrice := big.NewInt(0)
		dealInfos := make([]*models.DealInfo, 0)
		mysql.GetClient().Model(&models.DealInfo{}).Where("phone = ? AND filecid = ?", current.Phone, f.Filecid).Find(&dealInfos)
		for _, dl := range dealInfos {
			p, err := big.FromString(dl.TotalPrice)
			if err != nil {
				continue
			}
			totalPrice = types.BigAdd(totalPrice, p)
		}
		fileResponse = append(fileResponse, &models.DealResponse{
			File:       files[i],
			DealInfos:  dealInfos,
			TotalPrice: totalPrice.String(),
		})
	}
	return fileResponse, nil
}
