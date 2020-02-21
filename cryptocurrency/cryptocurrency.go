package cryptocurrency

import (
	"encoding/hex"
	"errors"
	"github.com/shopspring/decimal"
)

type CryptoCurrency interface {
	//basic
	Name() string
	Symbol() string
	Decimal() int64

	//account
	AllocAccount(password, salt string) (addr, priv string, err error)
	IsValidAccount(addr string) bool
	BalanceOf(addr string) (b decimal.Decimal, err error)

	//block
	LastBlockNumber() (blkNum uint64, err error)
	BlockByNumber(blkNum uint64) (bi interface{}, err error)
	BlockByHash(blkHash string) (bi interface{}, err error)

	//transaction
	TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error)
	Transfer(from, to map[string]decimal.Decimal) (txHash string, err error)

	//token
	TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error)
	//others
	EstimateFee(map[string]interface{}) (fee decimal.Decimal, err error)
}

type TransactionRecord struct {
	TokenFlag   string          `json:"token_flag"` //token 识别标记，如果为主链币交易，为空
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
	ciperText, err := utils.AesEncrypt(buf, []byte(pwd))
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
	origText, err := utils.AesDecrypt(ciperText, []byte(pwd))
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
