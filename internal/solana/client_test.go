package solana_test

//func TestNew(t *testing.T) {
//	c := client.NewClient("https://api.devnet.solana.com/")
//	ctx := context.Background()
//
//	userPub := "B2KhBdBCcKWexFob3wrdcfbjaQ31kZ3r7mrQxaqNLVh9"
//	// user    := types.AccountFromPrivateKeyBytes([]byte{0xa, 0x51, 0xfd, 0xbe, 0xde, 0x59, 0xb7, 0x1c, 0x2c, 0x9e, 0x56, 0x8a, 0xad, 0x9, 0x57, 0xc3, 0x19, 0x25, 0xfa, 0xca, 0x6f, 0x17, 0xe1, 0xec, 0x11, 0x5d, 0xd5, 0xad, 0x20, 0xd6, 0xe, 0xc2, 0x94, 0xeb, 0x96, 0xa, 0xe7, 0xd4, 0xec, 0x4c, 0x87, 0xb4, 0x34, 0x38, 0xd9, 0x73, 0xa6, 0x48, 0xaf, 0xbe, 0xa0, 0xf7, 0xa5, 0x52, 0x2e, 0x6b, 0xcf, 0x9f, 0xa7, 0xda, 0x78, 0x89, 0x9b, 0x10})
//
//	r, err := c.GetBalance(ctx, userPub)
//	require.NotNil(t, r)
//
//	resp, err := c.GetConfirmedSignaturesForAddress(ctx, userPub, client.GetConfirmedSignaturesForAddressConfig{Limit: 25})
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	var post, pre int64
//	var signs []string
//
//	for _, r := range resp {
//		tx, err := c.GetConfirmedTransaction(ctx, r.Signature)
//		require.NoError(t, err)
//		require.NotNil(t, tx)
//
//		signs = append(signs, r.Signature)
//
//		for _, s := range tx.Meta.PostBalances {
//			post += s
//		}
//		for _, p := range tx.Meta.PreBalances {
//			pre += p
//		}
//	}
//
//	status, err := c.GetSignatureStatuses(ctx, signs)
//	require.NoError(t, err)
//	require.NotNil(t, status)
//
//	diff := pre - post
//
//	tb, err := c.GetTokenAccountBalance(ctx, userPub, client.CommitmentFinalized)
//	require.NoError(t, err)
//	require.NotNil(t, tb)
//	require.EqualValues(t, tb, diff)
//
//}

