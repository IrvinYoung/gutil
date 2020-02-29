package cryptocurrency

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

const (
	BitcoinKeyTypeP2PK  = "P2PK"
	BitcoinKeyTypeP2PKH = "P2PKH"

	BitcoinKeyTypeP2SH  = "P2SH"
	BitcoinKeyTypeP2SHH = "P2SHH"

	BitcoinKeyTypeP2WPKH = "P2WPKH"
	BitcoinKeyTypeP2WSH  = "P2WSH"
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
func (b *Bitcoin) AllocAccount(password, salt string, addrType interface{}) (addr, priv string, err error) {
	if addrType == nil {
		err = errors.New("params lost")
		return
	}
	pubKeyType := addrType.(string)

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return
	}
	wif, err := btcutil.NewWIF(privateKey, b.net, b.isAddrCompress)
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
		if p2pkh, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		addr = p2pkh.EncodeAddress()
	case BitcoinKeyTypeP2SH: // 3***
		var p2sh *btcutil.AddressScriptHash
		if p2sh, err = btcutil.NewAddressScriptHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
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
		if p2wpkh, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
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
	log.Printf("%+v\n", bi)
	blkNum = uint64(bi.Height)
	return
}
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

func ToSatoshi() {}
