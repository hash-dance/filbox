package models

import (
	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/big"

	// ltypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/jinzhu/gorm"
)

type File struct {
	types.BaseModel
	Phone    string `json:"phone"`
	Filename string `json:"filename"`
	Filecid  string `json:"filecid"`
	Size     int64  `json:"size"`
}

type DealInfo struct {
	types.BaseModel
	Phone      string                          `json:"phone"`
	Filecid    string                          `json:"filecid"`
	Dealcid    string                          `json:"dealcid"`
	Dealid     uint64                          `json:"dealid"`
	Miner      string                          `json:"miner"`
	Price      string                          `json:"price"`
	Duration   int                             `json:"duration"`
	TotalPrice string                          `json:"totalPrice"`
	Wallet     string                          `json:"wallet"`
	Status     storagemarket.StorageDealStatus `json:"status"`
	Statusmsg  string                          `json:"statusmsg"`
}

type DealResponse struct {
	File       *File       `json:"file"`
	DealInfos  []*DealInfo `json:"dealInfos"`
	TotalPrice string      `json:"totalPrice"`
}

func CreateFile(tx *gorm.DB, file *File) error {
	return tx.Create(file).Error
}

func UpdateDealInfo(tx *gorm.DB, dealOnline *types.DealOnline, fields ...string) error {
	deal := &DealInfo{
		Dealid:     uint64(dealOnline.Dealid),
		TotalPrice: big.Mul(big.MustFromString(dealOnline.Priceperepoch), big.NewInt(int64(dealOnline.Duration))).String(),
		Status:     uint64(dealOnline.State),
		Statusmsg:  storagemarket.DealStates[storagemarket.StorageDealStatus(dealOnline.State)],
	}
	return tx.Model(&DealInfo{}).Where("dealcid = ?", dealOnline.Proposalcid.NAMING_FAILED).Select(fields).Updates(deal).Error
}

func UpdateDealInfos(tx *gorm.DB, id uint, data map[string]interface{}, fields ...string) error {
	return tx.Model(&DealInfo{}).Where("id = ?", id).Select(fields).Updates(data).Error
}
