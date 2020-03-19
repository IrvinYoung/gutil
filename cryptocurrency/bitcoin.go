package cryptocurrency

//PAY ATTENTION
//implement by btc.com
//not bitcoin wallet RPC

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Bitcoin struct {
	net            *chaincfg.Params
	isAddrCompress bool
	Host           string
}

type RespData struct {
	Data   json.RawMessage `json:"data"`
	ErrNo  int             `json:"err_no"` //0 正常,1 找不到该资源,2 参数错误
	ErrMsg string          `json:"err_msg"`
}

type BtcAddrInfo struct {
	Address            string `json:"address"`              // address: string 地址
	Received           int64  `json:"received"`             // received: int 总接收
	Sent               int64  `json:"sent"`                 // sent: int 总支出
	Balance            int64  `json:"balance"`              // balance: int 当前余额
	TxCount            int64  `json:"tx_count"`             // tx_count: int 交易数量
	UnconfirmedTxCount int64  `json:"unconfirmed_tx_count"` // unconfirmed_tx_count: int 未确认交易数量
	Unconfirmed        int64  `json:"unconfirmed_received"` // unconfirmed_received: int 未确认总接收
	UnconfirmedSent    int64  `json:"unconfirmed_sent"`     // unconfirmed_sent: int 未确认总支出
	UnspentTxCount     int64  `json:"unspent_tx_count"`     // unspent_tx_count: int 未花费交易数量
}

type BtcBlockInfo struct {
	Height           int64                  `json:"height"`             // height: int 块高度
	Version          int64                  `json:"version"`            // version: int 块版本
	MrklRoot         string                 `json:"mrkl_root"`          // mrkl_root: string Merkle Root
	CurrMaxTimestamp int64                  `json:"curr_max_timestamp"` // curr_max_timestamp: int 块最大时间戳
	TimeStamp        int64                  `json:"timestamp"`          // timestamp: int 块时间戳
	Bits             int64                  `json:"bits"`               // bits: int bits
	Nonce            int64                  `json:"nonce"`              // nonce: int nonce
	Hash             string                 `json:"hash"`               // hash: string 块哈希
	PrevBlockHash    string                 `json:"prev_block_hash"`    // prev_block_hash: string 前向块哈希，如不存在，则为 null
	NextBlockHash    string                 `json:"next_block_hash"`    // next_block_hash: string 后向块哈希，如不存在，则为 null
	Size             int64                  `json:"size"`               //size: int 块体积
	PoolDifficulty   float64                `json:"pool_difficulty"`    // pool_difficulty: int 矿池难度
	Difficulty       float64                `json:"difficulty"`         // difficulty: int 块难度
	TxCount          int64                  `json:"tx_count"`           // tx_count: int 块奖励
	RewardBlock      int64                  `json:"reward_block"`       // reward_block: int 块奖励
	RewardFees       int64                  `json:"reward_fees"`        // reward_fees: int 块手续费
	CreatedAt        int64                  `json:"created_at"`         // created_at: int 该记录系统处理时间，无业务含义
	Confirmations    int64                  `json:"confirmations"`      // confirmations: int 确认数
	Extras           map[string]interface{} `json:"extras"`
}

type BtcTxInput struct {
	PrevAddress  []string `json:"prev_addresses"` //"prev_addresses": Array<String> 输入地址
	PrevPosition int64    `json:"prev_position"`  // "prev_position": int 前向交易的输出位置
	PrevTxHash   string   `json:"prev_tx_hash"`   // "prev_tx_hash": string 前向交易哈希
	PrevValue    int64    `json:"prev_value"`     // "prev_value": int 前向交易输入金额
	ScriptASM    string   `json:"script_asm"`     // "script_asm": string Script Asm
	ScriptHex    string   `json:"script_hex"`     // "script_hex": string Script Hex
	Sequence     int64    `json:"sequence"`       // "sequence": int Sequence
}

