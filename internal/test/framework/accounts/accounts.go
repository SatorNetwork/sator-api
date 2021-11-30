package accounts

import (
	"encoding/base64"

	"github.com/portto/solana-go-sdk/types"
)

const (
	feePayerPrivateKey    = "tg3BEHU1lH24lo9JccmqLL13DLOzLMptxh0aa3NXJUtL4PVdkvwOmbpCqMTFG7a8CJles911d0uu7SYeuck8Uw=="
	tokenHolderPrivateKey = "I52q0J0qsUY2NLTSScSKre1lH6XZRu69FGS0pa3xypsNYtRHIr9ICfw0SXUd1Vcr0sf3tqQuG3whne/UvJfBNQ=="
	assetPrivateKey       = "bHTM9fYUQAjX6IsbMYAEbhR0aUpL+3s28GC1N5nyFc8sJHK42tvmt+4FWvuaw6csfDGI3CoZc8D48o5eo9RKUQ=="
)

func GetFeePayer() types.Account {
	feePayerPrivateKeyBytes, err := base64.StdEncoding.DecodeString(feePayerPrivateKey)
	if err != nil {
		panic(err)
	}
	feePayer := types.AccountFromPrivateKeyBytes(feePayerPrivateKeyBytes)
	return feePayer
}

func GetTokenHolder() types.Account {
	tokenHolderPrivateKeyBytes, err := base64.StdEncoding.DecodeString(tokenHolderPrivateKey)
	if err != nil {
		panic(err)
	}
	tokenHolder := types.AccountFromPrivateKeyBytes(tokenHolderPrivateKeyBytes)
	return tokenHolder
}

func GetAsset() types.Account {
	assetPrivateKeyBytes, err := base64.StdEncoding.DecodeString(assetPrivateKey)
	if err != nil {
		panic(err)
	}
	asset := types.AccountFromPrivateKeyBytes(assetPrivateKeyBytes)
	return asset
}

func GetAccounts() (types.Account, types.Account, types.Account) {
	return GetFeePayer(), GetTokenHolder(), GetAsset()
}
