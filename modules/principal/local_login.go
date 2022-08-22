// Package principal handler principal logic
package principal

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/modules/util"
)

// ValidateLoginUser validate user password
func ValidateLoginUser(username, password string) (*models.Principal, error) {
	principal, err := models.GetPrincipalByPhone(username)
	if err != nil {
		logrus.Errorf("GetPrincipalByPhone error: [%s]", err.Error())
		return nil, err
	}
	if util.CFBDecrypter(principal.Password) == password {
		return principal, nil
	}
	return nil, fmt.Errorf("user [%s] password invalid", username)
}
