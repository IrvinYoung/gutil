package cryptocurrency

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type Ethereum struct {
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

func (e *Ethereum) IsValidAccount(addr string) bool                      { return false }
func (e *Ethereum) BalanceOf(addr string) (b decimal.Decimal, err error) { return }

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
