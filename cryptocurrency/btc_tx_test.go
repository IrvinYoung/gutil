package cryptocurrency

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"
	"testing"
)

//out
// addr=	mk92td6Dm6LZ9AgMLSBaTEBgSo2FhZhMN8
// priv=	cSXSbV34fHjr4NL7ScVCEVGcDpak497SepBRM9uAPaXDnb5UDTRF
// puk=	0233a616993bf37e15781a30438e2e90ee67a6b1725ab33f816ae9e022ee4b2101
// pubKeyHash=	32b3507768cb095ba3de7418e9e0e4db52c862ea
//
// addr=	2N4eAFwmfErArLn7FF3rxiN4xQdgpLLeWjZ
// priv=	cNRJWdvd76VS3HPNCeyny9wx7b4kKgSKy8sUcn2w6XfQecgVk9G4
// puk=	03bf5136c3935ef9b1d72610f29e8ae916ce594e7de64d70a4de946ca012dea9ba
// pubKeyHash=	7cffd64cffc47d1d5f3f87be2eefba232bbf9f18
//
// addr=	tb1qdg8acm6hg5q0gumxqhzpp69852pupehh9f4tac
// priv=	cQCbugctcCW76PQ5xWcVUVCa9HHcRBMV3KTbnmnnFxP8D9xNCJ27
// puk=	03db8044abce66d7151cfb2349e575fc86fa9272e75b09a2d775460bf69d73c2fe
// pubKeyHash=	6a0fdc6f574500f4736605c410e8a7a283c0e6f7

