package cryptocurrency

import (
	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"log"
)

//wrap ERC20.go
//because: that is created by abigen

type EthToken struct {
	*Ethereum

	Contract string
	name     string
	symbol   string
	dcm      int64

	token *ERC20.ERC20
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

//account
func (et *EthToken) AllocAccount(password, salt string) (addr, priv string, err error) {
	return et.Ethereum.AllocAccount(password, salt)
}
func (et *EthToken) IsValidAccount(addr string) bool {
	return et.Ethereum.IsValidAccount(addr)
}
func (et *EthToken) BalanceOf(addr string, blkNum uint64) (b decimal.Decimal, err error) {
	amount, err := et.token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr))
	if err != nil {
		return
	}
	b = ToDecimal(amount, et.Decimal())
	return
}

//block
func (et *EthToken) LastBlockNumber() (blkNum uint64, err error)             { return }
func (et *EthToken) BlockByNumber(blkNum uint64) (bi interface{}, err error) { return }
func (et *EthToken) BlockByHash(blkHash string) (bi interface{}, err error)  { return }

//transaction
func (et *EthToken) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) {
	return
}
func (et *EthToken) Transfer(from, to map[string]decimal.Decimal) (txHash string, err error) { return }

//token
func (et *EthToken) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) { return }
func (et *EthToken) IsToken() bool                                                      { return false }

//others
func (et *EthToken) EstimateFee(map[string]interface{}) (fee decimal.Decimal, err error) { return }
