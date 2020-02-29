package cryptocurrency

import (
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestCryptoCurrencyBitcoin(t *testing.T) {
	var cc CryptoCurrency
	cc = &Bitcoin{}
	t.Log("name=", cc.CoinName())

	a, p, err := cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2PK,
		"isCompress": true,
	})
	t.Logf("P2PK\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)
	a, p, err = cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2PKH,
		"isCompress": true,
	})
	t.Logf("P2PKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)
	a, p, err = cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2SH,
		"isCompress": true,
	})
	t.Logf("P2SH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)
	a, p, err = cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2SHH,
		"isCompress": true,
	})
	t.Logf("P2SHH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)
	a, p, err = cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2WPKH,
		"isCompress": true,
	})
	t.Logf("P2WPKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)
	a, p, err = cc.AllocAccount("passwordpassword", "salt", map[string]interface{}{
		"net":        &chaincfg.MainNetParams,
		"keyType":    BitcoinKeyTypeP2WSH,
		"isCompress": true,
	})
	t.Logf("P2WSH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	//check addr valid
	addrs := [...]string{
		"12sQsrgs3Ypo6MqbaYHRKs7ADh8oMphWhC",
		"3D1GrdhTXCDG5vtAcwXt7vuM9nFLuzcEiH",
		"bc1ql3eym3gl875t9hdgu9at9ce93h2rgcg6nnvt5e",
		//wrong
		"19WX95dxY3v92qqPLumvWedgnukh5UgSQB",
		"3Joo2Hm2pxkMq1ztheeuHAUQLC311yjGxs",
		"bc1q5r3zc0swmtstrcnkhld9mzynhemgtuja54e488",
	}
	for _, v := range addrs {
		if cc.IsValidAccount(v){
			t.Logf("valid\t%s",v)
		}else{
			t.Logf("invalid\t%s",v)
		}
	}
}
