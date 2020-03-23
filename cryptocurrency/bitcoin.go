package cryptocurrency

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type Bitcoin struct {
	net            *chaincfg.Params
	isAddrCompress bool
	Host           string
	FeePerBytes    int64
}

type BtcSecretSource struct {
	Inputs []*TxFrom
	Net    *chaincfg.Params
}

const (
	BitcoinAddrTypeLegacy = "legacy"      //P2PKH
	BitcoinAddrTypeP2SH   = "p2sh-segwit" //P2SH-P2WPKH

	BitcoinAddrTypeBench32 = "bech32" //P2WPKH
)

func InitBitcoinClient(host string, isAddrCompress bool, defaultNet *chaincfg.Params) (cc CryptoCurrency, err error) {
	if host == "" || defaultNet == nil {
		err = errors.New("params error")
		return
	}
	b := &Bitcoin{
		net:            defaultNet,
		isAddrCompress: isAddrCompress,
		Host:           host,
	}
	if strings.Contains(host, "btc.com") {
		cc, err = InitBitcoinBtcComClient(b)	//btc.com
	}else{
		cc, err = InitBitcoinCoreRpcClient(b)	//wallet RPC
	}
	return
}

//basic
func (b *Bitcoin) ChainName() string            { return ChainBTC }
func (b *Bitcoin) CoinName() string             { return "Bitcoin" }
func (b *Bitcoin) Symbol() string               { return "btc" }
func (b *Bitcoin) Decimal() int64               { return 8 }
func (b *Bitcoin) TotalSupply() decimal.Decimal { return decimal.NewFromInt(21000000) }

//account
func (b *Bitcoin) AllocAccount(password, salt string, addressType interface{}) (addr, priv string, err error) {
	if addressType == nil {
		err = errors.New("params lost")
		return
	}
	addrType := addressType.(string)

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return
	}
	wif, err := btcutil.NewWIF(privateKey, b.net, b.isAddrCompress)
	if err != nil {
		return
	}
	privKey := wif.String()

	switch addrType {
	case BitcoinAddrTypeLegacy: // 1***
		var p2pkh *btcutil.AddressPubKeyHash
		if p2pkh, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2pkh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubKeyHash=\t%s\n\n", hex.EncodeToString(p2pkh.ScriptAddress()))

		addr = p2pkh.EncodeAddress()
	case BitcoinAddrTypeP2SH: // 3***	//P2SH-P2WPKH
		data := btcutil.Hash160(wif.SerializePubKey())
		builder := txscript.NewScriptBuilder()
		builder.AddOp(txscript.OP_0).AddData(data)
		var redeemScript []byte
		if redeemScript, err = builder.Script(); err != nil {
			return
		}
		var p2sh *btcutil.AddressScriptHash
		if p2sh, err = btcutil.NewAddressScriptHash(redeemScript, b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2sh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubKeyHash=\t%s\n\n", hex.EncodeToString(p2sh.ScriptAddress()))

		addr = p2sh.EncodeAddress()
	case BitcoinAddrTypeBench32: // bc1***
		var p2wpkh *btcutil.AddressWitnessPubKeyHash
		if p2wpkh, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2wpkh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubKeyHash=\t%s\n\n", hex.EncodeToString(p2wpkh.ScriptAddress()))

		addr = p2wpkh.EncodeAddress()
	default:
		err = fmt.Errorf("key type %s is unsupport", addrType)
		return
	}
	//encrypt private key
	priv, err = encryptPrivKey(password, salt, privKey)
	return
}
func (b *Bitcoin) IsValidAccount(addr string) bool {
	if addr == "" {
		return false
	}
	_, err := btcutil.DecodeAddress(addr, b.net)
	if err != nil {
		return false
	}
	return true
}

