package principal

import (
	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"gitee.com/szxjyt/filbox-backend/models"
)

// ListPrincipal list principals
func ListPrincipal(queryInfo *mysql.PageQueryInfo, current *models.Principal) (*[]models.Principal, error) {
	principal := make([]models.Principal, 0)
	db := mysql.NewPagingDB().Model(models.Principal{})
	if models.IsManagerRole(current.Role) {
		db = db.Where("role = ? OR role = ?", models.ROLEMANAGER, models.ROLEUSER)
	} else if models.IsUserRole(current.Role) {
		db = db.Where("id = ?", current.ID)
	}
	switch queryInfo.ConditionKey {
	case mysql.Query:
		db = db.ConditionAND(queryInfo.Condition)
	case mysql.Search:
		db = db.ConditionOR(queryInfo.Condition)
	}
	if err := db.Paging(queryInfo.Pagination).Find(&principal).Error; err != nil {
		return nil, err
	}
	return &principal, nil
}
