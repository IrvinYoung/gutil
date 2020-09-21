package cryptocurrency

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

type Tron struct {
	ctx     context.Context
	//c       *ethclient.Client
	//t       *ERC20.ERC20
	//chainID *big.Int

	Host string
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

	return
}