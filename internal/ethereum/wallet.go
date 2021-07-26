package ethereum

import (
	"crypto"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/zeebo/errs"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey `json:"private_key"`
	PublicKey  crypto.PublicKey  `json:"public_key"`
	Address    string            `json:"address"`
}

func CreateWallet() (Wallet, error) {
	privateKey, err := crypto2.GenerateKey()
	if err != nil {
		return Wallet{}, err
	}

	privateKeyBytes := crypto2.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Wallet{}, errs.New("internal error")
	}

	address := crypto2.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	return Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}
