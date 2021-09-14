package solana_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/portto/solana-go-sdk/types"
)

var (
	//feePayer = types.AccountFromPrivateKeyBytes([]byte{0xc2, 0x30, 0x16, 0x0, 0x95, 0x1b, 0xf8, 0x86, 0xf8, 0x71, 0x31, 0xab, 0x7d, 0x9d, 0x3b, 0x9d, 0x74, 0x6, 0x8d, 0xa6, 0xe1, 0xf0, 0x3, 0xd7, 0xdb, 0x26, 0xca, 0x5d, 0x98, 0x32, 0x2e, 0x4b, 0x35, 0x4, 0x1, 0x3b, 0xf, 0xdc, 0xe0, 0x52, 0x7e, 0x1c, 0x1f, 0xfc, 0x96, 0x68, 0x5f, 0xdc, 0x1d, 0xdd, 0x26, 0x7, 0xbf, 0x33, 0x1b, 0x1b, 0x84, 0xef, 0xf8, 0xd4, 0xec, 0x7d, 0xb7, 0xa6})

	feePayer = types.AccountFromPrivateKeyBytes([]byte{115, 91, 202, 172, 215, 254, 239, 102, 127, 239, 39, 117, 165, 14, 239, 60, 242, 138, 216, 4, 183, 230, 36, 122, 133, 128, 12, 201, 176, 200, 144, 182, 17, 64, 8, 222, 37, 225, 40, 90, 140, 94, 207, 194, 215, 172, 41, 156, 184, 231, 78, 111, 144, 102, 2, 211, 156, 35, 90, 19, 91, 13, 43, 209})

	issuer = types.AccountFromPrivateKeyBytes([]byte{0x7d, 0x36, 0x17, 0xd5, 0x2c, 0xc8, 0x64, 0xf2, 0x9a, 0x39, 0x2f, 0x8b, 0xb6, 0x40, 0x4e, 0xf9, 0xcd, 0x4c, 0x85, 0xa8, 0x9a, 0xbe, 0x3c, 0xfe, 0xa9, 0xe1, 0xad, 0xbc, 0xb5, 0x40, 0x2a, 0xf9, 0x1, 0xb6, 0x4b, 0x6f, 0x7e, 0x76, 0xc3, 0x3d, 0x4b, 0xf6, 0xcf, 0xc6, 0xb4, 0x6, 0xd8, 0x1f, 0xcf, 0x96, 0xe1, 0x67, 0x5a, 0xdf, 0xd3, 0x22, 0xbf, 0xe2, 0x8a, 0xa6, 0x92, 0xa, 0xee, 0x2f})

	asset = types.AccountFromPrivateKeyBytes([]byte{0x2b, 0xe7, 0x8c, 0x5, 0xbd, 0x7f, 0x6f, 0x7a, 0xb4, 0xd6, 0x68, 0x7a, 0xfa, 0xf3, 0xd6, 0x14, 0x9c, 0xce, 0x9a, 0xff, 0x72, 0x6a, 0x9, 0x40, 0x52, 0x16, 0x54, 0xe7, 0xe5, 0x75, 0xe0, 0x15, 0xc1, 0xc7, 0x6b, 0x43, 0x40, 0xe9, 0xdf, 0xc3, 0x9, 0x8a, 0x4f, 0xbd, 0x30, 0x99, 0xc4, 0x5d, 0x64, 0xcd, 0x43, 0xf5, 0xdf, 0x82, 0xf4, 0xc6, 0x4b, 0x6c, 0x5, 0x1c, 0xdc, 0xbc, 0x45, 0xd})
	//asset = common.PublicKeyFromString("13kBuVtxUT7CeddDgHfe61x3YdpBWTCKeB2Zg2LC4dab")
)

func TestScNew(t *testing.T) {
	c := solana.New("https://api.devnet.solana.com/")
	ctx := context.Background()

	_, err := c.RequestAirdrop(ctx, feePayer.PublicKey.ToBase58(), 10)
	time.Sleep(3 * time.Second)
	_, err = c.RequestAirdrop(ctx, issuer.PublicKey.ToBase58(), 10)
	time.Sleep(3 * time.Second)

	tx, stakePool, stakeAuthority, stakeTokenAccountPool, err := c.InitializeStakePool(ctx, feePayer, asset.PublicKey)
	println(tx)
	println(err)

	// wallet := c.PublicKeyFromString("7Cw8GWHV2EgbZto7v4prxi6LytqV237tCSVvBm8Z6WMQ")
	wallet := types.NewAccount()
	_, err = c.RequestAirdrop(ctx, wallet.PublicKey.ToBase58(), 10)
	time.Sleep(5 * time.Second)

	tx, err = c.IssueAsset(ctx, feePayer, issuer, asset, issuer, 20000)
	println(err)

	_, err = c.InitAccountToUseAsset(ctx, feePayer, issuer, asset, wallet)
	println(err)
	time.Sleep(10 * time.Second)
	for i := 0; i < 5; i++ {
		if tx, err := c.SendAssets(ctx, feePayer, issuer, asset, issuer, wallet.PublicKey.ToBase58(), 1000); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
		} else {
			log.Println(tx)
			break
		}
	}
	println(tx)
	time.Sleep(5 * time.Second)

	tx, _, err = c.Stake(ctx, feePayer, stakePool.PublicKey, stakeAuthority, wallet.PublicKey, stakeTokenAccountPool, 1000, 100)
	println(tx)
	println(err)
}

func TestScStake(t *testing.T) {
	c := solana.New("https://api.devnet.solana.com/")
	ctx := context.Background()

	stakePool := c.PublicKeyFromString("2c8X3S9PjENeU4JzD3A7AehTH9dhQGtVXnTc5VhNuhYj")
	stakeAuthority := c.PublicKeyFromString("A8DKHrWwoDjUj4TpvpKWY6vEyaqv6y9HSsN6kHpSwYyx")
	wallet := c.PublicKeyFromString("7Cw8GWHV2EgbZto7v4prxi6LytqV237tCSVvBm8Z6WMQ")
	stakeTokenAccountPool := c.PublicKeyFromString("F2DbdxnEieiM3BTDAaQyZWN5jH9iDdXe2V2rYvY79jCe")

	tx, acc, err := c.Stake(ctx, feePayer, stakePool, stakeAuthority, wallet, stakeTokenAccountPool, 470000, 1000)
	println(tx)
	println(err)

	tx, err := c.Unstake(ctx, feePayer, stakePool, wallet, stakeTokenAccountPool, acc, stakeAuthority)
	println(tx)
	println(err)
}
