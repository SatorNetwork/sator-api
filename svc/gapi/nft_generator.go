package gapi

import (
	"encoding/json"
	"fmt"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/dmitrymomot/random"
	"github.com/segmentio/ksuid"
)

// Generates NFT item from NFTPackInfo
// returns NFTInfo or error
func generateNFT(nftPack repository.UnityGameNftPack) (NFTInfo, error) {
	var dropChances DropChances
	if err := json.Unmarshal(nftPack.DropChances, &dropChances); err != nil {
		return NFTInfo{}, fmt.Errorf("failed to unmarshal drop chances: %w", err)
	}

	nftType := random.GetRandomMapItemWithProbabilities(dropChances.ToMap())

	return NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: getRandomNFTLevel(),
		NftType:  NFTType(nftType),
	}, nil
}

// Craft new NFT from users NFTs list
// returns NFTInfo or error
func craftNFT(nftsToCraft []repository.UnityGameNft) (*NFTInfo, error) {
	if len(nftsToCraft) < 2 {
		return nil, ErrNotEnoughNFTsToCraft
	}

	return &NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: getNextNFTLevel(nftsToCraft[0].MaxLevel),
		NftType:  getNextNFTType(NFTType(nftsToCraft[0].NftType)),
	}, nil
}
