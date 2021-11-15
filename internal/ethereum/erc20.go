package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/SatorNetwork/sator-api/internal/ethereum/erc20"
)

const (
	contractAddressHex = "0xA60288C54653aea211457613aF30B96CaA49fa24"
)

func (c *Client) TransferERC20(keyedTransactor *bind.TransactOpts, toAddr common.Address, amount uint64) error {
	conn, err := erc20.NewErc20(common.HexToAddress(contractAddressHex), &c.client)
	if err != nil {
		return err
	}

	_, err = conn.Transfer(keyedTransactor, toAddr, big.NewInt(int64(amount)))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) BalanceOf(ctx context.Context, address common.Address) (uint64, error) {
	conn, err := erc20.NewErc20(common.HexToAddress(contractAddressHex), &c.client)
	if err != nil {
		return 0, err
	}

	balance, err := conn.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}
