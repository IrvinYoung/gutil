package cryptocurrency

import (
	"errors"
	"github.com/IrvinYoung/gutil/cryptocurrency/tron_lib"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestBasic(t *testing.T) {
	trx := &Tron{}

	t.Log("coin=", trx.CoinName())
	t.Log("symbol=", trx.Symbol())
	t.Log("chain=", trx.ChainName())
	t.Log("decimal=", trx.Decimal())
}

func TestAccount(t *testing.T) {
	trx := &Tron{}

	a, p, e := trx.AllocAccount("8765432187654321", "tron", nil)
	t.Log("addr=", a)
	t.Log("priv=", p)
	t.Log("err=", e)

	is := trx.IsValidAccount(a)
	t.Log(a, "check result=", is)
}

func checkPriv2Addr(privKey, addr string) (err error) {
	priv, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return
	}
	a := tron_lib.EncodeCheck(tron_lib.PubkeyToAddressBytes(priv.PublicKey).Bytes())
	if a != addr {
		err = errors.New("private key do not match address")
		return
	}
	return
}

func TestBalance(t *testing.T) {
	trx := &Tron{Host: "https://api.trongrid.io"}
	addr := "TMhDGbyPn17fraYfvMjH58Zrfaix2ZCxz3"
	b, err := trx.BalanceOf(addr, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("balance of", addr, b.String())
}

func TestBlock(t *testing.T) {
	trx := &Tron{Host: "https://api.trongrid.io"}

	b, err := trx.LastBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("last blk=", b)
}
