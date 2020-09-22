package cryptocurrency

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/IrvinYoung/gutil/cryptocurrency/tron_lib"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
)

type Tron struct {
	//t       *ERC20.ERC20
	//chainID *big.Int

	Host string
}

func InitTronClient(host string) (t *Tron, err error) {
	t = &Tron{Host: host}
	return
}

func (t *Tron) IsToken() bool {
	return false
}

func (t *Tron) Close() {
}

func (t *Tron) TotalSupply() decimal.Decimal {
	return decimal.Zero
}

//basic
func (t *Tron) CoinName() string {
	return "Tron"
}

func (t *Tron) ChainName() string {
	return ChainTRX
}

func (t *Tron) Symbol() string {
	return "trx"

}
func (t *Tron) Decimal() int64 {
	return 6
}

//account
func (t *Tron) AllocAccount(password, salt string, params interface{}) (addr, priv string, err error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	//private key
	privateKeyData := crypto.FromECDSA(privateKeyECDSA)
	//priv = hexutil.Encode(privateKeyData)
	priv = hex.EncodeToString(privateKeyData) //without "0x"
	println(priv)
	//address
	addr = tron_lib.EncodeCheck(tron_lib.PubkeyToAddressBytes(privateKeyECDSA.PublicKey).Bytes())
	//encrypt private key
	priv, err = encryptPrivKey(password, salt, priv)
	return
}

func (t *Tron) IsValidAccount(addr string) bool {
	if len(addr) != 34 {
		return false
	}
	if addr[0:1] != "T" {
		return false
	}
	_, err := tron_lib.DecodeCheck(addr)
	if err != nil {
		return false
	}
	return true
}

func (t *Tron) BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error) {
	if !t.IsValidAccount(addr) {
		err = errors.New("address is invalid")
		return
	}
	data, err := t.requestPost("/walletsolidity/getaccount", map[string]interface{}{
		"address": addr,
		"visible": true,
	})
	if err != nil {
		return
	}
	var a tron_lib.AccountInfoData
	if err = json.Unmarshal(data, &a); err != nil {
		return
	}
	b = decimal.New(a.Balance, int32(0-t.Decimal()))
	return
}

//block
func (t *Tron) LastBlockNumber() (blkNum uint64, err error) {
	data, err := t.requestGet("/wallet/getnowblock", nil)
	if err != nil {
		return
	}
	var b tron_lib.BlockData
	if err = json.Unmarshal(data, &b); err != nil {
		return
	}
	blkNum = uint64(b.BlockHeader.RawData.Number)
	return
}

func (t *Tron) BlockByNumber(blkNum uint64) (bi interface{}, err error) {
	data, err := t.requestPost("/wallet/getblockbynum", map[string]interface{}{
		"num": blkNum,
	})
	if err != nil {
		return
	}
	var b tron_lib.BlockData
	if err = json.Unmarshal(data, &b); err != nil {
		return
	}
	bi = b
	return
}

func (t *Tron) BlockByHash(blkHash string) (bi interface{}, err error) {
	data, err := t.requestPost("/wallet/getblockbyid", map[string]interface{}{
		"value": blkHash,
	})
	if err != nil {
		return
	}
	var b tron_lib.BlockData
	if err = json.Unmarshal(data, &b); err != nil {
		return
	}
	bi = b
	return
}

//transaction
func (t *Tron) Transaction(txHash, blkHash string) (txs []*TransactionRecord, err error) {
	//todo: not implement
	return
	data, err := t.requestPost("/wallet/gettransactionbyid", map[string]interface{}{
		"value":   txHash,
		"visible": true,
	})
	if err != nil {
		return
	}
	var tx tron_lib.TransactionData
	if err = json.Unmarshal(data, &tx); err != nil {
		return
	}
	return
}

func (t *Tron) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
	//todo: not implement
	return
}

func (t *Tron) MakeTransaction(from []*TxFrom, to []*TxTo, params interface{}) (txSigned interface{}, err error) {
	//todo: not implement
	return
}

func (t *Tron) SendTransaction(tx interface{}) (txHash string, txData string, err error) {
	//todo: not implement
	return
}

func (t *Tron) DecodeRawTransaction(txData string) (from []*TxFrom, to []*TxTo, txHash string, err error) {
	//todo: not implement
	return
}

func (t *Tron) MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo, params interface{}) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}

func (t *Tron) ApproveAgent(owner *TxFrom, agent *TxTo) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}

func (t *Tron) Allowance(owner, agent string) (remain decimal.Decimal, err error) {
	err = errors.New("not support")
	return
}

func (t *Tron) EstimateFee(from []*TxFrom, to []*TxTo, d interface{}) (fee decimal.Decimal, limit uint64, err error) {
	//todo: not implement
	return
}

//token
func (t *Tron) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) {
	//todo: not implement
	return
}

//internal
func (t *Tron) requestGet(url string, d interface{}) (data json.RawMessage, err error) {
	resp, err := http.Get(t.Host + url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

func (t *Tron) requestPost(url string, d interface{}) (data json.RawMessage, err error) {
	buf, err := json.Marshal(d)
	if err != nil {
		return
	}
	r := bytes.NewReader(buf)
	resp, err := http.Post(t.Host+url, "application/json", r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}
