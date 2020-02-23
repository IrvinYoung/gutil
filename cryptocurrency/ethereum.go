package cryptocurrency

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
)

type Ethereum struct {
	ctx     context.Context
	c       *ethclient.Client
	t       *ERC20.ERC20
	chainID *big.Int

	Host string
}

var (
	ReEthereumAccount = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

func InitEthereumClient(host string) (e *Ethereum, err error) {
	e = &Ethereum{Host: host}
	err = e.Init()
	return
}

func (e *Ethereum) Init() (err error) {
	e.ctx = context.Background()
	if e.c, err = ethclient.DialContext(e.ctx, e.Host); err != nil {
		return
	}
	e.chainID, err = e.c.NetworkID(e.ctx)
	return
}

func (e *Ethereum) IsToken() bool {
	return false
}

func (e *Ethereum) Close() {
	e.c.Close()
}

func (e *Ethereum) TotalSupply() decimal.Decimal {
	return decimal.Zero
}

//basic
func (e *Ethereum) Name() string {
	return "Ethereum"
}

func (e *Ethereum) Symbol() string {
	return "eth"

}
func (e *Ethereum) Decimal() int64 {
	return 18
}

//account
func (e *Ethereum) AllocAccount(password, salt string) (addr, priv string, err error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	//private key
	privateKeyData := crypto.FromECDSA(privateKeyECDSA)
	priv = hexutil.Encode(privateKeyData)
	//address
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	addr = address.String()
	//encrypt private key
	priv, err = encryptPrivKey(password, salt, priv)
	return
}

func (e *Ethereum) IsValidAccount(addr string) bool {
	return ReEthereumAccount.MatchString(addr)
}

func (e *Ethereum) BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error) {
	var blk *big.Int
	if blkNum == 0 {
		blk = nil
	} else {
		blk = big.NewInt(int64(blkNum))
	}
	amount, err := e.c.BalanceAt(e.ctx, common.HexToAddress(addr), blk)
	if err != nil {
		return
	}
	b, err = ToDecimal(amount, e.Decimal())
	return
}

//block
func (e *Ethereum) LastBlockNumber() (blkNum uint64, err error) {
	//s, err := e.c.SyncProgress(e.ctx)	//XXX: not work? why?
	c, err := rpc.DialContext(e.ctx, e.Host)
	if err != nil {
		return
	}
	var raw json.RawMessage
	if err = c.CallContext(e.ctx, &raw, "eth_blockNumber"); err != nil {
		return
	}
	var num string
	if err = json.Unmarshal(raw, &num); err != nil {
		return
	}
	blkNum = hexutil.MustDecodeUint64(num)
	return
}

func (e *Ethereum) BlockByNumber(blkNum uint64) (bi interface{}, err error) {
	b, err := e.c.BlockByNumber(e.ctx, big.NewInt(int64(blkNum)))
	if err != nil {
		return
	}
	b.SanityCheck()
	bi = b
	return
}
func (e *Ethereum) BlockByHash(blkHash string) (bi interface{}, err error) {
	b, err := e.c.BlockByHash(e.ctx, common.HexToHash(blkHash))
	if err != nil {
		return
	}
	bi = b
	return
}

//transaction
func (e *Ethereum) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
	if from > to {
		err = errors.New("params error")
		return
	}
	txs = make([]*TransactionRecord, 0)
	var tmp []*TransactionRecord
	for i := from; i <= to; i++ {
		if tmp, err = e.getBlkTxs(i); err != nil {
			return
		}
		txs = append(txs, tmp...)
	}
	return
}

func (e *Ethereum) getBlkTxs(blk uint64) (txs []*TransactionRecord, err error) {
	txs = make([]*TransactionRecord, 0)
	b, err := e.c.BlockByNumber(e.ctx, big.NewInt(int64(blk)))
	if err != nil {
		return
	}
	var (
		msg    types.Message
		amount decimal.Decimal
	)
	for k, v := range b.Transactions() {
		//todo: maybe some transaction is invalid
		if msg, err = v.AsMessage(types.NewEIP155Signer(e.chainID)); err != nil {
			log.Println("get tx msg failed,", v.Hash().Hex(), err)
			continue
		}
		if amount, err = ToDecimal(v.Value(), e.Decimal()); err != nil {
			log.Println("get tx value failed,", v.Hash().Hex(), err)
			continue
		}
		tx := &TransactionRecord{
			TokenFlag:   e.Symbol(),
			Index:       uint64(k),
			From:        msg.From().Hex(),
			Value:       amount,
			BlockHash:   b.Hash().Hex(),
			TxHash:      v.Hash().Hex(),
			BlockNumber: b.NumberU64(),
			TimeStamp:   int64(b.Time()),
		}
		if msg.To() != nil {
			tx.To = msg.To().Hex()
		} else {
			tx.To = "" //new contract
		}
		txs = append(txs, tx)
		//log.Printf("%d\t%s : %s -> %s %s\n",
		//	k, tx.TxHash, tx.From, tx.To, tx.Value.String())
	}
	return
}

