package cryptocurrency

import (
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestCryptoCurrencyBitcoin(t *testing.T) {
	var cc CryptoCurrency
	cc = &Bitcoin{}
	t.Log("name=",cc.CoinName())

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
}
