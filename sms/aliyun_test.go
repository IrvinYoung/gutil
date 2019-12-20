package sms

import "testing"

func TestAliyunSMS(t *testing.T) {
	//testQueryDetail(t)
	testSend(t)
}

func testQueryDetail(t *testing.T) {
	s, err := Using("aliyun", map[string]interface{}{
		"regionId":     "cn-hangzhou",
		"accessKeyId":  "AK",
		"accessSecret": "SK"})
	if err != nil {
		t.Fatal(err)
	}
	d, err := s.GetDetail(map[string]interface{}{
		"PhoneNumber":"13012345678",
		"SendDate":"20191212",
		"BizId":"1231231231231231",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v/\n", d)
}

func testSend(t *testing.T) {
	s, err := Using("aliyun", map[string]interface{}{
		"regionId":     "cn-hangzhou",
		"accessKeyId":  "AK",
		"accessSecret": "SK"})
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.SendSMS(map[string]interface{}{
		"PhoneNumber":"13012345678",
		"SignName":"20191219",
		"TemplateCode":"",
		"TemplateParam":"1231231231231231",
	})
	t.Log(r, err)
}
