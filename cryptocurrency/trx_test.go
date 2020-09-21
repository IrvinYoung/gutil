package cryptocurrency

import "testing"

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
