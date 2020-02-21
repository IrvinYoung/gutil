package cryptocurrency

import (
	"context"
	"github.com/IrvinYoung/gutil/cryptocurrency/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"regexp"
)

type Ethereum struct {
	ctx context.Context
	c   *ethclient.Client
	t   *ERC20.ERC20

	isToken  bool
	Host     string
	Contract string
	dcm      int64
}

var (
	ReEthereumAccount = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

func InitEthereumClient(host string) (e *Ethereum, err error) {
	e = &Ethereum{Host: host, isToken: false, dcm: -1}
	e.ctx = context.Background()
	e.c, err = ethclient.DialContext(e.ctx, e.Host)
	return
}

func (e *Ethereum) Init() (err error) {
	e.isToken, e.dcm = false, -1
	e.ctx = context.Background()
	e.c, err = ethclient.DialContext(e.ctx, e.Host)
	return
}

func (e *Ethereum) IsToken() bool {
	return e.isToken
}

func (e *Ethereum) Close() {
	e.c.Close()
}

//basic
func (e *Ethereum) Name() string {
	if !e.isToken {
		return "Ethereum"
	}
	name, err := e.t.Name(&bind.CallOpts{})
	if err != nil {
		log.Println("get token name failed,", err)
	}
	return name
}

func (e *Ethereum) Symbol() string {
	if !e.isToken {
		return "eth"
	}
	symbol, err := e.t.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Println("get token symbol failed,", err)
	}
	return symbol
}
func (e *Ethereum) Decimal() int64 {
	if !e.isToken {
		e.dcm = 18
		return e.dcm
	}
	if e.dcm >= 0 {
		return e.dcm
	}
	d, err := e.t.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Println("get token decimal failed,", err)
		return 18
	}
	e.dcm = d.Int64()
	return e.dcm
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
	var amount *big.Int
	if !e.isToken {
		amount, err = e.c.BalanceAt(e.ctx, common.HexToAddress(addr), big.NewInt(int64(blkNum)))
	} else {
		amount, err = e.t.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr))
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	b = decimal.NewFromBigInt(amount, 0)
	return
}

//block
func (e *Ethereum) LastBlockNumber() (blkNum uint64, err error)             { return }
func (e *Ethereum) BlockByNumber(blkNum uint64) (bi interface{}, err error) { return }
func (e *Ethereum) BlockByHash(blkHash string) (bi interface{}, err error)  { return }

//transaction
func (e *Ethereum) TransactionsInBlocks(from, to uint64) (txs []*TransactionRecord, err error) { return }
func (e *Ethereum) Transfer(from, to map[string]decimal.Decimal) (txHash string, err error)    { return }

//token
func (e *Ethereum) TokenInstance(tokenInfo interface{}) (cc CryptoCurrency, err error) { return }

//others
func (e *Ethereum) EstimateFee(map[string]interface{}) (fee decimal.Decimal, err error) { return }
