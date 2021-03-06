package cryptocurrency

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"testing"
)

func TestEthAccount(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}
	a := "0x131231231"
	t.Logf("%s check=%v\n", a, cc.IsValidAccount(a))
	b := "0x53E11118300f77E8AA6a81FD658e27c5a0d88C77"
	t.Logf("%s check=%v\n", b, cc.IsValidAccount(b))
}

func TestCryptoCurrencyEthereum(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//get account
	a, p, err := cc.AllocAccount("passwordpassword", "salt", nil)
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
		cc.CoinName(), cc.Symbol(), cc.Decimal(), cc.TotalSupply().String())

	//token
	token, err := e.TokenInstance(contractAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	a, p, err = token.AllocAccount("passwordpassword", "salt", nil)
	t.Logf("account: addr=%s priv=%s err=%v\n", a, p, err)
	t.Logf("%s is valid %v\n", a, token.IsValidAccount(a))

	//token info
	t.Logf("token name=%s symbol=%s decimal=%d total_supply=%s\n",
		token.CoinName(), token.Symbol(), token.Decimal(), token.TotalSupply().String())

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

	//token approve
	owner := &TxFrom{
		From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
		PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
	}
	agent := &TxTo{
		To:    "0xa5B93c3694b1c9CcFeACcaEebB0E6EA9F13930cC",
		Value: decimal.New(5, 3),
	}
	tx, err := token.ApproveAgent(owner, agent)
	remain, _ := token.Allowance(owner.From, agent.To)
	t.Logf("before approve=%s\n", remain.String())
	txHash, txData, err := token.SendTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("new token approve tx=", txHash, txData)
	remain, _ = token.Allowance(owner.From, agent.To)
	t.Logf("after approve=%s\n", remain.String())

	//token transfer from
	from := []*TxFrom{
		&TxFrom{
			From:       "0xa5B93c3694b1c9CcFeACcaEebB0E6EA9F13930cC",
			PrivateKey: "71d86e526f9ed61088df3c6080821ba3476d5ca2008dff05c2176940b5505cb6",
		},
	}
	to := []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(1, 3),
		},
	}
	tx, err = token.MakeAgentTransaction("0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", from, to,0)
	if err != nil {
		t.Fatal(err)
	}
	tokenBalance, _ := token.BalanceOf("0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", 0)
	t.Logf("before token balance=%s\n", tokenBalance.String())
	txHash, txData, err = token.SendTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	tokenBalance, _ = token.BalanceOf("0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", 0)
	t.Log("new token approve tx=", txHash, txData)
	t.Logf("after token balance=%s\n", tokenBalance.String())
	remain, _ = token.Allowance(owner.From, agent.To)
	t.Logf("after approve=%s\n", remain.String())
}

