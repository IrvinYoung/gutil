package cryptocurrency

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
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
func (e *Ethereum) CoinName() string {
	return "Ethereum"
}

func (e *Ethereum) ChainName() string {
	return ChainETH
}

func (e *Ethereum) Symbol() string {
	return "eth"

}
func (e *Ethereum) Decimal() int64 {
	return 18
}

//account
func (e *Ethereum) AllocAccount(password, salt string, params interface{}) (addr, priv string, err error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	//private key
	privateKeyData := crypto.FromECDSA(privateKeyECDSA)
	//priv = hexutil.Encode(privateKeyData)
	priv = hex.EncodeToString(privateKeyData) //without "0x"
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
	if !e.IsValidAccount(addr) {
		err = errors.New("address is invalid")
		return
	}
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
func (e *Ethereum) Transaction(txHash, blkHash string) (txs []*TransactionRecord, err error) {
	b, err := e.c.BlockByHash(e.ctx, common.HexToHash(blkHash))
	if err != nil {
		return
	}
	var (
		msg    types.Message
		amount decimal.Decimal
		h      = common.HexToHash(txHash)
	)
	for k, v := range b.Transactions() {
		if v.Hash() != h {
			continue
		}
		if msg, err = v.AsMessage(types.NewEIP155Signer(e.chainID)); err != nil {
			//log.Println("get tx msg failed,", v.Hash().Hex(), err)
			continue
		}
		if amount, err = ToDecimal(v.Value(), e.Decimal()); err != nil {
			//log.Println("get tx value failed,", v.Hash().Hex(), err)
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
			Data:        v.Data(),
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
			Data:        v.Data(),
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

func (e *Ethereum) MakeTransaction(from []*TxFrom, to []*TxTo, params interface{}) (txSigned interface{}, err error) {
	if len(from) != 1 || len(to) != 1 || params == nil {
		err = errors.New("params error")
		return
	}
	var gasLimit uint64
	switch params.(type) {
	case uint64:
		gasLimit = params.(uint64)
	default:
		err = errors.New("invalid params format")
		return
	}
	if !e.IsValidAccount(from[0].From) || !e.IsValidAccount(to[0].To) {
		err = errors.New("address is invalid")
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
	amount, err := ToWei(to[0].Value, e.Decimal())
	if err != nil {
		return
	}
	//1. get nonce
	nonce, err := e.c.PendingNonceAt(e.ctx, addrFrom)
	if err != nil {
		return
	}
	//2. gas price
	gasPrice, err := e.c.SuggestGasPrice(e.ctx)
	if err != nil {
		return
	}
	println("net price=", gasPrice.String(), "set to:", 100000000000)
	gasPrice = big.NewInt(100000000000)
	//3. check balance
	cost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	cost = new(big.Int).Add(cost, amount)
	balance, err := e.c.BalanceAt(e.ctx, addrFrom, nil)
	if err != nil {
		return
	}
	if balance.Cmp(cost) < 0 {
		err = errors.New("no more balance")
		return
	}
	//4. make tx
	tx := types.NewTransaction(nonce, addrTo, amount, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(e.chainID), priv)
	if err != nil {
		return
	}
	txSigned = signedTx
	return
}

func (e *Ethereum) SendTransaction(tx interface{}) (txHash string, txData string, err error) {
	signedTx := tx.(*types.Transaction)
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return
	}
	txData = hexutil.Encode(data)
	if err = e.c.SendTransaction(e.ctx, signedTx); err != nil {
		return
	}
	txHash = signedTx.Hash().Hex()
	return
}

func (e *Ethereum) DecodeRawTransaction(txData string) (from []*TxFrom, to []*TxTo, txHash string, err error) {
	data, err := hexutil.Decode(txData)
	if err != nil {
		return
	}
	tx := &types.Transaction{}
	if err = rlp.DecodeBytes(data, &tx); err != nil {
		return
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(e.chainID))
	if err != nil {
		return
	}
	amount, err := ToDecimal(msg.Value().String(), e.Decimal())
	if err != nil {
		return
	}
	from = []*TxFrom{
		&TxFrom{
			From: msg.From().Hex(),
		},
	}
	to = []*TxTo{
		&TxTo{
			To:    msg.To().Hex(),
			Value: amount,
		},
	}
	txHash = tx.Hash().Hex()
	return
}

func (e *Ethereum) MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo,params interface{}) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}

func (e *Ethereum) ApproveAgent(owner *TxFrom, agent *TxTo) (txSigned interface{}, err error) {
	err = errors.New("not support")
	return
}

func (e *Ethereum) Allowance(owner, agent string) (remain decimal.Decimal, err error) {
	err = errors.New("not support")
	return
}

func (e *Ethereum) EstimateFee(from []*TxFrom, to []*TxTo, d interface{}) (fee decimal.Decimal, limit uint64, err error) {
	if len(from) != 1 || len(to) != 1 {
		err = errors.New("params error")
		return
	}
	if !e.IsValidAccount(from[0].From) || !e.IsValidAccount(to[0].To) {
		err = errors.New("address is invalid")
		return
	}
	addrTo := common.HexToAddress(to[0].To)
	amount, err := ToWei(to[0].Value, e.Decimal())
	if err != nil {
		return
	}
	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from[0].From),
		To:    &addrTo,
		Value: amount,
	}
	if limit, err = e.c.EstimateGas(e.ctx, msg); err != nil {
		return
	}
	price, err := e.c.SuggestGasPrice(e.ctx)
	if err != nil {
		return
	}
	fee = decimal.NewFromBigInt(price, 0).Mul(decimal.NewFromInt(int64(limit)))
	fee, err = ToDecimal(fee, 18)
	return
}

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

func (e *Ethereum) GetTokenInstance(tokenInfo interface{}) (token *EthToken, err error) {
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
	token = &EthToken{
		Ethereum: nec,
		Contract: addr,
	}
	token.token, err = ERC20.NewERC20(common.HexToAddress(addr), nec.c)
	return
}

func ToDecimal(ivalue interface{}, decimals int64) (d decimal.Decimal, err error) {
	var value string
	switch v := ivalue.(type) {
	case string:
		value = v
	case decimal.Decimal:
		value = v.String()
	case *big.Int:
		value = v.String()
	default:
		err = errors.New("not support type")
		return
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

func HexToDecimal(d []byte, decimals int64) (amount decimal.Decimal, err error) {
	str := "0x" + strings.TrimLeft(common.Bytes2Hex(d), "0")
	value, err := hexutil.DecodeBig(str)
	if err != nil {
		return
	}
	amount, err = ToDecimal(value.String(), decimals)
	return
}