func (e *Ethereum) MakeTransaction(from []*TxFrom, to []*TxTo) (txHash string, err error) {
	if len(from) != 1 || len(to) != 1 {
		err = errors.New("params error")
		return
	}
	addrFrom := common.HexToAddress(from[0].From)
	priv, err := crypto.HexToECDSA(from[0].PrivateKey)
	if err != nil {
		return
	}
	if crypto.PubkeyToAddress(priv.PublicKey) != addrFrom {
		err = errors.New("private key do not match address")
		return
	}
	addrTo := common.HexToAddress(to[0].To)
	//if !from[0].Value.Equal(to[0].Value) {
	//	err = errors.New("amount FROM is not equal TO")
	//	return
	//}
	amount, err := ToWei(to[0].Value, e.Decimal())
	if err != nil {
		return
	}
	//1. get nonce
	nonce, err := e.c.PendingNonceAt(e.ctx, addrFrom)
	if err != nil {
		return
	}
	log.Println("nonce=", nonce)
	//2. gas price
	gasPrice, err := e.c.SuggestGasPrice(e.ctx)
	if err != nil {
		return
	}
	log.Println("gasPrice=", gasPrice)
	//3. gas limit	//compute again, not use default value: 21000
	msg := ethereum.CallMsg{From: addrFrom, To: &addrTo, GasPrice: gasPrice, Value: amount, Data: nil}
	gasLimit, err := e.c.EstimateGas(e.ctx, msg)
	if err != nil {
		return
	}
	log.Println("gasLimit=", gasLimit)

	//tx := types.NewTransaction(nonce, addrTo, amount, gasLimit, gasPrice, nil)

	return
}

func (e *Ethereum) SendTransaction() {}

//token
func (e *Ethereum) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) {
	var addr string
	switch tokenInfo.(type) {
	case string:
		addr = tokenInfo.(string)
	default:
		err = errors.New("need contract address")
		return
	}
	if !e.IsValidAccount(addr) {
		err = errors.New("contract address is invalid")
		return
	}

	nec, err := InitEthereumClient(e.Host)
	if err != nil {
		return
	}
	token := &EthToken{
		Ethereum: nec,
		Contract: addr,
	}
	token.token, err = ERC20.NewERC20(common.HexToAddress(addr), nec.c)
	cc = token
	return
}

//others
func (e *Ethereum) EstimateFee(map[string]interface{}) (fee decimal.Decimal, err error) { return }

func ToDecimal(ivalue interface{}, decimals int64) (d decimal.Decimal, err error) {
	var value string
	switch v := ivalue.(type) {
	case string:
		value = v
	case decimal.Decimal:
		value = v.String()
	}
	if value, err = shiftDot(value, int(0-decimals)); err != nil {
		return
	}
	d, err = decimal.NewFromString(value)
	return
}

func ToWei(iamount interface{}, decimals int64) (amount *big.Int, err error) {
	//todo: consider to use decimal.Decimal
	var value string
	switch v := iamount.(type) {
	case string:
		value = iamount.(string)
	case float64:
		//value = strconv.FormatFloat(iamount.(float64), 'f', int(decimals), 64)	//todo: precision error
		err = errors.New("not support float64")
		return
	case int64:
		value = strconv.FormatInt(iamount.(int64), 10)
	case int:
		value = strconv.Itoa(iamount.(int))
	case decimal.Decimal:
		value = iamount.(decimal.Decimal).String()
	case *decimal.Decimal:
		value = (*v).String()
	default:
		err = errors.New("not support type")
		return
	}
	if value, err = shiftDot(value, int(decimals)); err != nil {
		return
	}
	d, err := decimal.NewFromString(value)
	if err != nil {
		return
	}
	amount = d.Coefficient()
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