func TestEstimateEthFee(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//init client
	host := "http://127.0.0.1:7545"
	cc, err := InitEthereumClient(host)
	cc.(*Ethereum).Close()
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	cc = e

	//get fee
	from := []*TxFrom{
		&TxFrom{
			From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
		},
	}
	to := []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(1, 0),
		},
	}
	fee, gasLimit, err := cc.EstimateFee(from, to, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("eth fee=%s gasLimit=%d\n", fee, gasLimit)

	//token fee
	contractAddress := "0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F"
	token, err := e.TokenInstance(contractAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	from = []*TxFrom{
		&TxFrom{
			From:       "0x1B49AC04074F4f3513197Eaa1D6e4fBeea8b7f51",
		},
	}
	to = []*TxTo{
		&TxTo{
			To:    "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			Value: decimal.New(10, -3),
		},
	}
	fee, gasLimit, err = token.EstimateFee(from, to, "0x1B489011a53200bB4d78a04980b5bA2e6af563a5")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token fee=%s gasLimit=%d\n", fee, gasLimit)
}

func TestMakeEthTx(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//init client
	host := "http://127.0.0.1:7545"
	cc, err := InitEthereumClient(host)
	cc.(*Ethereum).Close()
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	cc = e

	//get fee
	from := []*TxFrom{
		&TxFrom{
			From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
		},
	}
	to := []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(1, 0),
		},
	}
	fee, gasLimit, err := cc.EstimateFee(from, to, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("eth fee=%s gasLimit=%d\n", fee, gasLimit)

	//eth transaction
	tx, err := cc.MakeTransaction(from, to, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("txid=%s\n", tx.(*types.Transaction).Hash().Hex())

	//send eth raw transaction
	txHash, txData, err := cc.SendTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("new eth tx=", txHash, txData)

	//token fee
	contractAddress := "0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F"
	token, err := e.TokenInstance(contractAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer token.(*EthToken).Close()
	from = []*TxFrom{
		&TxFrom{
			From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
		},
	}
	to = []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(1, -3),
		},
	}
	fee, gasLimit, err = token.EstimateFee(from, to, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token fee=%s gasLimit=%d\n", fee, gasLimit)

	//token tx
	//make token raw transaction
	tokenBalance, _ := token.BalanceOf("0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", 0)
	t.Logf("before: from - token balance %s %s\n", "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", tokenBalance)
	tokenBalance, _ = token.BalanceOf("0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", 0)
	t.Logf("before: to - token balance %s %s\n", "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", tokenBalance)
	from = []*TxFrom{
		&TxFrom{
			From:       "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1",
			PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
		},
	}
	to = []*TxTo{
		&TxTo{
			To:    "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2",
			Value: decimal.New(1, -3),
		},
	}
	tx, err = token.MakeTransaction(from, to, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token txid=%s\n", tx.(*types.Transaction).Hash().Hex())

	//send token raw transaction
	txHash, txData, err = token.SendTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("new token tx=", txHash, txData)
	tokenBalance, _ = token.BalanceOf("0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", 0)
	t.Logf("after: from - token balance %s %s\n", "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1", tokenBalance)
	tokenBalance, _ = token.BalanceOf("0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", 0)
	t.Logf("after: to - token balance %s %s\n", "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2", tokenBalance)
}

func TestGetAddr(t *testing.T) {
	privStr := "4e7a0e32045d7c732bac92cc36d3d2e8b1bbdc155ccc2394fb8af1b798aa59af"
	priv, err := crypto.HexToECDSA(privStr)
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(priv.PublicKey)
	t.Log("eth address=", address.String())
}

func TestERC20TokenTransaction(t *testing.T) {
	var cc CryptoCurrency
	cc = &Ethereum{}

	//init client
	host := "http://127.0.0.1:7545"
	var err error
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()

	contractAddress := "0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F"
	if cc, err = e.TokenInstance(contractAddress); err != nil {
		t.Fatal(err)
	}

	//balance
	agent := "0xc056b439F3cC83F7631Fd9fa791B1523dadEc2a1"
	b, err := cc.BalanceOf(agent, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("agnet-->", agent, b)
	}
	fAddr := "0xAbe3716570020Dc0734a6ffbA2e8EBd4042C9Db2"
	b, err = cc.BalanceOf(fAddr, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("from-->", fAddr, b)
	}
	tAddr := "0xfcdE17BA66F8EA6084a37AA04FD888d4Fd9a3847"
	//tAddr := "0x1B49AC04074F4f3513197Eaa1D6e4fBeea8b7f51"
	b, err = cc.BalanceOf(tAddr, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("to-->", tAddr, b)
	}
	//set approve
	tx, err := cc.ApproveAgent(
		&TxFrom{
			From:       fAddr,
			PrivateKey: "90916eb92e2ed3d3c2f92713f6becb3a5b25d52c16dc5c3e3eff2b5d82f1204f",
		},
		&TxTo{
			To:    agent,
			Value: cc.TotalSupply(),
		})
	if err != nil {
		t.Fatal(err)
	}
	txid, txData, err := cc.SendTransaction(tx)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("approve txid=", txid, txData)
	}
	//get allowance
	remain, err := cc.Allowance(fAddr, agent)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("remain->", fAddr, remain)
	}

	//make tx
	tx, err = cc.MakeAgentTransaction(fAddr,
		[]*TxFrom{
			&TxFrom{
				From:       agent,
				PrivateKey: "c821b8cdfe1b7dd195ffb00d17245f945ab893253ee846d987e362658a92585c",
			}},
		[]*TxTo{
			&TxTo{
				To:    tAddr,
				Value: decimal.New(8, -3),
			}},0)
	if err != nil {
		t.Fatal(err)
	}
	if txid, txData, err = cc.SendTransaction(tx); err != nil {
		t.Fatal(err)
	} else {
		t.Log("token txid=", txid, txData)
	}

	//balance
	b, err = cc.BalanceOf(agent, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("agent-->", agent, b)
	}
	b, err = cc.BalanceOf(fAddr, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("from-->", fAddr, b)
	}
	b, err = cc.BalanceOf(tAddr, 0)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("to-->", tAddr, b)
	}
	return
}

