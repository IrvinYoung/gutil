package cryptocurrency

import (
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"strings"

	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

//wrap ERC20.go
//because: that is created by abigen

type EthToken struct {
	*Ethereum

	Contract    string
	name        string
	symbol      string
	dcm         int64
	totalSupply decimal.Decimal

	token *ERC20.ERC20
}

func InitEthereumTokenClient(host, addr string) (et *EthToken, err error) {
	nec, err := InitEthereumClient(host)
	if err != nil {
		return
	}
	if !nec.IsValidAccount(addr) {
		err = errors.New("contract address is invalid")
		return
	}
	et = &EthToken{
		Ethereum: nec,
		Contract: addr,
	}
	et.token, err = ERC20.NewERC20(common.HexToAddress(addr), nec.c)
	return
}

func (et *EthToken) Close() {
	et.Ethereum.Close()
}

func (et *EthToken) TotalSupply() (total decimal.Decimal) {
	if et.totalSupply.IsPositive() {
		return et.totalSupply
	}
	amount, err := et.token.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Println("get token name failed,", err)
		return
	}
	et.totalSupply, _ = ToDecimal(amount, et.Decimal())
	total = et.totalSupply
	return
}

//basic
func (et *EthToken) CoinName() string {
	if et.name != "" {
		return et.name
	}
	name, err := et.token.Name(&bind.CallOpts{})
	if err != nil {
		log.Println("get token name failed,", err)
	}
	et.name = name
	return et.name
}
func (et *EthToken) Symbol() string {
	if et.symbol != "" {
		return et.symbol
	}
	symbol, err := et.token.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Println("get token symbol failed,", err)
	}
	et.symbol = symbol
	return et.symbol
}
func (et *EthToken) Decimal() int64 {
	if et.dcm > 0 {
		return et.dcm //todo: maybe token decimal = 0
	}
	d, err := et.token.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Println("get token decimal failed,", err)
		return 18
	}
	et.dcm = d.Int64()
	return et.dcm
}

func (et *EthToken) BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error) {
	if !et.IsValidAccount(addr) {
		err = errors.New("address is invalid")
		return
	}
	amount, err := et.token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr))
	if err != nil {
		return
	}
	b, err = ToDecimal(amount, et.Decimal())
	return
}

//transaction
func (et *EthToken) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
	if from > to {
		err = errors.New("params error")
		return
	}
	txs = make([]*TransactionRecord, 0)
	ti, err := et.token.FilterTransfer(&bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: et.ctx,
	}, nil, nil)
	if err != nil {
		return
	}
	defer ti.Close()
	var (
		amount decimal.Decimal
	)
	for ti.Next() {
		if ti.Event.Raw.Removed {
			continue
		}
		if amount, err = ToDecimal(ti.Event.Value, et.Decimal()); err != nil {
			log.Println("get token value failed,", ti.Event.Raw.TxHash.Hex(), err)
			continue
		}
		tx := &TransactionRecord{
			TokenFlag:   et.Symbol(),
			Index:       uint64(ti.Event.Raw.TxIndex),
			LogIndex:    uint64(ti.Event.Raw.Index),
			From:        ti.Event.From.Hex(),
			To:          ti.Event.To.Hex(),
			Value:       amount,
			BlockHash:   ti.Event.Raw.BlockHash.Hex(),
			TxHash:      ti.Event.Raw.TxHash.Hex(),
			BlockNumber: ti.Event.Raw.BlockNumber,
			Data:        ti.Event.Raw.Data,
		}
		txs = append(txs, tx)
	}
	return
}

