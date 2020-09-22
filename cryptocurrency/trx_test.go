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

	bi, err := trx.BlockByNumber(23468545)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("blk=%+v\n", bi)

	bi, err = trx.BlockByHash("0000000001661a018734c0c1ef3f1fabd01dedfc01219e5b3110fbe8757c9b5e")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("blk=%+v\n", bi)
}

func TestTransaction(t *testing.T) {
	trx := &Tron{Host: "https://api.trongrid.io"}

	//txs, err := trx.Transaction("ade3abbe97afdbc145a6622437dfde346ede59dd0a09ddfc03b0805acd233c4b", "")
	txs, err := trx.Transaction("5eccfb8ac1d300b0a5e3a3c8784e813a1eaebae54d8ea3de881b34653cade629", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", txs)
}
