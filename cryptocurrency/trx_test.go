package cryptocurrency

import "testing"

func TestBasic(t *testing.T) {
	trx := &Tron{}

	t.Log("coin=",trx.CoinName())
	t.Log("symbol=",trx.Symbol())
	t.Log("chain=",trx.ChainName())
	t.Log("decimal=",trx.Decimal())
}