func (et *EthToken) MakeTransaction(from []*TxFrom, to []*TxTo, params interface{}) (txSigned interface{}, err error) {
	//make raw transaction, don't run token transfer
	if len(from) != 1 || len(to) != 1 {
		err = errors.New("params error")
		return
	}
	if !et.IsValidAccount(from[0].From) || !et.IsValidAccount(to[0].To) {
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
	amount, err := ToWei(to[0].Value, et.Decimal())
	if err != nil {
		return
	}
	//1. get nonce
	nonce, err := et.c.PendingNonceAt(et.ctx, addrFrom)
	if err != nil {
		return
	}
	//2. gas price
	gasPrice, err := et.c.SuggestGasPrice(et.ctx)
	if err != nil {
		return
	}
	//3. gas limit
	parsed, err := abi.JSON(strings.NewReader(ERC20.ERC20ABI))
	if err != nil {
		return
	}
	data, err := parsed.Pack("transfer", addrTo, amount)
	if err != nil {
		return
	}
	addrToken := common.HexToAddress(et.Contract)
	msg := ethereum.CallMsg{From: addrFrom, To: &addrToken, GasPrice: gasPrice, Value: nil, Data: data}
	gasLimit, err := et.c.EstimateGas(et.ctx, msg)
	if err != nil {
		return
	}
	//4. check eth balance
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	ethBalance, err := et.c.BalanceAt(et.ctx, addrFrom, nil)
	if err != nil {
		return
	}
	if ethBalance.Cmp(fee) < 0 {
		err = errors.New("no more fee")
		return
	}
	//5. check token balance
	balance, err := et.token.BalanceOf(&bind.CallOpts{}, addrFrom)
	if err != nil {
		return
	}
	if balance.Cmp(amount) < 0 {
		err = errors.New("no more balance")
		return
	}
	//6. make tx
	tx := types.NewTransaction(nonce, addrToken, nil, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(et.chainID), priv)
	if err != nil {
		return
	}
	txSigned = signedTx
	return
}

//token
func (et *EthToken) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) {
	cc, err = nil, errors.New("current instance is token, can not init another")
	return
}

func (et *EthToken) IsToken() bool { return true }

func (et *EthToken) MakeAgentTransaction(from string, agent []*TxFrom, to []*TxTo) (txSigned interface{}, err error) {
	if !et.IsValidAccount(from) {
		err = errors.New("address is invalid")
		return
	}
	if len(agent) != 1 || len(to) != 1 {
		err = errors.New("params error")
		return
	}
	if !et.IsValidAccount(agent[0].From) || !et.IsValidAccount(to[0].To) {
		err = errors.New("address is invalid")
		return
	}
	addrFrom := common.HexToAddress(from)
	addrAgent := common.HexToAddress(agent[0].From)
	priv, err := crypto.HexToECDSA(agent[0].PrivateKey)
	if err != nil {
		return
	}
	if crypto.PubkeyToAddress(priv.PublicKey) != addrAgent {
		err = errors.New("private key do not match address")
		return
	}
	addrTo := common.HexToAddress(to[0].To)
	amount, err := ToWei(to[0].Value, et.Decimal())
	if err != nil {
		return
	}
	//1. get nonce
	nonce, err := et.c.PendingNonceAt(et.ctx, addrAgent)
	if err != nil {
		return
	}
	//2. gas price
	gasPrice, err := et.c.SuggestGasPrice(et.ctx)
	if err != nil {
		return
	}
	//3. gas limit
	parsed, err := abi.JSON(strings.NewReader(ERC20.ERC20ABI))
	if err != nil {
		return
	}
	data, err := parsed.Pack("transferFrom", addrFrom, addrTo, amount)
	if err != nil {
		return
	}
	addrToken := common.HexToAddress(et.Contract)
	msg := ethereum.CallMsg{From: addrAgent, To: &addrToken, GasPrice: gasPrice, Value: nil, Data: data}
	gasLimit, err := et.c.EstimateGas(et.ctx, msg)
	if err != nil {
		return
	}
	//4. check eth balance
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	ethBalance, err := et.c.BalanceAt(et.ctx, addrAgent, nil)
	if err != nil {
		return
	}
	if ethBalance.Cmp(fee) < 0 {
		err = errors.New("no more fee")
		return
	}
	//5. check token balance
	balance, err := et.token.BalanceOf(&bind.CallOpts{}, addrFrom)
	if err != nil {
		return
	}
	if balance.Cmp(amount) < 0 {
		err = errors.New("no more balance")
		return
	}
	//6. make tx
	tx := types.NewTransaction(nonce, addrToken, nil, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(et.chainID), priv)
	if err != nil {
		return
	}
	txSigned = signedTx
	return
}

