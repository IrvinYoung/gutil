package cryptocurrency

//PAY ATTENTION
//implement by btc.com
//not bitcoin wallet RPC

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"io/ioutil"
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

const (
	BitcoinKeyTypeP2PKH = "P2PKH"
	BitcoinKeyTypeP2SH  = "P2SH"

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
	case BitcoinKeyTypeP2PKH: // 1***
		var p2pkh *btcutil.AddressPubKeyHash
		if p2pkh, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2pkh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubHash=\t%s\n\n", hex.EncodeToString(p2pkh.ScriptAddress()))

		addr = p2pkh.EncodeAddress()
	case BitcoinKeyTypeP2SH: // 3***
		var p2sh *btcutil.AddressScriptHash
		if p2sh, err = btcutil.NewAddressScriptHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2sh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubHash=\t%s\n\n", hex.EncodeToString(p2sh.ScriptAddress()))

		addr = p2sh.EncodeAddress()
	case BitcoinKeyTypeP2WPKH: // bc1***
		var p2wpkh *btcutil.AddressWitnessPubKeyHash
		if p2wpkh, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), b.net); err != nil {
			return
		}
		//log.Printf("addr=\t%s\n", p2wpkh.EncodeAddress())
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("pubHash=\t%s\n\n", hex.EncodeToString(p2wpkh.ScriptAddress()))

		addr = p2wpkh.EncodeAddress()
	case BitcoinKeyTypeP2WSH:
		var p2wsh *btcutil.AddressWitnessScriptHash
		buf := wif.SerializePubKey()
		if p2wsh, err = btcutil.NewAddressWitnessScriptHash(buf[1:], b.net); err != nil {
			return
		}
		//log.Printf("priv=\t%s\n", privKey)
		//log.Printf("puk=\t%s\n", hex.EncodeToString(wif.SerializePubKey()))
		//log.Printf("addr=\t%s\n", p2wsh.EncodeAddress())
		//log.Printf("pubHash=\t%s\n\n", hex.EncodeToString(p2wsh.ScriptAddress()))

		addr = p2wsh.EncodeAddress()
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
	txs, err = b.getBlkTxs(from)
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
		//don't consider of page size,because:i don't know how to setup
		//log.Printf("%d --> %d\n", blk, page)
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
				if strings.HasPrefix(o.PayType, BitcoinKeyTypeP2WPKH) {
					if buf, err = hex.DecodeString(tx.To); err != nil {
						continue
					}
					if awph, err = btcutil.NewAddressWitnessPubKeyHash(buf, b.net); err != nil {
						continue
					}
					tx.To = awph.EncodeAddress()
				}
				if strings.HasPrefix(o.PayType, BitcoinKeyTypeP2WSH) {
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
	//log.Printf("%+v\n", btp)
	return
}

func (b *Bitcoin) MakeTransaction([]*TxFrom, []*TxTo) (txSigned interface{}, err error) { return }
func (b *Bitcoin) SendTransaction(txSigned interface{}) (txHash string, err error)      { return }
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
