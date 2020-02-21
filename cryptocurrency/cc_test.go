package cryptocurrency

import "testing"

func TestCryptoCurrency(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//get account
	a, p, err := cc.AllocAccount("passwordpassword", "salt")
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)
	//decrypt private key
	priv,err := DecryptPrivKey("passwordpassword", "salt", p)
	t.Logf("priv=%s err=%v\n",priv,err)

	//account check
	t.Logf("%s check=%v\n",a,cc.IsValidAccount(a))

}
