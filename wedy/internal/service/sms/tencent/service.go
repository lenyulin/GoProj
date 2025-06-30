package tencent

import (
	"context"
	"fmt"
	ekit "github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SMSService struct {
	client   *sms.Client
	appId    string
	signName string
}

func (s *SMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = ekit.ToPtr(s.appId)
	request.SignName = ekit.ToPtr(s.signName)
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	response, err := s.client.SendSms(request)
	if err != nil {
		fmt.Println("SendSms err")
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("SendSms err, code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *SMSService) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

func NewSMSService(client *sms.Client, appId, signName string) *SMSService {
	return &SMSService{
		client:   client,
		appId:    appId,
		signName: signName,
	}
}
