//go:build !mock_solana

package client

import (
	"context"
	"encoding/base64"
	"log"
	"testing"
	"time"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

func TestScNew(t *testing.T) {
	t.Skip()

	c := New("https://api.devnet.solana.com/", Config{
		SystemProgram:  common.SystemProgramID.ToBase58(),
		SysvarRent:     common.SysVarRentPubkey.ToBase58(),
		SysvarClock:    common.SysVarClockPubkey.ToBase58(),
		SplToken:       common.TokenProgramID.ToBase58(),
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	})

	feePayerPrivate, err := base64.StdEncoding.DecodeString("MeFkg3Y/Ssa+CwfoZO6SvunvBDvxc/y/Jk/Ux/7G6F2vhaKpJyy/+5dPzj7iwO4hNlBxa3rtoRmeJzRLPPhZ5A==")
	if err != nil {
		log.Fatalf("feePayerPk base64 decoding error: %v", err)
	}

	IssuerPrivate, err := base64.StdEncoding.DecodeString("ICjjYoQ8RTqN8jc8Nn27Nda6miIA2xnvs4om1hh1TDn51/MOHeUiZeoIVb9Q4/csWYfLqBCOCMb6x5OUufGjBQ==")
	if err != nil {
		log.Fatalf("IssuerPk base64 decoding error: %v", err)
	}

	feePayerPub := common.PublicKeyFromString("CpAY2VpdxVK5kEQvoKCqNtaMomUMr8iVSXRPJFBZPZtf")
	issuerPub := common.PublicKeyFromString("HpHWCqBPRm7QCZDkR39WuoZp9xHF351T5U3AfHWXj8RA")
	asset := common.PublicKeyFromString("FBDfbe7CFXHHNzDpNBYf4Evcg5GKrThYNjk4wP2xwjwA")

	feePayer := types.Account{
		PublicKey:  feePayerPub,
		PrivateKey: feePayerPrivate,
	}

	issuer := types.Account{
		PublicKey:  issuerPub,
		PrivateKey: IssuerPrivate,
	}

	ctx := context.Background()
	wallet := types.NewAccount()

	tx, err := c.CreateAccountWithATA(ctx, asset.ToBase58(), wallet.PublicKey.ToBase58(), feePayer)
	if err != nil {
		log.Println(err.Error())
	}

	for i := 0; i < 5; i++ {
		tx, err = c.SendAssetsWithAutoDerive(ctx, asset.ToBase58(), feePayer, issuer, wallet.PublicKey.ToBase58(), 2, 0, false)
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
		log.Println(err.Error())
	}

	time.Sleep(5 * time.Second)

	tx, stakePool, err := c.InitializeStakePool(ctx, feePayer, issuer, asset)
	if err != nil {
		println(err.Error())
	}
	log.Println(tx)

	time.Sleep(time.Second * 5)

	for i := 0; i < 5; i++ {
		if tx, err = c.Stake(ctx, feePayer,
			wallet, stakePool.PublicKey, asset, 10, 1); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 20)
		} else {
			log.Println(tx)
			break
		}
	}

	time.Sleep(20 * time.Second)
	for i := 0; i < 5; i++ {
		if tx, err := c.Unstake(ctx, feePayer, wallet, stakePool.PublicKey, asset); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 20)
		} else {
			log.Println(tx)
			break
		}
	}
}