type BtcTxOutput struct {
	Addresses []string `json:"addresses"`  // addresses: Array<String> 输出地址
	Value     int64    `json:"value"`      // value: int 输出金额
	PayType   string   `json:"type"`       // "type": "P2PKH"
	ScriptASM string   `json:"script_asm"` // "script_asm": string Script Asm
	ScriptHex string   `json:"script_hex"` // "script_hex": string Script Hex
	//"script_asm": "OP_DUP OP_HASH160 89950b23ec273834ca18a75fc28580d48b0232a8 OP_EQUALVERIFY OP_CHECKSIG",
	//"script_hex": "76a91489950b23ec273834ca18a75fc28580d48b0232a888ac",
	SpentByTx         string `json:"spent_by_tx"`          // "spent_by_tx": null
	SpentByTxPosition int64  `json:"spent_by_tx_position"` // "spent_by_tx_position": -1
}

type BtcTransaction struct {
	BlockHeight  int64          `json:"block_height"` // block_height: int 所在块高度
	BlockTime    int64          `json:"block_time"`   // block_time: int 所在块时间
	CreatedAt    int64          `json:"created_at"`   // created_at: int 该记录系统处理时间，没有业务含义
	Fee          int64          `json:"fee"`          // fee: int 该交易的手续费
	Hash         string         `json:"hash"`         // hash: string 交易哈希
	Inputs       []*BtcTxInput  `json:"inputs"`
	InputsCount  int64          `json:"inputs_count"` // inputs_count: int 输入数量
	InputsValue  int64          `json:"inputs_value"` // inputs_value: int 输入金额
	IsCoinbase   bool           `json:"is_coinbase"`  // is_coinbase: boolean 是否为 coinbase 交易
	LockTime     int64          `json:"lock_time"`    // lock_time: int lock time
	Outputs      []*BtcTxOutput `json:"outputs"`
	OutputsCount int64          `josn:"outputs_count"` // outputs_count: int 输出数量
	OutputsValue int64          `json:"outputs_value"` // outputs_value: int 输出金额
	Size         int64          `json:"size"`          // size: int 交易体积
	Version      int64          `json:"version"`       // version: int 交易版本号
	//others data from response
	BlockHash     string `json:"block_hash"`
	Confirmations int64  `json:"confirmations"`
	IsDoubleSpent bool   `json:"is_double_spend"`
	IsSWTx        bool   `json:"is_sw_tx"`
	Weight        int64  `json:"weight"`
	Vsize         int64  `json:"vsize"`
	WitnessHash   string `json:"witness_hash"`
	SigOps        int64  `json:"sigops"`
}

