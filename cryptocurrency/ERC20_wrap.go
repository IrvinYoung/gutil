package cryptocurrency

import (
	"errors"
	"log"

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
func (et *EthToken) Name() string {
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
func (et *EthToken) Transfer(from, to map[string]decimal.Decimal) (txHash string, err error) { return }

//token
func (et *EthToken) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) { return }
func (et *EthToken) IsToken() bool                                                      { return true }

//others
func (et *EthToken) EstimateFee(map[string]interface{}) (fee decimal.Decimal, err error) { return }