func (et *EthToken) ApproveAgent(owner *TxFrom, agent *TxTo) (txSigned interface{}, err error) {
	if owner == nil || agent == nil {
		err = errors.New("params error")
		return
	}
	if !et.IsValidAccount(owner.From) || !et.IsValidAccount(agent.To) {
		err = errors.New("address is invalid")
		return
	}
	addrOwner := common.HexToAddress(owner.From)
	priv, err := crypto.HexToECDSA(owner.PrivateKey)
	if err != nil {
		return
	}
	if crypto.PubkeyToAddress(priv.PublicKey) != addrOwner {
		err = errors.New("private key do not match address")
		return
	}
	addrAgent := common.HexToAddress(agent.To)
	amount, err := ToWei(agent.Value, et.Decimal())
	if err != nil {
		return
	}
	//1. get nonce
	nonce, err := et.c.PendingNonceAt(et.ctx, addrOwner)
	if err != nil {
		return
	}
	//2. gas price
	gasPrice, err := et.c.SuggestGasPrice(et.ctx)
	if err != nil {
		return
	}
	//3. gas limit
	parsed, err := abi.JSON(strings.NewReader(ERC20.ERC20ABI))
	if err != nil {
		return
	}
	data, err := parsed.Pack("approve", addrAgent, amount)
	if err != nil {
		return
	}
	addrToken := common.HexToAddress(et.Contract)
	msg := ethereum.CallMsg{From: addrOwner, To: &addrToken, GasPrice: gasPrice, Value: nil, Data: data}
	gasLimit, err := et.c.EstimateGas(et.ctx, msg)
	if err != nil {
		return
	}
	//4. check eth balance
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	ethBalance, err := et.c.BalanceAt(et.ctx, addrOwner, nil)
	if err != nil {
		return
	}
	if ethBalance.Cmp(fee) < 0 {
		err = errors.New("no more fee")
		return
	}
	//5. make tx
	tx := types.NewTransaction(nonce, addrToken, nil, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(et.chainID), priv)
	if err != nil {
		return
	}
	txSigned = signedTx
	return
}

func (et *EthToken) Allowance(owner, agent string) (remain decimal.Decimal, err error) {
	if !et.IsValidAccount(owner) || !et.IsValidAccount(agent) {
		err = errors.New("address is invalid")
		return
	}
	a, err := et.token.Allowance(&bind.CallOpts{}, common.HexToAddress(owner), common.HexToAddress(agent))
	if err != nil {
		log.Println(err)
		return
	}
	remain, err = ToDecimal(a, et.Decimal())
	return
}

func (et *EthToken) EstimateFee(from []*TxFrom, to []*TxTo, d interface{}) (fee decimal.Decimal, err error) {
	if len(from) != 1 || len(to) != 1 {
		err = errors.New("params error")
		return
	}
	if !et.IsValidAccount(from[0].From) || !et.IsValidAccount(to[0].To) {
		err = errors.New("address is invalid")
		return
	}
	amount, err := ToWei(to[0].Value, et.Decimal())
	if err != nil {
		return
	}
	parsed, err := abi.JSON(strings.NewReader(ERC20.ERC20ABI))
	if err != nil {
		return
	}
	addrTo := common.HexToAddress(to[0].To)
	data, err := parsed.Pack("transfer", addrTo, amount)
	if err != nil {
		return
	}
	addrToken := common.HexToAddress(et.Contract)

	msg := ethereum.CallMsg{
		From: common.HexToAddress(from[0].From),
		To:   &addrToken,
		Data: data,
	}
	limit, err := et.c.EstimateGas(et.ctx, msg)
	if err != nil {
		return
	}
	price, err := et.c.SuggestGasPrice(et.ctx)
	if err != nil {
		return
	}
	fee = decimal.NewFromBigInt(price, 0).Mul(decimal.NewFromInt(int64(limit)))
	fee, err = ToDecimal(fee, 18)
	return
}
