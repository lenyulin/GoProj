package ioc

import (
	"GoProj/wedy/internal/service/sms"
	"GoProj/wedy/internal/service/sms/MemSMS"
	"GoProj/wedy/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService() sms.Service {
	return MemSMS.NewMemSMS("appId", "signName")
}
func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("TENCENT_SECRET_ID")
	if !ok {
		panic("TENCENT_SECRET_ID is not set")
	}
	secretKey, ok := os.LookupEnv("TENCENT_SECRET_KEY")
	if !ok {
		panic("TENCENT_SECRET_KEY is not set")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewSMSService(c, "sms_login", "default")
}
