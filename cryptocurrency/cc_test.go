package cryptocurrency

import "testing"

func TestCryptoCurrencyEthereum(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//get account
	a, p, err := cc.AllocAccount("passwordpassword", "salt")
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)
	//decrypt private key
	priv, err := DecryptPrivKey("passwordpassword", "salt", p)
	t.Logf("priv=%s err=%v\n", priv, err)

	//account check
	t.Logf("%s check=%v\n", a, cc.IsValidAccount(a))

	//init client
	cc, err = InitEthereumClient("http://127.0.0.1:7545")
	t.Logf("%s %v\n", cc.Name(), err)
	cc.(*Ethereum).Close()
	e := &Ethereum{Host: "http://127.0.0.1:7545"}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	cc = e
	t.Logf("%s\n", cc.Symbol())

	//token

	//decimal
	t.Logf("decimal of %s = %d\n", cc.Symbol(), cc.Decimal())

}
