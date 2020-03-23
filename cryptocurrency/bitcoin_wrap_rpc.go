package cryptocurrency

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//implement : bitcoin core RPC

type BitcoinCore struct {
	*Bitcoin
	cli *rpcclient.Client
}

func InitBitcoinCoreRpcClient(btc *Bitcoin) (b *BitcoinCore, err error) {
	if btc == nil {
		err = errors.New("params is invalid")
		return
	}
	b = &BitcoinCore{
		Bitcoin: btc,
	}
	accountAndHost := strings.Split(b.Host, "@")
	if len(accountAndHost) != 2 {
		err = errors.New("host format should be USER:PASSWORD@HOST[:PORT][:PATH]")
		return
	}
	userAndPwd := strings.Split(accountAndHost[0], ":")
	if len(userAndPwd) != 2 {
		err = errors.New("host format should be USER:PASSWORD@HOST[:PORT][:PATH]")
		return
	}
	conn := &rpcclient.ConnConfig{
		Host:         accountAndHost[1],
		User:         userAndPwd[0],
		Pass:         userAndPwd[1],
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	if b.cli, err = rpcclient.New(conn, nil); err != nil {
		return
	}
	//defer b.cli.Shutdown()	//todo:?
	chainInfo, err := b.cli.GetBlockChainInfo()
	if err != nil {
		return
	}
	if !strings.HasPrefix(b.net.Name, chainInfo.Chain) {
		err = errors.New("bitcoin network error")
	}
	return
}

//account
func (b *BitcoinCore) BalanceOf(addr string, blkNum uint64) (d decimal.Decimal, err error) {
	err = errors.New("nonsupport by bitcoin core RPC")
	return
}

//block
func (b *BitcoinCore) LastBlockNumber() (blkNum uint64, err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	count, err := b.cli.GetBlockCount()
	if err != nil {
		return
	}
	blkNum = uint64(count)
	return
}

func (b *BitcoinCore) BlockByNumber(blkNum uint64) (bi interface{}, err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	h, err := b.cli.GetBlockHash(int64(blkNum))
	if err != nil {
		return
	}
	blk, err := b.cli.GetBlock(h)
	if err != nil {
		return
	}
	bi = blk //todo:?
	return
}

func (b *BitcoinCore) BlockByHash(blkHash string) (bi interface{}, err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	h, err := chainhash.NewHashFromStr(blkHash)
	if err != nil {
		return
	}
	blk, err := b.cli.GetBlock(h)
	if err != nil {
		return
	}
	bi = blk //todo:?
	return
}

//transaction
func (b *BitcoinCore) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
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

func (b *BitcoinCore) SendTransaction(txSigned interface{}) (txHash string, err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	tx := txSigned.(*btcutil.Tx)
	h, err := b.cli.SendRawTransaction(tx.MsgTx(), false)
	if err != nil {
		return
	}
	txHash = h.String()
	return
}

func (b *BitcoinCore) EstimateFee(from []*TxFrom, to []*TxTo, params interface{}) (fee decimal.Decimal, txSize uint64, err error) {
	if err = b.estimateFee(); err != nil {
		return
	}
	if b.FeePerBytes <= 0 {
		err = errors.New("fee-per-byte is invalid")
		return
	}
	txSign, err := b.MakeTransaction(from, to, nil)
	if err != nil {
		return
	}
	var (
		msg = txSign.(*btcutil.Tx).MsgTx()
		buf bytes.Buffer
	)
	if err = msg.BtcEncode(&buf, 70002, wire.WitnessEncoding); err != nil {
		return
	}
	txSize = uint64(buf.Len())
	fee, err = ToBtc(int64(buf.Len()) * b.FeePerBytes)
	return
}

//internal
func (b *BitcoinCore) getBlkTxs(blk uint64) (txs []*TransactionRecord, err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	var (
		amount   decimal.Decimal
		pkScript txscript.PkScript
		pkType   txscript.ScriptClass
		addr     btcutil.Address
	)
	bi, err := b.BlockByNumber(blk)
	if err != nil {
		return
	}
	blkInfo := bi.(*wire.MsgBlock)
	for _, v := range blkInfo.Transactions {
		for index, out := range v.TxOut {
			if out.PkScript[0] == txscript.OP_RETURN {
				continue
			}
			pkType = txscript.GetScriptClass(out.PkScript)
			switch pkType {
			case txscript.NonStandardTy: // None of the recognized forms.
				continue
			case txscript.NullDataTy: // Empty data-only (provably prunable).
				continue
			case txscript.MultiSigTy: // Multi signature.
				continue
			case txscript.PubKeyTy: // Pay pubkey.
				if addr, err = btcutil.NewAddressPubKey(out.PkScript[1:34], &chaincfg.TestNet3Params); err != nil {
					return
				}
			default:
				//PubKeyHashTy                             // Pay pubkey hash.
				//WitnessV0PubKeyHashTy                    // Pay witness pubkey hash.
				//ScriptHashTy                             // Pay to script hash.
				//WitnessV0ScriptHashTy                    // Pay to witness script hash.
				if pkScript, err = txscript.ParsePkScript(out.PkScript); err != nil {
					return
				}
				if addr, err = pkScript.Address(&chaincfg.TestNet3Params); err != nil {
					return
				}
			}
			if amount, err = ToBtc(out.Value); err != nil {
				return
			}
			tx := &TransactionRecord{
				TokenFlag:   b.Symbol(),
				Index:       uint64(index),
				LogIndex:    0,
				From:        "",
				To:          addr.EncodeAddress(),
				Value:       amount,
				BlockHash:   blkInfo.BlockHash().String(),
				TxHash:      v.TxHash().String(),
				BlockNumber: blk,
				TimeStamp:   blkInfo.Header.Timestamp.Unix(),
				Data:        nil,
			}
			txs = append(txs, tx)
		}
	}
	return
}

func (b *BitcoinCore) estimateFee() (err error) {
	if b.cli == nil {
		err = errors.New("client is invalid")
	}
	var feePerBytes int64
	//prev bitcoin core version : using estimatefee
	fee, err := b.cli.EstimateFee(1) //todo: test with wallet which is low version
	if err == nil {
		if feePerBytes, err = ToSatoshi(decimal.NewFromFloat(fee)); err != nil {
			return
		}
		b.FeePerBytes = feePerBytes / 1000
		return
	}
	if !strings.Contains(err.Error(), "Method not found") {
		return
	}
	//new bitcoin core version: estimatesmartfee
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()
	query := fmt.Sprintf(`{"jsonrpc":"1.0","id":"%d","method":"estimatesmartfee","params":[1,"CONSERVATIVE"]}`, id)
	resp, err := http.Post("http://"+b.Host, "text/plain", strings.NewReader(query))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var info struct {
		Result struct {
			FeeRate float64 `json:"feerate"`
			Blocks  int64   `json:"blocks"`
		} `json:"result"`
		Error json.RawMessage `json:"error"`
		Id    string          `json:"id"`
	}
	if err = json.Unmarshal(buf, &info); err != nil {
		return
	}
	if info.Error != nil && string(info.Error) != "null" {
		err = errors.New(string(info.Error))
		return
	}
	if info.Id != strconv.FormatInt(id, 10) {
		err = errors.New("wrong response data")
		return
	}
	if feePerBytes, err = ToSatoshi(decimal.NewFromFloat(info.Result.FeeRate)); err != nil {
		return
	}
	b.FeePerBytes = feePerBytes / 1000
	return
}
