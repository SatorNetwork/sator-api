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
func generateNFT(nftPack repository.UnityGameNftPack) (*NFTInfo, error) {
	var dropChances DropChances
	if err := json.Unmarshal(nftPack.DropChances, &dropChances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drop chances: %w", err)
	}

	nftType := random.GetRandomMapItemWithPrecent(dropChances.ToMap())
	nt := NFTType(nftType)
	if !nt.IsValid() {
		return nil, fmt.Errorf("invalid nft type: %s", nt)
	}

	return &NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: getNFTLevelByType(nt),
		NftType:  nt,
	}, nil
}

// Craft new NFT from users NFTs list
// returns NFTInfo or error
func craftNFT(nftsToCraft []repository.UnityGameNft) (*NFTInfo, error) {
	if len(nftsToCraft) < 2 {
		return nil, ErrNotEnoughNFTsToCraft
	}

	var nftType string
	for _, nft := range nftsToCraft {
		if nftType != "" && nftType != nft.NftType {
			return nil, ErrNFTsToCraftHaveDifferentTypes
		}

		if nft.NftType == NFTTypeLegend.String() {
			return nil, ErrNFTTypeLegendCannotBeCrafted
		}

		nftType = nft.NftType
	}

	nt := getNextNFTType(NFTType(nftType))

	return &NFTInfo{
		ID:       ksuid.New().String(),
		MaxLevel: getNFTLevelByType(nt),
		NftType:  nt,
	}, nil
}
