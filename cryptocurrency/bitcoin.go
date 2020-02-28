package cryptocurrency

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
)

type Bitcoin struct {
	Host string
}

const (
	BitcoinKeyTypeP2PK  = "P2PK"
	BitcoinKeyTypeP2PKH = "P2PKH"

	BitcoinKeyTypeP2SH  = "P2SH"
	BitcoinKeyTypeP2SHH = "P2SHH"

	BitcoinKeyTypeP2WPKH = "P2WPKH"
	BitcoinKeyTypeP2WSH  = "P2WSH"
)

//basic
func (b *Bitcoin) ChainName() string            { return ChainBTC }
func (b *Bitcoin) CoinName() string             { return "Bitcoin" }
func (b *Bitcoin) Symbol() string               { return "btc" }
func (b *Bitcoin) Decimal() int64               { return 8 }
func (b *Bitcoin) TotalSupply() decimal.Decimal { return decimal.NewFromInt(21000000) }

//account
func (b *Bitcoin) AllocAccount(password, salt string, params map[string]interface{}) (addr, priv string, err error) {
	if params == nil || len(params) != 3 {
		err = errors.New("params lost")
		return
	}
	bitNet := params["net"].(*chaincfg.Params)
	pubKeyType := params["keyType"].(string)
	isCompressed := params["isCompress"].(bool)

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return
	}
	wif, err := btcutil.NewWIF(privateKey, bitNet, isCompressed)
	if err != nil {
		return
	}
	privKey := wif.String()

	switch pubKeyType {
	//case BitcoinKeyTypeP2PK:
	//	var p2pk *btcutil.AddressPubKey
	//	if p2pk,err = btcutil.NewAddressPubKey(btcutil.Hash160(wif.SerializePubKey()),bitNet);err!=nil{
	//		return
	//	}
	//	addr = p2pk.EncodeAddress()
	case BitcoinKeyTypeP2PKH: // 1***
		var p2pkh *btcutil.AddressPubKeyHash
		if p2pkh, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), bitNet); err != nil {
			return
		}
		addr = p2pkh.EncodeAddress()
	case BitcoinKeyTypeP2SH: // 3***
		var p2sh *btcutil.AddressScriptHash
		if p2sh, err = btcutil.NewAddressScriptHash(btcutil.Hash160(wif.SerializePubKey()), bitNet); err != nil {
			return
		}
		addr = p2sh.EncodeAddress()
	//case BitcoinKeyTypeP2SHH:	// 3***
	//	var p2shh *btcutil.AddressScriptHash
	//	if p2shh, err = btcutil.NewAddressScriptHashFromHash(btcutil.Hash160(wif.SerializePubKey()), bitNet); err != nil {
	//		return
	//	}
	//	addr = p2shh.EncodeAddress()
	case BitcoinKeyTypeP2WPKH: // bc1***
		var p2wpkh *btcutil.AddressWitnessPubKeyHash
		if p2wpkh, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), bitNet); err != nil {
			return
		}
		addr = p2wpkh.EncodeAddress()
	//case BitcoinKeyTypeP2WSH:
	//	var p2wsh *btcutil.AddressWitnessScriptHash
	//	if p2wsh,err = btcutil.NewAddressWitnessScriptHash(btcutil.Hash160(wif.SerializePubKey()),bitNet);err!=nil{
	//		return
	//	}
	//	addr = p2wsh.EncodeAddress()
	default:
		err = fmt.Errorf("key type %s is unsupport", pubKeyType)
		return
	}
	//encrypt private key
	priv, err = encryptPrivKey(password, salt, privKey)
	return
}
func (b *Bitcoin) IsValidAccount(addr string) bool                                     { return false }
func (b *Bitcoin) BalanceOf(addr string, blkNum uint64) (d decimal.Decimal, err error) { return }

//block
func (b *Bitcoin) LastBlockNumber() (blkNum uint64, err error)             { return }
func (b *Bitcoin) BlockByNumber(blkNum uint64) (bi interface{}, err error) { return }
func (b *Bitcoin) BlockByHash(blkHash string) (bi interface{}, err error)  { return }

//transaction
func (b *Bitcoin) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) { return }
func (b *Bitcoin) MakeTransaction([]*TxFrom, []*TxTo) (txSigned interface{}, err error)       { return }
func (b *Bitcoin) SendTransaction(txSigned interface{}) (txHash string, err error)            { return }
func (b *Bitcoin) MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo) (txSigned interface{}, err error) {
	return
}
func (b *Bitcoin) ApproveAgent(*TxFrom, *TxTo) (txSigned interface{}, err error)     { return }
func (b *Bitcoin) Allowance(owner, agent string) (remain decimal.Decimal, err error) { return }

//token
func (b *Bitcoin) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) { return }
func (b *Bitcoin) IsToken() bool                                                      { return false }