func TestMakeTx(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	from := []*TxFrom{
		&TxFrom{
			//addr:n3dquXGC7y46H9JvDMZ4zVtqygsLXJqaRT
			//amount: 0.00022
			From:       "n3dquXGC7y46H9JvDMZ4zVtqygsLXJqaRT",
			TxHash:     "c10c33df9736739623684eb6515743b7032585c13f7b28a5d6fc76606eef9922",
			PrivateKey: "cRkzCe5a5mubPHdnrNNBExsxj46NPJSNypVknKRCebQq72jgto97",
			Index:      0,
			Amount:     decimal.New(22, -5),
		},
		&TxFrom{
			//addr: 2N2fjVuY29MNeb7bA7PLg7uyZrCCnyQU22b
			//amount: 0.00071
			From:       "2N2fjVuY29MNeb7bA7PLg7uyZrCCnyQU22b",
			TxHash:     "0be8589ad971b36d421cdb6a792eb0da06b206b0be63b7bbcdb73a927e771990",
			PrivateKey: "cN4pjWkdGZfjQ6p5bgjTKPJTKt3kTGn33hSY8THs3XLPGWgRC2vc",
			Index:      0,
			Amount:     decimal.New(71, -5),
		},
		&TxFrom{
			//addr: tb1q3w4ealuv2ptyccm5ntawn9n7td7km0ydy3079j
			//amount: 0.00093
			From:       "tb1q3w4ealuv2ptyccm5ntawn9n7td7km0ydy3079j",
			TxHash:     "2fe793d7375d0cc493c66dad083319c488e7f054c55067108c0bd8a8ad842975",
			PrivateKey: "cPNuBTZztqXvk1JEyDEngRqbzyHh54rK4Ph2djP79ULCRYLwvTRb",
			Index:      0,
			Amount:     decimal.New(93, -5),
		},
	}

	to := []*TxTo{
		&TxTo{ //P2KH
			To:    "mk92td6Dm6LZ9AgMLSBaTEBgSo2FhZhMN8",
			Value: decimal.New(6, -4),
		},
		&TxTo{ //P2SH
			To:    "2N4eAFwmfErArLn7FF3rxiN4xQdgpLLeWjZ",
			Value: decimal.New(6, -4),
		},
		&TxTo{ //P2WPKH
			To:    "tb1qdg8acm6hg5q0gumxqhzpp69852pupehh9f4tac",
			Value: decimal.New(6, -4),
		},
	}
	txSign, err := cc.MakeTransaction(from, to, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", txSign)
}

func TestMakeTx1(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	from := []*TxFrom{
		&TxFrom{
			From:       "mk92td6Dm6LZ9AgMLSBaTEBgSo2FhZhMN8",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cSXSbV34fHjr4NL7ScVCEVGcDpak497SepBRM9uAPaXDnb5UDTRF",
			Index:      0,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "2N4eAFwmfErArLn7FF3rxiN4xQdgpLLeWjZ",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cNRJWdvd76VS3HPNCeyny9wx7b4kKgSKy8sUcn2w6XfQecgVk9G4",
			Index:      1,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "tb1qdg8acm6hg5q0gumxqhzpp69852pupehh9f4tac",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cQCbugctcCW76PQ5xWcVUVCa9HHcRBMV3KTbnmnnFxP8D9xNCJ27",
			Index:      2,
			Amount:     decimal.New(6, -4),
		},
	}

	to := []*TxTo{
		&TxTo{ //P2KH
			To:    "2NGZrVvZG92qGYqzTLjCAewvPZ7JE8S8VxE",
			Value: decimal.New(17, -5),
		},
	}

	txSign, err := cc.MakeTransaction(from, to, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", txSign)
}

func TestPublishTx(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	//make transaction
	from := []*TxFrom{
		&TxFrom{
			From:       "mk92td6Dm6LZ9AgMLSBaTEBgSo2FhZhMN8",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cSXSbV34fHjr4NL7ScVCEVGcDpak497SepBRM9uAPaXDnb5UDTRF",
			Index:      0,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "2N4eAFwmfErArLn7FF3rxiN4xQdgpLLeWjZ",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cNRJWdvd76VS3HPNCeyny9wx7b4kKgSKy8sUcn2w6XfQecgVk9G4",
			Index:      1,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "tb1qdg8acm6hg5q0gumxqhzpp69852pupehh9f4tac",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cQCbugctcCW76PQ5xWcVUVCa9HHcRBMV3KTbnmnnFxP8D9xNCJ27",
			Index:      2,
			Amount:     decimal.New(6, -4),
		},
	}
	to := []*TxTo{
		&TxTo{ //P2KH
			To:    "2NGZrVvZG92qGYqzTLjCAewvPZ7JE8S8VxE",
			Value: decimal.New(17, -5),
		},
	}
	txSign, err := cc.MakeTransaction(from, to, 0)
	if err != nil {
		t.Fatal(err)
	}
	//publish it
	txHash, err := cc.SendTransaction(txSign)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("txid=", txHash)
}

func TestEstimateFee(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	//make transaction
	from := []*TxFrom{
		&TxFrom{
			From:       "mk92td6Dm6LZ9AgMLSBaTEBgSo2FhZhMN8",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cSXSbV34fHjr4NL7ScVCEVGcDpak497SepBRM9uAPaXDnb5UDTRF",
			Index:      0,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "2N4eAFwmfErArLn7FF3rxiN4xQdgpLLeWjZ",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cNRJWdvd76VS3HPNCeyny9wx7b4kKgSKy8sUcn2w6XfQecgVk9G4",
			Index:      1,
			Amount:     decimal.New(6, -4),
		},
		&TxFrom{
			From:       "tb1qdg8acm6hg5q0gumxqhzpp69852pupehh9f4tac",
			TxHash:     "5e8c892ce0b951f03f2c3652b8577570b8b4f415bed3df5775c360a5bebe6d25",
			PrivateKey: "cQCbugctcCW76PQ5xWcVUVCa9HHcRBMV3KTbnmnnFxP8D9xNCJ27",
			Index:      2,
			Amount:     decimal.New(6, -4),
		},
	}
	to := []*TxTo{
		&TxTo{ //P2KH
			To:    "2NGZrVvZG92qGYqzTLjCAewvPZ7JE8S8VxE",
			Value: decimal.New(17, -5),
		},
	}
	fee, txSize, err := cc.EstimateFee(from, to, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("fee=", fee, "txSize=", txSize)
}

func TestGetTxByBtcCom(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	txs, isPending, err := cc.Transaction("a0417c6c3fdb0d969fe0648c3dc761f541ab92f3405cf195a0612330b714e461",
		"000000000000000000124c0502c5f20d3a48ce906b98e5455d4d9ff74ebfd547")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("is pending=",isPending)
	for k, v := range txs {
		t.Logf("\t%d -> tx=%+v\n", k, v)
	}
}

func TestGetTxByRPC(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("myusername:12345678@192.168.3.14:18332", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	txs, isPending, err := cc.Transaction("56f87210814c8baef7068454e517a70da2f2103fc3ac7f687e32a228dc80e115",
		"00000000000008c4e99525336570ce0817625cb9b9d73ddab5579c32dbb96fb8")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("is pending=",isPending)
	for k, v := range txs {
		t.Logf("\t%d -> tx=%+v\n", k, v)
	}
}
