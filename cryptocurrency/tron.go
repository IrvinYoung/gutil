package cryptocurrency

import (
	"context"
	"encoding/hex"
	"github.com/IrvinYoung/gutil/cryptocurrency/tron_lib"
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
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	//private key
	privateKeyData := crypto.FromECDSA(privateKeyECDSA)
	//priv = hexutil.Encode(privateKeyData)
	priv = hex.EncodeToString(privateKeyData) //without "0x"
	println(priv)
	//address
	addr = tron_lib.EncodeCheck(tron_lib.PubkeyToAddressBytes(privateKeyECDSA.PublicKey))
	//encrypt private key
	priv, err = encryptPrivKey(password, salt, priv)
	return
}

func (t *Tron) IsValidAccount(addr string) bool {
	if len(addr) != 34 {
		return false
	}
	if addr[0:1] != "T" {
		return false
	}
	_, err := tron_lib.DecodeCheck(addr)
	if err != nil {
		return false
	}
	return true
}

