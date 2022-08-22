package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// 配置
const (
	SignName     = "苏州星际云通"        // 应用名称
	TemplateCode = "SMS_150786099" // 模板名称
)

var tmp_cache = cache.New(5*time.Minute, 10*time.Minute)

//CreateCLient 获取client
func createCLient(region, keyID, keySecret string) (*dysmsapi.Client, error) {
	client, err := dysmsapi.NewClientWithAccessKey(region, keyID, keySecret)
	return client, err
}

//SendSms 发送短信
func SendSms(phone string, region, keyID, keySecret string) error {
	client, err := createCLient(region, keyID, keySecret)
	if err != nil {
		fmt.Println(err)
		return err
	}
	code := GenerateCode()
	code = code[:6]
	jcode := fmt.Sprintf(`{code:"%s"}`, code)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = SignName
	request.TemplateCode = TemplateCode
	request.TemplateParam = jcode
	// request.OutId = "bsdfl45"

	_, err = client.SendSms(request)
	if err != nil {
		logrus.Errorf("sms error:%s", err.Error())
		return err
	}
	tmp_cache.Set(phone, code, time.Minute*5)
	return nil
}

//GenerateCode 生成验证码
func GenerateCode() (code string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code = fmt.Sprintf("%d", r.Int31())
	return
}

//CodeIsEq 判断验证码是否匹配
func CodeIsEq(phone, code string) bool {
	value, ok := tmp_cache.Get(phone)
	if !ok {
		return false
	}
	if value.(string) != code {
		return false
	}
	tmp_cache.Delete(phone)
	return true
}
