package cryptocurrency

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"testing"
)

func TestCryptoCurrencyBitcoin(t *testing.T) {
	var (
		cc  CryptoCurrency
		err error
	)
	cc, err = InitBitcoinClient("https://chain.api.btc.com/v3", true, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name=", cc.CoinName())

	a, p, err := cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2PKH)
	t.Logf("P2PKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2SH)
	t.Logf("P2SH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2WPKH)
	t.Logf("P2WPKH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

	a, p, err = cc.AllocAccount("passwordpassword", "salt", BitcoinKeyTypeP2WSH)
	t.Logf("P2WSH\taddr=%s\tpriv=%s\terr=%v\n", a, p, err)

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
		b, err = cc.BalanceOf(v, 0)
		if err != nil {
			t.Logf("balance\t%s failed,%v\n", v, err)
			continue
		}
		t.Logf("balance\t%s %s\n", v, b.String())
	}

	//get latest block number
	blkNum, err := cc.LastBlockNumber()
	t.Logf("last block= %d %v\n", blkNum, err)

	//get block by number
	bi, err := cc.BlockByNumber(blkNum - 1)
	t.Logf("blk content: %+v %v\n", bi, err)

	//get block by hash
	bi, err = cc.BlockByHash("00000000000000000010ff7ad8443865c89f2de3047e0c5d7f84dedd44e666b5")
	t.Logf("blk content: %+v %v\n", bi, err)

	//get tx in blocks
	txs, err := cc.TransactionsInBlocks(blkNum-3, blkNum)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range txs {
		t.Logf("txid=%s to=%s index=%d amount=%s\n", v.TxHash, v.To, v.Index, v.Value.String())
	}
}

func TestDecodeAddress(t *testing.T) {
	ad, er := btcutil.DecodeAddress("bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej", &chaincfg.MainNetParams)
	t.Log(ad, er)
	t.Log(hex.EncodeToString(ad.ScriptAddress()))
	ad, er = btcutil.DecodeAddress("bc1qj382exhwmuys2ps7jhpjdxfdjeylws0qxw8tql", &chaincfg.MainNetParams)
	t.Log(ad, er)
	t.Log(hex.EncodeToString(ad.ScriptAddress()))
	ad, er = btcutil.DecodeAddress("3CaqMax7uJpjo28w8t9n15FwhzRNQaaE1f", &chaincfg.MainNetParams)
	t.Log(ad, er)
	t.Log(hex.EncodeToString(ad.ScriptAddress()))
	ad, er = btcutil.DecodeAddress("164ibT2TVy6CD7xuzXvcqCPMrMm3jGVtxS", &chaincfg.MainNetParams)
	t.Log(ad, er)
	t.Log(hex.EncodeToString(ad.ScriptAddress()))

	return
}

func TestEncodePubkey(t *testing.T) {
	//e2b765cf8ff1d1738b3d63951c98e04d2af1da7f		1MfmR3rMYnUBzGfpjnUdNW6RGFrMBgv7Jx
	serializedPubKey, er := hex.DecodeString("e2b765cf8ff1d1738b3d63951c98e04d2af1da7f")
	if er != nil {
		t.Fatal(er)
	}
	t.Log("1--->", len(serializedPubKey))
	buf, er := btcutil.NewAddressPubKeyHash(serializedPubKey, &chaincfg.MainNetParams)
	t.Log("2--->", buf, er)

	//944eac9aeedf0905061e95c326992d9649f741e0		bc1qj382exhwmuys2ps7jhpjdxfdjeylws0qxw8tql
	serializedPubKey, er = hex.DecodeString("944eac9aeedf0905061e95c326992d9649f741e0")
	if er != nil {
		t.Fatal(er)
	}
	t.Log("1--->", len(serializedPubKey))
	buf1, er := btcutil.NewAddressWitnessPubKeyHash(serializedPubKey, &chaincfg.MainNetParams)
	t.Log("2--->", buf1, er)

	//701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d		bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej
	serializedPubKey, er = hex.DecodeString("701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d")
	if er != nil {
		t.Fatal(er)
	}
	t.Log("1--->", len(serializedPubKey))
	buf2, er := btcutil.NewAddressWitnessScriptHash(serializedPubKey, &chaincfg.MainNetParams)
	t.Log("2--->", buf2, er)
	return
}
