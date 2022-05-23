package gapi

import (
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/segmentio/ksuid"
)

// Generates NFT item from NFTPackInfo
// returns NFTInfo or error
func generateNFT(nftPack repository.UnityGameNftPack) (NFTInfo, error) {
	rand.Seed(time.Now().Unix())

	return NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: rand.Intn(3) + 1,
		NftType:  nftTypesSlice[rand.Intn(len(nftTypesSlice))],
	}, nil
}

// Craft new NFT from users NFTs list
// returns NFTInfo or error
func craftNFT(nftsToCraft []repository.UnityGameNft) (*NFTInfo, error) {
	rand.Seed(time.Now().Unix())

	return &NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: rand.Intn(3) + 1,
		NftType:  nftTypesSlice[rand.Intn(len(nftTypesSlice))],
	}, nil
}
