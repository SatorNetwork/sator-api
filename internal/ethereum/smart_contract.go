package ethereum

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func MintNFT() error {
	client, err := ethclient.Dial("http://geth-dev.sator.io:8545")
	if err != nil {
		log.Fatal(err)
	}

	instance, err := NewEthereum(common.HexToAddress("0x06C46089EC98ed594aaec58f323F85D377755987"), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("c0764f8178607f6e285ea5cb5d4740d928f568387bd650068fdfd39dc960b3e6")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 4 for rinkeby, more chainID's here: https://chainlist.org/
	chainID := big.NewInt(4)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(21000) // in units
	auth.GasPrice = gasPrice

	tx, err := instance.Mint(auth, common.HexToAddress("0x6323324e4376a2babb1f9096a6eab66326fc60b8"), "testToken_123")
	if err != nil {
		log.Fatal(err)
	}
	print(tx)

	return nil
}
