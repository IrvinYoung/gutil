package sms

import "testing"

func TestAliyunSMS(t *testing.T) {
	testQueryDetail(t)
	testSend(t)
}

func testQueryDetail(t *testing.T) {
	s, err := Using("aliyun", "cn-hangzhou", "AK", "SK")
	if err != nil {
		t.Fatal(err)
	}
	d, err := s.GetDetail("13012345678", "20191219", "1231231231231231")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v/\n", d)
}

func testSend(t *testing.T) {
	s, err := Using("aliyun", "cn-hangzhou", "AK", "SK")
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.SendSMS(
		"13012345678",
		"SMSTEST",
		"SMS_123123",
		`{"code":"123123"}`,
	)
	t.Log(r, err)
}
