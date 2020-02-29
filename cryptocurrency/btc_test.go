package cryptocurrency

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"
	"testing"
)

func TestCryptoCurrencyBitcoin(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	a, p, err := cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2PK)
	t.Logf("P2PK\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2PKH)
	t.Logf("P2PKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2SH)
	t.Logf("P2SH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2SHH)
	t.Logf("P2SHH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2WPKH)
	t.Logf("P2WPKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2WSH)
	t.Logf("P2WSH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	//check addr valid
	addrs := []string{
		"12sQsrgs3Ypo6MqbaYHRKs7ADh8oMphWhC",
		"3D1GrdhTXCDG5vtAcwXt7vuM9nFLuzcEiH",
		"bc1ql3eym3gl875t9hdgu9at9ce93h2rgcg6nnvt5e",
		//wrong
		"19WX95dxY3v92qqPLumvWedgnukh5UgSQB",
		"3Joo2Hm2pxkMq1ztheeuHAUQLC311yjGxs",
		"bc1q5r3zc0swmtstrcnkhld9mzynhemgtuja54e488",
	}
	for _, v := range addrs {
		if cc.IsValidAccount(v) {
			t.Logf("valid\t%s", v)
		} else {
			t.Logf("invalid\t%s", v)
		}
	}

	//get balance
	addrs = []string{
		"1HtuUatKrJSR8PYs2qSxnxvPuYhf8UiCpB",
		"3K3kszcpPfA79rAP8T9zXQQTLwUYuk5qEc",
		"bc1qnsupj8eqya02nm8v6tmk93zslu2e2z8chlmcej",
	}
	var b decimal.Decimal
	for _, v := range addrs {
		b, err = cc.BalanceOf(v, 0)
		if err != nil {
			t.Logf("balance\t%s failed,%v\n", v, err)
			continue
		}
		t.Logf("balance\t%s %s\n", v, b.String())
	}

	//get block number
	blkNum, err := cc.LastBlockNumber()
	t.Logf("last block= %d %v\n", blkNum, err)
}
