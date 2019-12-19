package captcha

type Captcha interface {
	InitCaptcha(params ...interface{}) (instance Captcha, err error)
	NewCaptcha(params ...interface{}) (id string, data []byte)
	VerifyCaptcha(id, digits string) bool
}

