package cryptocurrency

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"testing"
)

//"address": "n3dquXGC7y46H9JvDMZ4zVtqygsLXJqaRT",
//private key : cRkzCe5a5mubPHdnrNNBExsxj46NPJSNypVknKRCebQq72jgto97
//"scriptPubKey": "76a914f2a0589d12b18c074d71a6d579e27b2a70ca43b788ac",
//"pubkey": "027eb2d8276c39595ea0731acdb1bee1aa74a1114ff836b08a656fdf2dc8d724f4",

func TestP2PKH(t *testing.T) {
	//from private key
	wif, err := btcutil.DecodeWIF("cRkzCe5a5mubPHdnrNNBExsxj46NPJSNypVknKRCebQq72jgto97")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("private key =", wif.String())
	//public key
	t.Log("public key =", hex.EncodeToString(wif.SerializePubKey()))
	//address
	address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("address =", address.String())
	//public key hash
	t.Log("public key hash =", hex.EncodeToString(address.ScriptAddress()))
	//scriptPubKey
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(address.ScriptAddress()).
		AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG)

	scriptPubKey, err := builder.Script()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("scriptPubKey =", hex.EncodeToString(scriptPubKey))
}

//"address": "2N2fjVuY29MNeb7bA7PLg7uyZrCCnyQU22b",
//private:cN4pjWkdGZfjQ6p5bgjTKPJTKt3kTGn33hSY8THs3XLPGWgRC2vc
//"scriptPubKey": "a914675bbfdc7b8743ace123746cd89942abe91d08ed87",
//"isscript": true,
//"iswitness": false,
//"script": "witness_v0_keyhash",
//"hex": "0014e39a619e30d9c9fbe7e98186ed32512479772948",
//"pubkey": "03c028976ef2fb59efe1216e728f95979b38297cb5acd13966628f2aa605d7b088",
//"embedded": {
	//"isscript": false,
	//"iswitness": true,
	//"witness_version": 0,
	//"witness_program": "e39a619e30d9c9fbe7e98186ed32512479772948",
	//"pubkey": "03c028976ef2fb59efe1216e728f95979b38297cb5acd13966628f2aa605d7b088",
	//"address": "tb1quwdxr83sm8ylhelfsxrw6vj3y3uhw22gqmsa4r",
	//"scriptPubKey": "0014e39a619e30d9c9fbe7e98186ed32512479772948"
//},

func TestP2SH(t *testing.T) {
	//from private key
	wif, err := btcutil.DecodeWIF("cN4pjWkdGZfjQ6p5bgjTKPJTKt3kTGn33hSY8THs3XLPGWgRC2vc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("private key =", wif.String())
	//public key
	t.Log("public key =", hex.EncodeToString(wif.SerializePubKey()))

	data := btcutil.Hash160(wif.SerializePubKey())
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0).AddData(data)
	redeemScript, err := builder.Script()
	if err != nil {
		t.Fatal(err)
	}else{
		t.Log("redeemScript/hex/scriptPubKey =",hex.EncodeToString(redeemScript))
	}

	scriptAddr, err := btcutil.NewAddressScriptHash(redeemScript, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("P2SH address =", scriptAddr.String())

	t.Log("witness_program =",hex.EncodeToString(data))
	apsh,_:= btcutil.NewAddressWitnessPubKeyHash(data,&chaincfg.TestNet3Params)
	t.Log("P2WPHK address =",apsh.String())
}


//"address": "tb1q3w4ealuv2ptyccm5ntawn9n7td7km0ydy3079j",
//private key : cPNuBTZztqXvk1JEyDEngRqbzyHh54rK4Ph2djP79ULCRYLwvTRb
//"scriptPubKey": "00148bab9eff8c50564c63749afae9967e5b7d6dbc8d",
//"isscript": false,
//"iswitness": true,
//"witness_version": 0,
//"witness_program": "8bab9eff8c50564c63749afae9967e5b7d6dbc8d",
//"pubkey": "032bfbd9f081735104206855934bc0b8e4b33638181699bbb5d1bcc8e3d3d57fd0",

func TestP2WPKH(t *testing.T){
	//from private key
	wif, err := btcutil.DecodeWIF("cPNuBTZztqXvk1JEyDEngRqbzyHh54rK4Ph2djP79ULCRYLwvTRb")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("private key =", wif.String())
	//public key
	t.Log("public key =", hex.EncodeToString(wif.SerializePubKey()))

	data := btcutil.Hash160(wif.SerializePubKey())
	t.Log("witness_program =",hex.EncodeToString(data))
	apsh,_:= btcutil.NewAddressWitnessPubKeyHash(data,&chaincfg.TestNet3Params)
	t.Log("P2WPHK address =",apsh.String())
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0).AddData(data)
	redeemScript, err := builder.Script()
	if err != nil {
		t.Fatal(err)
	}else{
		t.Log("scriptPubKey =",hex.EncodeToString(redeemScript))
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
