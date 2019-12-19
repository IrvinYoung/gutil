package captcha

import (
	"fmt"
	"strings"
)

type Captcha interface {
	InitCaptcha(params ...interface{}) (instance Captcha, err error)
	NewCaptcha(params ...interface{}) (id string, data []byte)
	VerifyCaptcha(id, digits string) bool
}

func Using(captchaType string, params ...interface{}) (ct Captcha, err error) {
	switch strings.ToLower(captchaType) {
	case "dchest":
		c := &CaptchaDchest{}
		ct, err = c.InitCaptcha(params...)
	default:
		err = fmt.Errorf("unsupported captcha type %s", captchaType)
	}
	return
}