//transaction
func (b *Bitcoin) MakeTransaction(from []*TxFrom, to []*TxTo, params interface{}) (txSigned interface{}, err error) {
	if from == nil || len(from) == 0 || to == nil || len(to) == 0 {
		err = errors.New("params error")
		return
	}
	var lockTime uint32
	if params != nil {
		switch params.(type) {
		case int:
			lockTime = uint32(params.(int))
		case uint:
			lockTime = uint32(params.(uint))
		case int32:
			lockTime = uint32(params.(int32))
		case uint32:
			lockTime = params.(uint32)
		case int64:
			lockTime = uint32(params.(int64))
		case uint64:
			lockTime = uint32(params.(uint64))
		default:
			err = errors.New("params is invalid")
			return
		}
	}
	mtx := wire.NewMsgTx(wire.TxVersion)
	var (
		addr btcutil.Address

		txHash  *chainhash.Hash
		prevOut *wire.OutPoint
		txIn    *wire.TxIn

		satoshi  int64
		pkScript []byte
	)
	//inputs
	for _, input := range from {
		if txHash, err = chainhash.NewHashFromStr(input.TxHash); err != nil { //txid -> hash
			return
		}
		prevOut = wire.NewOutPoint(txHash, uint32(input.Index)) //make out point
		txIn = wire.NewTxIn(prevOut, []byte{}, nil)
		if params != nil && lockTime != 0 {
			txIn.Sequence = wire.MaxTxInSequenceNum - 1
		}
		mtx.AddTxIn(txIn)
	}
	//outputs
	for _, v := range to {
		if !v.Value.IsPositive() || v.Value.GreaterThanOrEqual(b.TotalSupply()) {
			err = fmt.Errorf("invalid amount: %s = %s", v.To, v.Value.String())
			return
		}
		if addr, err = btcutil.DecodeAddress(v.To, b.net); err != nil {
			err = fmt.Errorf("invalid address or key: %v", err)
			return
		}
		switch addr.(type) {
		case *btcutil.AddressPubKeyHash:
		case *btcutil.AddressScriptHash:
		case *btcutil.AddressWitnessPubKeyHash:
		//case *btcutil.AddressPubKey: //todo: ?
		default:
			err = fmt.Errorf("invalid address or key: %s", v.To)
			return
		}
		if !addr.IsForNet(b.net) {
			err = fmt.Errorf("invalid address: %s is for the wrong network", v.To)
			return
		}
		if pkScript, err = txscript.PayToAddrScript(addr); err != nil {
			err = fmt.Errorf("failed to generate pay-to-address script: %v", err)
			return
		}
		if satoshi, err = ToSatoshi(v.Value); err != nil {
			err = fmt.Errorf("failed to convert amount = %s", v.Value.String())
			return
		}
		txOut := wire.NewTxOut(satoshi, pkScript)
		mtx.AddTxOut(txOut)
	}
	if params != nil && lockTime != 0 {
		mtx.LockTime = lockTime
	}
	//sign
	if err = b.signTx(mtx, from); err != nil {
		return
	}
	tx := btcutil.NewTx(mtx)
	return tx, nil
}
func (b *Bitcoin) MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}
func (b *Bitcoin) ApproveAgent(*TxFrom, *TxTo) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}
func (b *Bitcoin) Allowance(owner, agent string) (remain decimal.Decimal, err error) {
	err = errors.New("not support")
	return
}

//token
func (b *Bitcoin) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) {
	err = errors.New("not support")
	return
}
func (b *Bitcoin) IsToken() bool { return false }

func (b *Bitcoin) signTx(msgTx *wire.MsgTx, from []*TxFrom) (err error) {
	var (
		prevPkScripts [][]byte
		inputValues   []btcutil.Amount
		amt           int64
		addr          btcutil.Address
		pps           []byte
	)
	for _, v := range from {
		if amt, err = ToSatoshi(v.Amount); err != nil {
			return
		}
		inputValues = append(inputValues, btcutil.Amount(amt))

		if addr, err = btcutil.DecodeAddress(v.From, b.net); err != nil {
			err = fmt.Errorf("invalid address or key: %v", err)
			return
		}
		if pps, err = txscript.PayToAddrScript(addr); err != nil {
			return
		}
		prevPkScripts = append(prevPkScripts, pps)
	}
	bss := &BtcSecretSource{
		Inputs: from,
		Net:    b.net,
	}
	return addAllInputScripts(msgTx, prevPkScripts, inputValues, bss)
}

func (bss *BtcSecretSource) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	var from *TxFrom
	for _, v := range bss.Inputs {
		if v.From != addr.String() {
			continue
		} else {
			from = v
		}
	}
	if from == nil {
		return nil, false, errors.New("nope")
	}
	//parse private key
	wif, err := btcutil.DecodeWIF(from.PrivateKey)
	if err != nil {
		return nil, false, err
	}
	return wif.PrivKey, wif.CompressPubKey, nil
}

func (bss *BtcSecretSource) GetScript(addr btcutil.Address) ([]byte, error) {
	var from *TxFrom
	for _, v := range bss.Inputs {
		if v.From != addr.String() {
			continue
		} else {
			from = v
		}
	}
	if from == nil {
		return nil, errors.New("nope")
	}
	//parse private key
	wif, err := btcutil.DecodeWIF(from.PrivateKey)
	if err != nil {
		return nil, err
	}

	data := btcutil.Hash160(wif.SerializePubKey())
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0).AddData(data)
	return builder.Script()
}

