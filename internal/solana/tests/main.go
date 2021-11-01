package main

import (
	"context"
	"log"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/types"
)

var (
	feePayerPub = "4Zx9h7C47oBB3duhXGn2GYGGtvKVRpuYTqNsoT9A94Ds"
	feePayer    = types.AccountFromPrivateKeyBytes([]byte{0xc2, 0x30, 0x16, 0x0, 0x95, 0x1b, 0xf8, 0x86, 0xf8, 0x71, 0x31, 0xab, 0x7d, 0x9d, 0x3b, 0x9d, 0x74, 0x6, 0x8d, 0xa6, 0xe1, 0xf0, 0x3, 0xd7, 0xdb, 0x26, 0xca, 0x5d, 0x98, 0x32, 0x2e, 0x4b, 0x35, 0x4, 0x1, 0x3b, 0xf, 0xdc, 0xe0, 0x52, 0x7e, 0x1c, 0x1f, 0xfc, 0x96, 0x68, 0x5f, 0xdc, 0x1d, 0xdd, 0x26, 0x7, 0xbf, 0x33, 0x1b, 0x1b, 0x84, 0xef, 0xf8, 0xd4, 0xec, 0x7d, 0xb7, 0xa6})

	issuerPub = "7gdYR2pRwM61cea9qmKF8fEXsgRa77KYtV5xh2o4xEA"
	issuer    = types.AccountFromPrivateKeyBytes([]byte{0x7d, 0x36, 0x17, 0xd5, 0x2c, 0xc8, 0x64, 0xf2, 0x9a, 0x39, 0x2f, 0x8b, 0xb6, 0x40, 0x4e, 0xf9, 0xcd, 0x4c, 0x85, 0xa8, 0x9a, 0xbe, 0x3c, 0xfe, 0xa9, 0xe1, 0xad, 0xbc, 0xb5, 0x40, 0x2a, 0xf9, 0x1, 0xb6, 0x4b, 0x6f, 0x7e, 0x76, 0xc3, 0x3d, 0x4b, 0xf6, 0xcf, 0xc6, 0xb4, 0x6, 0xd8, 0x1f, 0xcf, 0x96, 0xe1, 0x67, 0x5a, 0xdf, 0xd3, 0x22, 0xbf, 0xe2, 0x8a, 0xa6, 0x92, 0xa, 0xee, 0x2f})

	assetPub = "E3S44kKvw4ssUV71oLSvcNfKH77XwwFNUz6XHEjkRwec"
	asset    = types.AccountFromPrivateKeyBytes([]byte{0x2b, 0xe7, 0x8c, 0x5, 0xbd, 0x7f, 0x6f, 0x7a, 0xb4, 0xd6, 0x68, 0x7a, 0xfa, 0xf3, 0xd6, 0x14, 0x9c, 0xce, 0x9a, 0xff, 0x72, 0x6a, 0x9, 0x40, 0x52, 0x16, 0x54, 0xe7, 0xe5, 0x75, 0xe0, 0x15, 0xc1, 0xc7, 0x6b, 0x43, 0x40, 0xe9, 0xdf, 0xc3, 0x9, 0x8a, 0x4f, 0xbd, 0x30, 0x99, 0xc4, 0x5d, 0x64, 0xcd, 0x43, 0xf5, 0xdf, 0x82, 0xf4, 0xc6, 0x4b, 0x6c, 0x5, 0x1c, 0xdc, 0xbc, 0x45, 0xd})

	userPub = "B2KhBdBCcKWexFob3wrdcfbjaQ31kZ3r7mrQxaqNLVh9"
	user    = types.AccountFromPrivateKeyBytes([]byte{0xa, 0x51, 0xfd, 0xbe, 0xde, 0x59, 0xb7, 0x1c, 0x2c, 0x9e, 0x56, 0x8a, 0xad, 0x9, 0x57, 0xc3, 0x19, 0x25, 0xfa, 0xca, 0x6f, 0x17, 0xe1, 0xec, 0x11, 0x5d, 0xd5, 0xad, 0x20, 0xd6, 0xe, 0xc2, 0x94, 0xeb, 0x96, 0xa, 0xe7, 0xd4, 0xec, 0x4c, 0x87, 0xb4, 0x34, 0x38, 0xd9, 0x73, 0xa6, 0x48, 0xaf, 0xbe, 0xa0, 0xf7, 0xa5, 0x52, 0x2e, 0x6b, 0xcf, 0x9f, 0xa7, 0xda, 0x78, 0x89, 0x9b, 0x10})
)

