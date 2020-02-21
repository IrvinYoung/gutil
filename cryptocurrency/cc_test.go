package cryptocurrency

import "testing"

func TestCryptoCurrency(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	a, p, err := cc.AllocAccount("password", "salt")
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)

}
