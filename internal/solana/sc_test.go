package solana_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/SatorNetwork/sator-api/internal/solana"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

var (
	feePayer = types.AccountFromPrivateKeyBytes([]byte{0xc2, 0x30, 0x16, 0x0, 0x95, 0x1b, 0xf8, 0x86, 0xf8, 0x71, 0x31, 0xab, 0x7d, 0x9d, 0x3b, 0x9d, 0x74, 0x6, 0x8d, 0xa6, 0xe1, 0xf0, 0x3, 0xd7, 0xdb, 0x26, 0xca, 0x5d, 0x98, 0x32, 0x2e, 0x4b, 0x35, 0x4, 0x1, 0x3b, 0xf, 0xdc, 0xe0, 0x52, 0x7e, 0x1c, 0x1f, 0xfc, 0x96, 0x68, 0x5f, 0xdc, 0x1d, 0xdd, 0x26, 0x7, 0xbf, 0x33, 0x1b, 0x1b, 0x84, 0xef, 0xf8, 0xd4, 0xec, 0x7d, 0xb7, 0xa6})
)

func TestScNew(t *testing.T) {
	c := solana.New("https://api.devnet.solana.com/")
	ctx := context.Background()

	wallet := common.PublicKeyFromString("7uWo5zDLDCyARCthtqsKhxGUMYfaViYUmK9rFMpZdJgS")

	payer := types.AccountFromPrivateKeyBytes([]byte{115, 91, 202, 172, 215, 254, 239, 102, 127, 239, 39, 117, 165, 14, 239, 60, 242, 138, 216, 4, 183, 230, 36, 122, 133, 128, 12, 201, 176, 200, 144, 182, 17, 64, 8, 222, 37, 225, 40, 90, 140, 94, 207, 194, 215, 172, 41, 156, 184, 231, 78, 111, 144, 102, 2, 211, 156, 35, 90, 19, 91, 13, 43, 209})
	println(payer.PublicKey.ToBase58())

	testAsset := types.AccountFromPrivateKeyBytes([]byte{213, 192, 236, 172, 129, 236, 157, 105, 169, 136, 46, 123, 109, 101, 48, 172, 124, 140, 128, 105, 10, 96, 229, 160, 116, 186, 58, 152, 181, 244, 123, 125, 162, 253, 106, 157, 104, 123, 65, 211, 209, 132, 130, 73, 185, 218, 92, 21, 65, 183, 177, 123, 72, 83, 37, 76, 144, 180, 119, 107, 90, 151, 97, 183})

	tx, stakePool, err := c.InitializeStakePool(ctx, payer, testAsset.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	println("tx = ")
	println(tx)

	time.Sleep(time.Second * 20)

	for i := 0; i < 5; i++ {
		if tx, err := c.Stake(ctx, payer, stakePool.PublicKey, wallet, 100, 100); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 20)
		} else {
			log.Println(tx)
			break
		}
	}

	time.Sleep(200 * time.Second)
	tx, err = c.Unstake(ctx, payer, stakePool.PublicKey, wallet)
	println(tx)
	println(err)
}

func TestScStake(t *testing.T) {
	c := solana.New("https://api.devnet.solana.com/")
	ctx := context.Background()

	stakePool := c.PublicKeyFromString("2c8X3S9PjENeU4JzD3A7AehTH9dhQGtVXnTc5VhNuhYj")
	wallet := c.PublicKeyFromString("7Cw8GWHV2EgbZto7v4prxi6LytqV237tCSVvBm8Z6WMQ")

	tx, err := c.Unstake(ctx, feePayer, stakePool, wallet)
	println(tx)
	println(err)
}
