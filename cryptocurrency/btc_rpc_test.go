package cryptocurrency

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"
	"testing"
)

func TestCryptoCurrencyBitcoinBtcCom(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("myusername:12345678@127.0.0.1:8334", true, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	a, p, err := cc.AllocAccount("passwordpassword", "salt", BitcoinAddrTypeLegacy)
	t.Logf("P2PKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinAddrTypeP2SH)
	t.Logf("P2SH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinAddrTypeBench32)
	t.Logf("P2WPKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

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
		if b, err = cc.BalanceOf(v, 0); err != nil {
			t.Logf("balance\t%s failed,%v\n", v, err)
			continue
		}
		t.Logf("balance\t%s %s\n", v, b.String())
	}

	//get latest block number
	blkNum, err := cc.LastBlockNumber()
	t.Logf("last block= %d %v\n", blkNum, err)

	//get block by number
	bi, err := cc.BlockByNumber(blkNum)
	t.Logf("blk content: %+v %v\n", bi, err)
	t.Log("fee per byte=", cc.(*BitcoinCore).FeePerBytes)

	//get block by hash
	bi, err = cc.BlockByHash("00000000000000000010ff7ad8443865c89f2de3047e0c5d7f84dedd44e666b5")
	t.Logf("blk content: %+v %v\n", bi, err)

	//get tx in blocks
	txs, err := cc.TransactionsInBlocks(blkNum, blkNum)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range txs {
		t.Logf("txid=%s to=%s index=%d amount=%s\n", v.TxHash, v.To, v.Index, v.Value.String())
	}
}

func TestGetBtcTx(t *testing.T) {
	cc, err := InitBitcoinClient("myusername:12345678@192.168.1.11:8332", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cc.CoinName(), cc.Symbol(), cc.Decimal())
	txs, err := cc.TransactionsInBlocks(1697341, 1697341)
	if err != nil {
		t.Fatal(err)
	}
	for _, tx := range txs {
		t.Log(tx.To)
	}
}