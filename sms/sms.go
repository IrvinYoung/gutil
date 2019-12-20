package sms

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type SMS interface {
	InitSMS(params ...interface{}) (SMS, error)
	SendSMS(params ...interface{}) (Result, error)
	GetDetail(params ...interface{}) (Receipt, error)
	SupportBy() string
	//others
}

var (
	ErrorSMSNoReceipt = errors.New("no receipt found")
)

type Result struct {
	ReceiptId string
	ErrCode   string
}

// Detail result of query detail
type Receipt struct {
	ErrCode string
	// no error: DELIVERED
	//errors: https://help.aliyun.com/document_detail/101347.html

	SendDate    time.Time
	ReceiveDate time.Time

	SendStat string
	//WAIT：等待回执。
	//FAIL：发送失败。
	//DONE：发送成功。
}

func Using(SMSType string, params ...interface{}) (sms SMS, err error) {
	switch strings.ToLower(SMSType) {
	case "aliyun":
		s := &AliyunSMS{}
		sms, err = s.InitSMS(params...)
	default:
		err = fmt.Errorf("unsupported captcha type %s", SMSType)
	}
	return
}