func (bss *BtcSecretSource) ChainParams() *chaincfg.Params {
	return bss.Net
}

func ToBtc(ivalue interface{}) (b decimal.Decimal, err error) {
	var v string
	switch ivalue.(type) {
	case string:
		v = ivalue.(string)
	case int64:
		v = strconv.FormatInt(ivalue.(int64), 10)
	default:
		err = errors.New("value type is not support")
		return
	}
	if v, err = shiftDot(v, -8); err != nil {
		return
	}
	b, err = decimal.NewFromString(v)
	return
}

func ToSatoshi(d decimal.Decimal) (amount int64, err error) {
	v, err := shiftDot(d.String(), 8)
	if err != nil {
		return
	}
	amount, err = strconv.ParseInt(v, 10, 64)
	return
}

func addAllInputScripts(tx *wire.MsgTx, prevPkScripts [][]byte, inputValues []btcutil.Amount,
	secrets txauthor.SecretsSource) error {

	inputs := tx.TxIn
	hashCache := txscript.NewTxSigHashes(tx)
	chainParams := secrets.ChainParams()

	if len(inputs) != len(prevPkScripts) {
		return errors.New("tx.TxIn and prevPkScripts slices must " +
			"have equal length")
	}
	for i := range inputs {
		pkScript := prevPkScripts[i]
		switch {
		case txscript.IsPayToScriptHash(pkScript):
			err := spendNestedWitnessPubKeyHash(inputs[i], pkScript,
				int64(inputValues[i]), chainParams, secrets,
				tx, hashCache, i)
			if err != nil {
				return err
			}
		case txscript.IsPayToWitnessPubKeyHash(pkScript):
			err := spendWitnessKeyHash(inputs[i], pkScript,
				int64(inputValues[i]), chainParams, secrets,
				tx, hashCache, i)
			if err != nil {
				return err
			}
		default:
			sigScript := inputs[i].SignatureScript
			script, err := txscript.SignTxOutput(chainParams, tx, i,
				pkScript, txscript.SigHashAll, secrets, secrets,
				sigScript)
			if err != nil {
				return err
			}
			inputs[i].SignatureScript = script
		}
	}
	return nil
}

func spendNestedWitnessPubKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, secrets txauthor.SecretsSource,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) error {
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
		chainParams)
	if err != nil {
		return err
	}
	privKey, compressed, err := secrets.GetKey(addrs[0])
	if err != nil {
		return err
	}
	pubKey := privKey.PubKey()

	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	// Next, we'll generate a valid sigScript that'll allow us to spend
	// the p2sh output. The sigScript will contain only a single push of
	// the p2wkh witness program corresponding to the matching public key
	// of this address.
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		return err
	}
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		return err
	}
	bldr := txscript.NewScriptBuilder()
	bldr.AddData(witnessProgram)
	sigScript, err := bldr.Script()
	if err != nil {
		return err
	}
	txIn.SignatureScript = sigScript

	// With the sigScript in place, we'll next generate the proper witness
	// that'll allow us to spend the p2wkh output.
	witnessScript, err := txscript.WitnessSignature(tx, hashCache, idx,
		inputValue, witnessProgram, txscript.SigHashAll, privKey, compressed)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript

	return nil
}

func spendWitnessKeyHash(txIn *wire.TxIn, pkScript []byte,
	inputValue int64, chainParams *chaincfg.Params, secrets txauthor.SecretsSource,
	tx *wire.MsgTx, hashCache *txscript.TxSigHashes, idx int) error {

	// First obtain the key pair associated with this p2wkh address.
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript,
		chainParams)
	if err != nil {
		return err
	}
	privKey, compressed, err := secrets.GetKey(addrs[0])
	if err != nil {
		return err
	}
	pubKey := privKey.PubKey()

	// Once we have the key pair, generate a p2wkh address type, respecting
	// the compression type of the generated key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		return err
	}

	// With the concrete address type, we can now generate the
	// corresponding witness program to be used to generate a valid witness
	// which will allow us to spend this output.
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		return err
	}
	witnessScript, err := txscript.WitnessSignature(tx, hashCache, idx,
		inputValue, witnessProgram, txscript.SigHashAll, privKey, true)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript

	return nil
}

func btcMsgToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, 70002, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", errors.New(context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}
