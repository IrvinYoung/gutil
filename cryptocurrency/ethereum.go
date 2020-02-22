package cryptocurrency

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
	"math/big"
	"regexp"
)

type Ethereum struct {
	ctx context.Context
	c   *ethclient.Client
	t   *ERC20.ERC20

	Host string
}

var (
	ReEthereumAccount = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

func InitEthereumClient(host string) (e *Ethereum, err error) {
	e = &Ethereum{Host: host}
	e.ctx = context.Background()
	e.c, err = ethclient.DialContext(e.ctx, e.Host)
	return
}

func (e *Ethereum) Init() (err error) {
	e.ctx = context.Background()
	e.c, err = ethclient.DialContext(e.ctx, e.Host)
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
	b = ToDecimal(amount, e.Decimal())
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
func (e *Ethereum) BlockByNumber(blkNum uint64) (bi interface{}, err error) { return }
func (e *Ethereum) BlockByHash(blkHash string) (bi interface{}, err error)  { return }

//transaction
func (e *Ethereum) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) { return }
func (e *Ethereum) Transfer(from, to map[string]decimal.Decimal) (txHash string, err error)    { return }

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

func ToDecimal(ivalue interface{}, decimals int64) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)
	return result
}

func ToWei(iamount interface{}, decimals int64) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)
	wei := new(big.Int)
	wei.SetString(result.String(), 10)
	return wei
}
