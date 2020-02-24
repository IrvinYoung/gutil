package cryptocurrency

import (
	"encoding/hex"
	"errors"
	"github.com/IrvinYoung/gutil/crypto"
	"github.com/shopspring/decimal"
)

type CryptoCurrency interface {
	//basic
	ChainName() string
	CoinName() string
	Symbol() string
	Decimal() int64
	TotalSupply() decimal.Decimal

	//account
	AllocAccount(password, salt string) (addr, priv string, err error)
	IsValidAccount(addr string) bool
	BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error)

	//block
	LastBlockNumber() (blkNum uint64, err error)
	BlockByNumber(blkNum uint64) (bi interface{}, err error)
	BlockByHash(blkHash string) (bi interface{}, err error)

	//transaction
	TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error)
	MakeTransaction([]*TxFrom, []*TxTo) (txSigned interface{}, err error)
	SendTransaction(txSigned interface{}) (txHash string, err error)
	MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo) (txSigned interface{}, err error)
	ApproveAgent(*TxFrom, *TxTo) (txSigned interface{},err error)
	Allowance(owner, agent string) (remain decimal.Decimal, err error)

	//token
	TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error)
	IsToken() bool
}

type TxFrom struct {
	From       string
	PrivateKey string
	Index      uint64 //for UTXO
}

type TxTo struct {
	To    string
	Value decimal.Decimal
}

type TransactionRecord struct {
	TokenFlag   string          `json:"token_flag"` //识别标记
	Index       uint64          `json:"-"`          //tx index
	LogIndex    uint64          `json:"-"`          //log index
	From        string          `json:"from"`
	To          string          `json:"to"`
	Value       decimal.Decimal `json:"value"`
	BlockHash   string          `json:"hash"`
	TxHash      string          `json:"tx_hash"`
	BlockNumber uint64          `json:"blk_num"`
	TimeStamp   int64           `json:"time_stamp"`
	Data        interface{}     `json:"-"` //bitcoin:OP_RETURN  ERC20:input data
}

const (
	ChainBTC = "bitcoin"
	ChainETH = "ethereum"
)

func PasswordCheck(pwd, salt string) (err error) {
	if len(pwd) != 16 {
		err = errors.New("password length must be 16")
		return
	}
	if len(salt) != 8 {
		err = errors.New("salt length must be 8")
		return
	}
	return
}

func encryptPrivKey(pwd, salt, from string) (to string, err error) {
	buf := append([]byte(salt), []byte(from)...)
	ciperText, err := crypto.AesEncrypt(buf, []byte(pwd))
	if err != nil {
		return
	}
	to = hex.EncodeToString(ciperText)
	return
}

func DecryptPrivKey(pwd, salt, from string) (to string, err error) {
	ciperText, err := hex.DecodeString(from)
	if err != nil {
		return
	}
	origText, err := crypto.AesDecrypt(ciperText, []byte(pwd))
	if err != nil {
		return
	}
	perfix := origText[:len([]byte(salt))]
	if string(perfix) != salt {
		err = errors.New("privKey decrypt failed")
		return
	}
	to = string(origText[len(salt):])
	return
}