func TestGetTransaction(t *testing.T) {
	//init client
	host := "http://127.0.0.1:7545"
	var err error
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()

	tx, err := e.Transaction("0x4e00243c7e763d85bed3291467bc24c8474f9f3757475282daa3e4dab065af1a",
		"0x11f8d385ebf47b9844bcb9c9db5b09b4bfebcb18e88393da0789ffdcab8c707e")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx=%+v\n", tx)
}

func TestGetTokenTransaction(t *testing.T) {
	//init client
	host := "http://127.0.0.1:7545"
	var err error
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()

	contractAddress := "0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F"
	token, err := e.TokenInstance(contractAddress)
	if err != nil {
		t.Fatal(err)
	}

	txs, err := token.Transaction("0x689fb8fc0c318c252b9b696b399f1886df3f56e7d5698bd5007f6ae8f75feb5d",
		"0x3f3a4138e9898e254938d6b9e34ca44b74ffa42d274187dd41ae913c5f8653d9")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range txs {
		t.Logf("\t%d -> tx=%+v\n", k, v)
	}
}

func TestDecodeEthTx(t *testing.T) {
	//init client
	host := "http://127.0.0.1:7545"
	var err error
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()

	f, to, txhash, err := e.DecodeRawTransaction("0xf86e048504a817c8008252089456d7ec6e9359eafb2a66d32072ecfb574fe240bc880f43fc2c04ee000080822d46a0ce0d0fe6f0814b04345f6106c0c0e93f969496a33b541e3633de4fcdc7c5d5aea02f808df575b952ea3fc6fd8bbc893436daec733b87f1eda08a2d6a13f65b72a6")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", txhash)
	t.Logf("%+v\n", f[0])
	t.Logf("%+v\n", to[0])
}

func TestDecodeEthTokenTx(t *testing.T) {
	//init client
	host := "http://127.0.0.1:7545"
	var err error
	e := &Ethereum{Host: host}
	if err = e.Init(); err != nil {
		t.Fatal("init ethereum failed,", err)
	}
	defer e.Close()
	et, err := e.TokenInstance("0x6aa0cfdEFFefDd4968Cf550f9160D78AF9afd65F")
	if err != nil {
		t.Fatal(err)
	}

	//transfer
	txData := "0xf8ab478504a817c800828fb4946aa0cfdeffefdd4968cf550f9160d78af9afd65f80b844a9059cbb000000000000000000000000abe3716570020dc0734a6ffba2e8ebd4042c9db200000000000000000000000000000000000000000000000000038d7ea4c68000822d45a0b1a237027739866a15a04f911c39761e5bcc28e3f6a8644477f41a0d18330f9aa03c71b17a7b7dac23f64e90040585902ceb36b234ab75a373d7faa4e6c2a7f3b4"
	//transfer from
	//txData := "0xf8cb488504a817c80082ac62946aa0cfdeffefdd4968cf550f9160d78af9afd65f80b86423b872dd000000000000000000000000abe3716570020dc0734a6ffba2e8ebd4042c9db2000000000000000000000000fcde17ba66f8ea6084a37aa04fd888d4fd9a3847000000000000000000000000000000000000000000000000001c6bf526340000822d45a05c463fab1b6554cc016d292467d9dbaff3405ce0ef7a7fe02d58c493d79a47c0a009a83e6a4dd7e4ad644440900c6da16385e1f1a65119a325235e9aa49d6c881b"
	f, to, txhash, err := et.DecodeRawTransaction(txData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", txhash)
	t.Logf("%+v\n", f[0])
	t.Logf("%+v\n", to[0])
}
