package solana_test

/*func TestNew(t *testing.T) {
	c := client.NewClient("https://api.devnet.solana.com/")
	ctx := context.Background()

	userPub := "B2KhBdBCcKWexFob3wrdcfbjaQ31kZ3r7mrQxaqNLVh9"
	// user    := types.AccountFromPrivateKeyBytes([]byte{0xa, 0x51, 0xfd, 0xbe, 0xde, 0x59, 0xb7, 0x1c, 0x2c, 0x9e, 0x56, 0x8a, 0xad, 0x9, 0x57, 0xc3, 0x19, 0x25, 0xfa, 0xca, 0x6f, 0x17, 0xe1, 0xec, 0x11, 0x5d, 0xd5, 0xad, 0x20, 0xd6, 0xe, 0xc2, 0x94, 0xeb, 0x96, 0xa, 0xe7, 0xd4, 0xec, 0x4c, 0x87, 0xb4, 0x34, 0x38, 0xd9, 0x73, 0xa6, 0x48, 0xaf, 0xbe, 0xa0, 0xf7, 0xa5, 0x52, 0x2e, 0x6b, 0xcf, 0x9f, 0xa7, 0xda, 0x78, 0x89, 0x9b, 0x10})

	r, err := c.GetBalance(ctx, userPub)
	require.NotNil(t, r)

	resp, err := c.GetConfirmedSignaturesForAddress(ctx, userPub, client.GetConfirmedSignaturesForAddressConfig{Limit: 25})
	require.NoError(t, err)
	require.NotNil(t, resp)

	var post, pre int64
	var signs []string

	for _, r := range resp {
		tx, err := c.GetConfirmedTransaction(ctx, r.Signature)
		require.NoError(t, err)
		require.NotNil(t, tx)

		signs = append(signs, r.Signature)

		for _, s := range tx.Meta.PostBalances {
			post += s
		}
		for _, p := range tx.Meta.PreBalances {
			pre += p
		}
	}

	status, err := c.GetSignatureStatuses(ctx, signs)
	require.NoError(t, err)
	require.NotNil(t, status)

	diff := pre - post

	tb, err := c.GetTokenAccountBalance(ctx, userPub, client.CommitmentFinalized)
	require.NoError(t, err)
	require.NotNil(t, tb)
	require.EqualValues(t, tb, diff)

}*/

//func TestSol(t *testing.T) {
//	feePayer := types.NewAccount()
//	issuer := types.NewAccount()
//	asset := types.NewAccount()
//
//	newUser := types.NewAccount()
//
//	ctx := context.Background()
//	sc := solana.New(client.DevnetRPCEndpoint, feePayer, asset, issuer)
//
//	assert.Equal(t, feePayer, sc.FeePayer)
//	assert.Equal(t, issuer, sc.Issuer)
//	assert.Equal(t, asset, sc.Asset)
//
//	for i := 0; i < 5; i++ {
//		tx, err := sc.RequestAirdrop(ctx, feePayer.PublicKey.ToBase58(), 10)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		log.Printf("airdrop for %s: %s", feePayer.PublicKey.ToBase58(), tx)
//		break
//	}
//
//	time.Sleep(time.Second * 15)
//
//	// convert account to asset
//	if tx, err := sc.CreateAsset(ctx); err != nil {
//		print(err)
//	} else {
//		log.Printf("convert account (%s) to asset: %s", asset.PublicKey.ToBase58(), tx)
//	}
//
//	time.Sleep(time.Second * 15)
//
//	// allow issuer account to hold asset
//	if tx, err := sc.InitAccountToUseAsset(
//		ctx,
//		sc.AccountFromPrivateKey(issuer.PrivateKey),
//	); err != nil {
//		print(err)
//	} else {
//		log.Printf("init issuer account (%s) to user asset: %s", issuer.PublicKey.ToBase58(), tx)
//	}
//
//	time.Sleep(time.Second * 15)
//
//	// issue tokens
//	if tx, err := sc.IssueAsset(
//		ctx,
//		sc.AccountFromPrivateKey(issuer.PrivateKey),
//		10000000,
//	); err != nil {
//		print(err)
//	} else {
//		log.Printf("issue asset (%s): %s", issuer.PublicKey.ToBase58(), tx)
//	}
//
//	time.Sleep(time.Second*10)
//
//	if info, err := sc.GetTokenAccountBalance(ctx, asset.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer: %+v", info)
//	}
//
//	// allow user account to hold asset
//	if tx, err := sc.InitAccountToUseAsset(ctx, newUser); err != nil {
//		log.Fatalln(err)
//	} else {
//		log.Println(tx)
//	}
//
//	time.Sleep(time.Second * 15)
//
//	if info, err := sc.GetTokenAccountBalance(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer token before send: %+v", info)
//	}
//	if info, err := sc.GetAccountBalanceSOL(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer sol before send: %+v", info)
//	}
//
//	time.Sleep(time.Second * 15)
//
//	// sends token
//	for i := 0; i < 5; i++ {
//		if tx, err := sc.SendAssets(ctx, issuer, newUser.PublicKey.ToBase58(), 10000); err != nil {
//			log.Println(err)
//			time.Sleep(time.Second * 10)
//		} else {
//			log.Println(tx)
//			break
//		}
//	}
//
//	time.Sleep(time.Second * 15)
//
//	if info, err := sc.GetTokenAccountBalance(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer token after send: %+v", info)
//	}
//	if info, err := sc.GetAccountBalanceSOL(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer sol after222 send: %+v", info)
//	}
//
//	time.Sleep(time.Second * 15)
//
//	for i := 0; i < 5; i++ {
//		if balance, err := sc.GetTokenAccountBalance(ctx, newUser.PublicKey.ToBase58()); err != nil {
//			log.Println(err)
//			time.Sleep(time.Second * 10)
//		} else {
//			log.Println("user token balance", balance)
//			break
//		}
//	}
//
//	if info, err := sc.GetAccountBalanceSOL(ctx, feePayer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("fee payer sol: %+v", info)
//	}
//	if info, err := sc.GetTokenAccountBalance(ctx, feePayer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("fee payer token: %+v", info)
//	}
//
//	if info, err := sc.GetTokenAccountBalance(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer token: %+v", info)
//	}
//	if info, err := sc.GetAccountBalanceSOL(ctx, issuer.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("issuer sol: %+v", info)
//	}
//
//	if info, err := sc.GetTokenAccountBalance(ctx, asset.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("asset token: %+v", info)
//	}
//	if info, err := sc.GetAccountBalanceSOL(ctx, asset.PublicKey.ToBase58()); err != nil {
//		log.Println(err)
//	} else {
//		log.Printf("asset sol: %+v", info)
//	}
//}
