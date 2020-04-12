package cryptocurrency

import (
	"encoding/hex"
	"errors"
	"github.com/IrvinYoung/gutil/crypto"
	"github.com/shopspring/decimal"
	"strings"
)

type CryptoCurrency interface {
	//basic
	ChainName() string
	CoinName() string
	Symbol() string
	Decimal() int64
	TotalSupply() decimal.Decimal

	//account
	AllocAccount(password, salt string, params interface{}) (addr, priv string, err error)
	IsValidAccount(addr string) bool
	BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error)

	//block
	LastBlockNumber() (blkNum uint64, err error)
	BlockByNumber(blkNum uint64) (bi interface{}, err error)
	BlockByHash(blkHash string) (bi interface{}, err error)

	//transaction
	Transaction(txHash, blkHash string) (txs []*TransactionRecord, err error)
	TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error)
	MakeTransaction([]*TxFrom, []*TxTo, interface{}) (txSigned interface{}, err error)
	SendTransaction(txSigned interface{}) (txHash string, err error)
	MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo) (txSigned interface{}, err error)
	ApproveAgent(*TxFrom, *TxTo) (txSigned interface{}, err error)
	Allowance(owner, agent string) (remain decimal.Decimal, err error)
	EstimateFee([]*TxFrom, []*TxTo, interface{}) (decimal.Decimal, uint64, error)

	//token
	TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error)
	IsToken() bool
}

type TxFrom struct {
	From       string //address
	PrivateKey string

	Amount decimal.Decimal //for segwit
	TxHash string          //for UTXO
	Index  uint64          //for UTXO
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

func shiftDot(f string, decimals int) (t string, err error) {
	lr := strings.Split(f, ".")
	if len(lr) > 2 || len(lr) < 1 {
		err = errors.New("transform value failed,invalid number:" + f)
		return
	}
	if decimals == 0 {
		t = f
		return
	}
	l, r := lr[0], ""
	if len(lr) == 2 {
		r = lr[1]
	}
	if decimals < 0 {
		decimals = 0 - decimals
		if decimals >= len(l) {
			t = "0." + strings.Repeat("0", decimals-len(l)) + l + r
		} else {
			t = l[:len(l)-decimals] + "." + l[len(l)-decimals:] + r
		}
	} else {
		if decimals >= len(r) {
			t = l + r + strings.Repeat("0", decimals-len(r))
			t = strings.TrimLeft(t, "0")
		} else {
			t = l + r[:decimals] + "." + r[decimals:]
		}
	}
	return
}
