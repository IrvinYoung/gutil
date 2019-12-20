package sms

/**
implement: refers to aliyun sms demo
*/

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"time"
)

type AliyunSMS struct {
	client *dysmsapi.Client
}

func (s *AliyunSMS) InitSMS(params map[string]interface{}) (instance SMS, err error) {
	regionId := params["region_id"].(string) //"cn-hangzhou"
	accessKeyId := params["access_key_id"].(string)
	accessSecret := params["access_secret"].(string)

	s.client, err = dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessSecret)
	if err != nil {
		return
	}
	instance = s
	return
}

func (s *AliyunSMS) SendSMS(params map[string]interface{}) (r Result, err error) {
	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"
	request.PhoneNumbers = params["PhoneNumber"].(string)
	request.SignName = params["SignName"].(string)
	request.TemplateCode = params["TemplateCode"].(string)
	request.TemplateParam = params["TemplateParam"].(string) //json string

	response, err := s.client.SendSms(request)
	if err != nil {
		return
	}
	//fmt.Printf("response is %#v\n", response)

	//error code : https://help.aliyun.com/document_detail/101346.html
	r.ErrCode = response.Message
	if response.Code != "OK" {
		err = errors.New(response.Message)
		return
	}
	r.ReceiptId = response.BizId
	return
}

func (s *AliyunSMS) GetDetail(params map[string]interface{}) (r Receipt, err error) {
	request := dysmsapi.CreateQuerySendDetailsRequest()

	request.Scheme = "https"
	request.PhoneNumber = params["PhoneNumber"].(string)
	request.SendDate = params["SendDate"].(string)
	request.BizId = params["BizId"].(string)
	request.CurrentPage = "1"
	request.PageSize = "1"

	response, err := s.client.QuerySendDetails(request)
	if err != nil {
		return
	}

	//fmt.Printf("response is %+v\n", response)
	if response.Code != "OK" {
		err = errors.New(response.Message)
		return
	}
	if response.TotalCount == "0" {
		err = ErrorSMSNoReceipt
		return
	} else if response.TotalCount != "1" {
		err = fmt.Errorf("details count(%s) error", response.TotalCount)
		return
	}

	// no error: DELIVERED
	//errors: https://help.aliyun.com/document_detail/101347.html
	r.ErrCode = response.SmsSendDetailDTOs.SmsSendDetailDTO[0].ErrCode

	if r.SendDate, err = time.ParseInLocation("2006-01-02 15:04:05",
		response.SmsSendDetailDTOs.SmsSendDetailDTO[0].SendDate,
		time.Local); err != nil {
		return
	}
	if r.ReceiveDate, err = time.ParseInLocation("2006-01-02 15:04:05",
		response.SmsSendDetailDTOs.SmsSendDetailDTO[0].ReceiveDate,
		time.Local); err != nil {
		return
	}
	//fixme: support custom timezone

	switch response.SmsSendDetailDTOs.SmsSendDetailDTO[0].SendStatus {
	case 1:
		r.SendStat = "WAIT" //等待回执
	case 2:
		r.SendStat = "FAIL" //发送失败
	case 3:
		r.SendStat = "DONE" //发送成功
	default:
		r.SendStat = "FAIL" //发送失败
	}
	return
}

func (s *AliyunSMS) SupportBy() string {
	return "aliyun"
}
