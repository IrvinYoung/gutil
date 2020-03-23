package cryptocurrency

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"testing"
)

func TestBTCRPC(t *testing.T) {
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
		User:         "myusername",
		Pass:         "12345678",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("Block count: %d\n", blockCount)
	}
	//get block hash
	blockCount = 1670838
	bh, err := client.GetBlockHash(blockCount)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("hash of block %d = %s\n", blockCount, bh)
	}
	//net
	cNet, err := client.GetCurrentNet()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cNet)
	}
	//get block by hash
	blk, err := client.GetBlock(bh)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("blk content of %d = %+v\n", blockCount, blk)
		for k, v := range blk.Transactions {
			t.Logf("\t%02d : %s\n", k, v.TxHash())
			for m, n := range v.TxIn {
				t.Logf("\t\t* txIn %d -> %s\n", m, n.PreviousOutPoint)
			}
			for m, n := range v.TxOut {
				//if n.PkScript[0] == txscript.OP_RETURN {
				//	continue
				//}
				var pkScript txscript.PkScript
				var addr btcutil.Address
				pkType := txscript.GetScriptClass(n.PkScript)
				switch pkType {
				case txscript.NullDataTy:
					t.Log("null data")
					continue
				case txscript.NonStandardTy:
					t.Log("nonsupport:", hex.EncodeToString(n.PkScript))
					continue
				case txscript.PubKeyTy:
					if addr, err = btcutil.NewAddressPubKey(n.PkScript[1:34], &chaincfg.TestNet3Params); err != nil {
						t.Log(err)
						continue
					}
				default:
					if pkScript, err = txscript.ParsePkScript(n.PkScript); err != nil {
						t.Log(hex.EncodeToString(n.PkScript))
						t.Log(err)
						continue
					}
					if addr, err = pkScript.Address(&chaincfg.TestNet3Params); err != nil {
						t.Log(err)
						continue
					}
				}
				t.Logf("\t\t# txOut %d %s -> %s => %d\n", m, pkType, addr.EncodeAddress(), n.Value)
			}
		}
	}
	//estimate fee
	fee, err := client.EstimateFee(1)
	if err != nil {
		t.Log("estimate fee error:", err)
	} else {
		t.Logf("fee: %f\n", fee)
	}
	//info
	info, err := client.GetBlockChainInfo()
	if err != nil {
		t.Log("block chain info:", err)
	} else {
		t.Logf("info=%+v\n", info)
	}
}

func TestEstimateBitcoinCoreFee(t *testing.T) {
	cc, err := InitBitcoinClient("myusername:12345678@127.0.0.1:8334", true, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	b := cc.(*BitcoinCore)
	if err = b.estimateFee(); err != nil {
		t.Fatal(err)
	}
	t.Log(b.FeePerBytes)
}