func main() {
	ctx := context.Background()
	sc := solana.New(client.DevnetRPCEndpoint, solana.Config{
		SystemProgram:   "11111111111111111111111111111111",
		SysvarRent:      "SysvarRent111111111111111111111111111111111",
		SysvarClock:     "SysvarC1ock11111111111111111111111111111111",
		SplToken:        "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		StakeProgramID:  "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
		RewardProgramID: "DajevvE6uo5HtST4EDguRUcbdEMNKNcLWjjNowMRQvZ1",
	})

	txList, err := sc.GetTransactions(ctx, "3Z6t2topTRBVeQjbLf8mExLnEvxUz38rMkyiNSnwMkrj")
	if err != nil {
		log.Fatalln(err)
	}

	// log.Printf("txList: %+v\n\n", txList)

	for _, tx := range txList {
		log.Printf("tx: %+v\n\n", tx)
	}

	//acc := sc.NewAccount()
	//log.Printf("account pub key: %#v", acc.PublicKey.ToBase58())
	//log.Printf("account private key: %#v", acc.PrivateKey)
	//
	//_, err = sc.RequestAirdrop(ctx, acc.PublicKey.ToBase58(), 10)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//// convert account to asset
	//_, err = sc.CreateAsset(ctx, feePayer, issuer, asset)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//// allow issuer account to hold asset
	//if tx, err := sc.InitAccountToUseAsset(ctx, feePayer, issuer, asset, issuer); err != nil {
	//	log.Fatalln(err)
	//} else {
	//	log.Println(tx)
	//}
	//
	//// issue tokens
	//if tx, err := sc.IssueAsset(ctx, feePayer, issuer, asset, issuer, 1200); err != nil {
	//	log.Fatalln(err)
	//} else {
	//	log.Println(tx)
	//}
	//
	//// allow user account to hold asset
	//if tx, err := sc.InitAccountToUseAsset(ctx, feePayer, issuer, asset, acc); err != nil {
	//	log.Fatalln(err)
	//} else {
	//	log.Println(tx)
	//}
	//
	//for i := 0; i < 5; i++ {
	//	if balance, err := sc.GetTokenAccountBalance(ctx, acc.PublicKey.ToBase58()); err != nil {
	//		log.Println(err)
	//		time.Sleep(time.Second * 10)
	//	} else {
	//		log.Println("init balance", balance)
	//		break
	//	}
	//}
	//
	//// sends token
	//for i := 0; i < 5; i++ {
	//	if tx, err := sc.SendAssets(ctx, feePayer, issuer, asset, issuer, user.PublicKey.ToBase58(), 1000); err != nil {
	//		log.Println(err)
	//		time.Sleep(time.Second * 10)
	//	} else {
	//		log.Println(tx)
	//		break
	//	}
	//}
	//
	//time.Sleep(time.Second * 20)
	//
	//for i := 0; i < 5; i++ {
	//	if balance, err := sc.GetTokenAccountBalance(ctx, acc.PublicKey.ToBase58()); err != nil {
	//		log.Println(err)
	//		time.Sleep(time.Second * 10)
	//	} else {
	//		log.Println("balance", balance)
	//		break
	//	}
	//}
	//
	//if info, err := sc.GetAccountBalanceSOL(ctx, feePayerPub); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("fee payer: %+v", info)
	//}
	//
	//if info, err := sc.GetTokenAccountBalance(ctx, issuerPub); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("issuer: %+v", info)
	//}
	//
	//if info, err := sc.GetTokenAccountBalance(ctx, assetPub); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("asset: %+v", info)
	//}
	//
	//if info, err := sc.GetTokenAccountBalance(ctx, userPub); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("user: %+v", info)
	//}
}
