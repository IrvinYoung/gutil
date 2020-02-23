package cryptocurrency

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
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

	host := "http://127.0.0.1:7545"
	contractAddress := "0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F"
	//init client
	cc, err = InitEthereumClient(host)
	cc.(*Ethereum).Close()
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	cc = e

	//eth info
	t.Logf("name=%s symbol=%s decimal=%d total_supply=%s\n",
		cc.Name(), cc.Symbol(), cc.Decimal(), cc.TotalSupply().String())

	//token
	token, err := e.TokenInstance(contractAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	a, p, err = token.AllocAccount("passwordpassword", "salt")
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)
	t.Logf("%s is valid %v\n", a, token.IsValidAccount(a))

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

	//get block by number
	block, err := cc.BlockByNumber(blk)
	t.Logf("eth - blk content number: %+v %v\n", block, err)
	block, err = token.BlockByNumber(blk)
	t.Logf("token - blk content number: %+v %v\n", block, err)

	//get block by hash
	block, err = cc.BlockByHash("0x650867ef48d96a1251d4950a1375fc810e50d7023dc8a7f003e3f4ab285d9958")
	t.Logf("eth - blk content by hash: %+v %v\n", block, err)
	block, err = token.BlockByHash("0x650867ef48d96a1251d4950a1375fc810e50d7023dc8a7f003e3f4ab285d9958")
	t.Logf("token - blk content by hash: %+v %v\n", block, err)

	//get transactions
	txs, err := cc.TransactionsInBlocks(blk-3, blk)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range txs {
		t.Logf("eth-txs: %+v\n", v)
	}
	//get token transaction
	txs, err = token.TransactionsInBlocks(blk-3, blk)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range txs {
		t.Logf("token-txs: %+v\n", v)
	}

	//make raw transaction
	from := []*TxFrom{
		&TxFrom{
			From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
		},
	}
	to := []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(50, 0),
		},
	}
	txHash, err := cc.MakeTransaction(from, to)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("txid=%s\n", txHash.(*types.Transaction).Hash().Hex())

}
