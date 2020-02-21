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

	//decimal
	t.Logf("decimal of %s = %d\n", cc.Symbol(), cc.Decimal())

	//token
	token, err := e.TokenInstance("0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F")
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	//token info
	t.Logf("token name=%s symbol=%s decimal=%d total_supply=%s\n",
		token.Name(), token.Symbol(), token.Decimal(),token.TotalSupply().String())


}
