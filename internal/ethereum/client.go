package ethereum

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO(evg): root ethereum address should be configurable
const RootEthAddressHex = "0x06Ea300809363d22cD1BbFE78D441C2608E3b496"

type Client struct {
	client ethclient.Client
}

func NewClient() (*Client, error) {
	//client, err := ethclient.Dial("https://mainnet.infura.io")
	//if err != nil {
	//	return &Client{}, err
	//}

	// TODO(evg): eth's url should be configurable
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/417ac592054c4b7584bfd9c6aa8dacf0")
	if err != nil {
		return &Client{}, err
	}

	return &Client{client: *client}, nil
}

func (c *Client) CreateAccount() (Wallet, error) {
	return CreateWallet()
}

func (c *Client) NewKeyedTransactor(privateKey *ecdsa.PrivateKey, address ethereum_common.Address) (*bind.TransactOpts, error) {
	nonce, err := c.client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(3000000) // in units
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}
