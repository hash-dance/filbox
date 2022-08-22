package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/types"
)

const (
	// ROLEADMIN admin
	ROLEADMIN = iota
	// ROLEMANAGER manager
	ROLEMANAGER
	// ROLEUSER user default role
	ROLEUSER
)

// Principal user's basic information
type Principal struct {
	types.BaseModel
	ExternalID string `json:"external_id,omitempty"` // 外部系统ID,唯一识别号
	Username   string `json:"username,omitempty"`    // 登录名称,也是单点登录系统的uuid
	Password   string `json:"password,omitempty"`
	Role       int    `json:"role,omitempty"`
	Phone      string `gorm:"unique" json:"phone"`
}

// CheckHasRoleAdmin if there is an adminRole in system
func CheckHasRoleAdmin() bool {
	if err := mysql.GetClient().Where("role = ?", ROLEADMIN).First(&Principal{}).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// CreatePrincipal add principal if not exists
func CreatePrincipal(tx *gorm.DB, principal *Principal) error {
	old := Principal{}
	err := tx.Where("phone = ?", principal.Phone).First(&old).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return tx.Create(principal).Error
		}
		return err
	}
	return fmt.Errorf("user [%s-%s] aleary exists", principal.Phone)
}

// CreateOrUpdatePrincipalByExternalID create or update by external_id
// used for external user, only change username
func CreateOrUpdatePrincipalByExternalID(tx *gorm.DB, principal *Principal) (*Principal, error) {
	old := Principal{}
	err := tx.Where("external_id = ?", principal.ExternalID).First(&old).Error
	if gorm.IsRecordNotFoundError(err) {
		return principal, tx.Create(principal).Error
	}
	if old.Username == principal.Username {
		return &old, nil
	}
	old.Username = principal.Username
	return &old, tx.Save(&old).Error
}

// GetPrincipalByPhone get principal info by phone
func GetPrincipalByPhone(phone string) (*Principal, error) {
	principal := Principal{}
	return &principal, mysql.GetClient().Where("phone = ?", phone).First(&principal).Error
}

// GetPrincipalByID get by id
func GetPrincipalByID(id uint) (*Principal, error) {
	principal := Principal{BaseModel: types.BaseModel{ID: id}}
	return &principal, mysql.GetClient().Find(&principal).Error
}

// GetUnscopedPrincipalByID get by id include deleted record
func GetUnscopedPrincipalByID(id uint) (*Principal, error) {
	principal := Principal{BaseModel: types.BaseModel{ID: id}}
	return &principal, mysql.GetClient().Unscoped().Find(&principal).Error
}

// ModifyPrincipalRoleByID principal modify
func ModifyPrincipalRoleByID(tx *gorm.DB, role int, id int) error {
	old, err := GetPrincipalByID(uint(id))
	if err != nil {
		return err
	}
	if old.Role == role {
		return nil
	}
	if IsRole(role) {
		old.Role = role
		return tx.Save(old).Error
	}
	return fmt.Errorf("role [%d] invalid", role)
}

// APIFormat format principal
func (p *Principal) APIFormat() *types.Principal {
	return &types.Principal{
		BaseModel:  p.BaseModel,
		ExternalID: p.ExternalID,
		Username:   p.Username,
		Role:       p.Role,
	}
}

// IsRole check role is valid
func IsRole(role int) bool {
	return role == ROLEADMIN || role == ROLEMANAGER || role == ROLEUSER
}

// IsAdminRole role is admin
func IsAdminRole(role int) bool {
	return role == ROLEADMIN
}

// IsManagerRole role is manager
func IsManagerRole(role int) bool {
	return role == ROLEMANAGER
}

// IsUserRole role is user
func IsUserRole(role int) bool {
	return role == ROLEUSER
}
