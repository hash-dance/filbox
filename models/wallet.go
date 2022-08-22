package models

import (
	"fmt"

	"gitee.com/szxjyt/filbox-backend/types"
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	types.BaseModel
	Phone     string `json:"wallet"`
	Address   string `json:"address" gorm:"unique"` // 钱包地址
	SecretKey string `json:"-"`                     // 秘钥
}

// CreatePrincipal add principal if not exists
func CreateWallet(tx *gorm.DB, wallet *Wallet) error {
	old := Wallet{}
	err := tx.Where("phone = ?", wallet.Phone).First(&old).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return tx.Create(wallet).Error
		}
		return err
	}
	return fmt.Errorf("use wallet %s aleary exists", wallet.Phone)
}

func GetWallet(tx *gorm.DB, phone string) (*Wallet, error) {
	wallet := Wallet{}
	return &wallet, tx.Where("phone = ?", phone).First(&wallet).Error
}
