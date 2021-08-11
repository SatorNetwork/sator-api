package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client ethclient.Client
}

func NewClient() (*Client, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		return &Client{}, err
	}

	return &Client{client: *client}, nil
}

func (c *Client) CreateAccount() (Wallet, error) {
	return CreateWallet()
}