type BtcTxPage struct {
	TotalCount int64             `json:"total_count"`
	Page       int64             `json:"page"`
	PageSize   int64             `json:"pagesize"`
	Txs        []*BtcTransaction `json:"list"`
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

func InitBitcoinClient(host string, isAddrCompress bool, defaultNet *chaincfg.Params) (b *Bitcoin, err error) {
	if host == "" || defaultNet == nil {
		err = errors.New("params error")
		return
	}
	b = &Bitcoin{
		net:            defaultNet,
		isAddrCompress: isAddrCompress,
		Host:           host,
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

func (b *Bitcoin) BalanceOf(addr string, blkNum uint64) (d decimal.Decimal, err error) {
	// https://chain.api.btc.com/v3/address/15urYnyeJe3gwbGJ74wcX89Tz7ZtsFDVew
	var ai BtcAddrInfo
	if err = b.request("/address/"+addr, &ai); err != nil {
		return
	}
	d, err = ToBtc(ai.Balance)
	return
}

//block
func (b *Bitcoin) LastBlockNumber() (blkNum uint64, err error) {
	// https://chain.api.btc.com/v3/block/latest
	var bi BtcBlockInfo
	if err = b.request("/block/latest", &bi); err != nil {
		return
	}
	blkNum = uint64(bi.Height)
	return
}

func (b *Bitcoin) BlockByNumber(blkNum uint64) (bi interface{}, err error) {
	// https://chain.api.btc.com/v3/block/3
	// https://chain.api.btc.com/v3/block/3,4,5,latest
	var bbi BtcBlockInfo
	if err = b.request(fmt.Sprintf("/block/%d", blkNum), &bbi); err != nil {
		return
	}
	bi = bbi
	return
}

func (b *Bitcoin) BlockByHash(blkHash string) (bi interface{}, err error) {
	// https://chain.api.btc.com/v3/block/3
	// https://chain.api.btc.com/v3/block/3,4,5,latest
	var bbi BtcBlockInfo
	if err = b.request(fmt.Sprintf("/block/%s", blkHash), &bbi); err != nil {
		return
	}
	bi = bbi
	return
}

//transaction
func (b *Bitcoin) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
	//https://chain.api.btc.com/v3/block/latest/tx?verbose=2
	//verbose，可选，默认为2，选择输出内容等级，含义分别如下：
	//等级 1，包含交易信息；
	//等级 2，包含等级 1、交易的输入、输出地址与金额；
	//等级 3，包含等级 2、交易的输入、输入 script 等信息。
	if from > to {
		err = errors.New("params error")
		return
	}
	txs = make([]*TransactionRecord, 0)
	var tmp []*TransactionRecord
	for i := from; i <= to; i++ {
		if tmp, err = b.getBlkTxs(i); err != nil {
			return
		}
		txs = append(txs, tmp...)
	}
	return
}

func (b *Bitcoin) getBlkTxs(blk uint64) (txs []*TransactionRecord, err error) {
	var (
		btp    *BtcTxPage
		page   int64 = 1
		amount decimal.Decimal
		buf    []byte

		awph *btcutil.AddressWitnessPubKeyHash
		awsh *btcutil.AddressWitnessScriptHash
	)
	for { //each page
		btp = &BtcTxPage{}
		if err = b.request(fmt.Sprintf("/block/%d/tx?verbose=3&page=%d", blk, page), btp); err != nil {
			return
		}
		for _, v := range btp.Txs { //each transaction
			for index, o := range v.Outputs { //each output
				if o.PayType == "NULL_DATA" {
					//log.Printf("txid=%s to=%v, value=%d type=%s\n", v.Hash, o.Addresses, o.Value, o.PayType)
					continue //skip NULL_DATA
				}
				if o.PayType == "P2PKH_MULTISIG" {
					//log.Printf("txid=%s to=%v, value=%d type=%s\n", v.Hash, o.Addresses, o.Value, o.PayType)
					continue // skip P2PKH_MULTISIG
				}
				if len(o.Addresses) > 1 {
					//log.Printf("txid=%s to=%v, value=%d type=%s\n", v.Hash, o.Addresses, o.Value, o.PayType)
					continue //skip multi sign
				}
				//log.Printf("txid=%s to=%v, value=%d type=%s\n", v.Hash, o.Addresses, o.Value, o.PayType)
				if amount, err = ToBtc(o.Value); err != nil {
					//log.Printf("transfer amount failed, txid=%s, o_index=%d\n", v.Hash, index)
					continue
				}
				tx := &TransactionRecord{
					TokenFlag:   b.Symbol(),
					Index:       uint64(index),
					Value:       amount,
					To:          o.Addresses[0],
					BlockHash:   v.BlockHash,
					TxHash:      v.Hash,
					BlockNumber: uint64(v.BlockHeight),
					TimeStamp:   v.BlockTime,
				}
				if strings.HasPrefix(o.PayType, "P2WPKH") {
					if buf, err = hex.DecodeString(tx.To); err != nil {
						continue
					}
					if awph, err = btcutil.NewAddressWitnessPubKeyHash(buf, b.net); err != nil {
						continue
					}
					tx.To = awph.EncodeAddress()
				}
				if strings.HasPrefix(o.PayType, "P2WSH") {
					if buf, err = hex.DecodeString(tx.To); err != nil {
						continue
					}
					if awsh, err = btcutil.NewAddressWitnessScriptHash(buf, b.net); err != nil {
						continue
					}
					tx.To = awsh.EncodeAddress()
				}
				txs = append(txs, tx)
			}
		}
		//break
		if btp.PageSize > int64(len(btp.Txs)) {
			break
		} else {
			page++
		}
	}
	return
}

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
	txStr, err := btcMsgToHex(tx.MsgTx())
	log.Printf("hasSegWit=%v txStr= %s %v\n", tx.HasWitness(), txStr, err)
	return tx, nil
}

func (b *Bitcoin) SendTransaction(txSigned interface{}) (txHash string, err error) { return }

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

func (b *Bitcoin) request(url string, d interface{}) (err error) {
	resp, err := http.Get(b.Host + url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var rd RespData
	if err = json.Unmarshal(buf, &rd); err != nil {
		return
	}
	if rd.ErrNo != 0 {
		err = errors.New(rd.ErrMsg)
		return
	}
	err = json.Unmarshal(rd.Data, &d)
	return
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

func btcMsgToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, 70002, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", errors.New(context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
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
