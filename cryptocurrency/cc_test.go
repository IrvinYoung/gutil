package cryptocurrency

import (
	"testing"
)

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
	cc.(*Ethereum).Close()
	e := &Ethereum{Host: "http://127.0.0.1:7545"}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	cc = e

	//eth info
	t.Logf("name=%s symbol=%s decimal=%d total_supply=%s\n",
		cc.Name(), cc.Symbol(), cc.Decimal(), cc.TotalSupply().String())

	//token
	token, err := e.TokenInstance("0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F")
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	a,p,err = token.AllocAccount("passwordpassword", "salt")
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)
	t.Logf("%s is valid %v\n",a,token.IsValidAccount(a))

	//token info
	t.Logf("token name=%s symbol=%s decimal=%d total_supply=%s\n",
		token.Name(), token.Symbol(), token.Decimal(), token.TotalSupply().String())

	//balance
	b, err := cc.BalanceOf("0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", 0)
	t.Logf("eth balance of %s -> %s %v\n",
		"0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", b.String(), err)
	b, err = token.BalanceOf("0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", 0)
	t.Logf("token balance of %s -> %s %v\n",
		"0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", b.String(), err)

	//last block number
	blk, err := cc.LastBlockNumber()
	t.Logf("last eth blk=%d %v\n", blk, err)
	blk, err = token.LastBlockNumber()
	t.Logf("last token blk=%d %v\n", blk, err)
}
