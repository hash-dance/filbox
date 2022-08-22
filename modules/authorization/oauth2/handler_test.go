package oauth2

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var js = `{
    "code": 0,
    "errMessage": "",
    "data": {
        "id": 1,
        "created_at": "2020-03-19T21:13:17+08:00",
        "updated_at": "2020-03-30T22:14:46+08:00",
        "deleted_at": null,
        "username": "guows",
        "description": "2020-03-19 13:13:17",
        "password": "a4501034544f755dce5594d98fcc5fa469feee21c671",
        "email": "guows@163.com",
        "phone": "111111111",
        "provider": "local",
        "enabled": true,
        "role": 0
    }
}`

func Test_parseJson(t *testing.T) {
	var field = "data.username"
	data := make(map[string]interface{})
	json.Unmarshal([]byte(js), &data)
	fields := strings.Split(field, ".")
	t.Log(GetValueFromFields(fields, data))
}

func GetValueFromFields(fields []string, source map[string]interface{}) (string, error) {
	if len(fields) == 1 {
		if v, ok := source[fields[0]]; ok {
			return fmt.Sprintf("%v", v), nil
		}
		return "", fmt.Errorf("no this field: %s", fields[0])
	}
	if next, ok := source[fields[0]].(map[string]interface{}); ok {
		return GetValueFromFields(fields[1:], next)
	}
	return "", fmt.Errorf("error parse field %s", fields[0])
}
