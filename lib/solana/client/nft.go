//go:build !mock_solana

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/rpc"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
)

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

type TokenAccountData struct {
	Parsed struct {
		Info struct {
			IsNative    bool        `json:"isNative"`
			Mint        string      `json:"mint"`
			Owner       string      `json:"owner"`
			State       string      `json:"state"`
			TokenAmount TokenAmount `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

// GetNFTMintAddrs returns mint addrs which owned by walletAddr
func (c *Client) GetNFTMintAddrs(ctx context.Context, walletAddr string) ([]string, error) {
	resp, err := c.solanaRpc.GetTokenAccountsByOwnerWithConfig(ctx, walletAddr, rpc.GetTokenAccountsByOwnerConfigFilter{
		ProgramId: common.TokenProgramID.ToBase58(),
	}, rpc.GetTokenAccountsByOwnerConfig{
		Encoding: rpc.GetTokenAccountsByOwnerConfigEncodingJsonParsed,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get token accounts by owner")
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("get token accounts by owner error: %s", resp.Error.Message)
	}

	mintAddr := make([]string, 0)
	for _, account := range resp.Result.Value {
		dataInJSON, err := json.Marshal(account.Account.Data)
		if err != nil {
			return nil, errors.Wrap(err, "can't marshal token account data")
		}

		var tokenAccountData TokenAccountData
		if err := json.Unmarshal(dataInJSON, &tokenAccountData); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal token account data")
		}

		if !isNFT(tokenAccountData.Parsed.Info.TokenAmount) {
			continue
		}

		mintAddr = append(mintAddr, tokenAccountData.Parsed.Info.Mint)
	}

	return mintAddr, nil
}

func isNFT(t TokenAmount) bool {
	return t.Decimals == 0 && t.UiAmount == 1
}

// GetNFTsByWalletAddress returns NFTs which owned by walletAddr
func (c *Client) GetNFTsByWalletAddress(ctx context.Context, walletAddr string) ([]*lib_solana.ArweaveNFTMetadata, error) {
	mintAddrs, err := c.GetNFTMintAddrs(ctx, walletAddr)
	if err != nil {
		return nil, errors.Wrap(err, "can't get nft mint addrs")
	}

	nfts := make([]*lib_solana.ArweaveNFTMetadata, 0, len(mintAddrs))
	for _, mintAddr := range mintAddrs {
		nftMetadata, err := c.getNFTMetadata(mintAddr)
		if err != nil {
			return nil, errors.Wrap(err, "can't get nft metadata")
		}

		arweaveNFTMetadata, err := loadArweaveNFTMetadata(nftMetadata.Data.Uri)
		if err != nil {
			return nil, errors.Wrap(err, "can't load arweave nft metadata")
		}

		nfts = append(nfts, arweaveNFTMetadata)
	}

	return nfts, nil
}

// GetNFTMetadata returns metadata of NFT
func (c *Client) GetNFTMetadata(mintAddr string) (*lib_solana.ArweaveNFTMetadata, error) {
	nftMetadata, err := c.getNFTMetadata(mintAddr)
	if err != nil {
		return nil, errors.Wrap(err, "can't get nft metadata")
	}

	arweaveNFTMetadata, err := loadArweaveNFTMetadata(nftMetadata.Data.Uri)
	if err != nil {
		return nil, errors.Wrap(err, "can't load arweave nft metadata")
	}

	return arweaveNFTMetadata, nil
}

func (c *Client) getNFTMetadata(mintAddr string) (*tokenmeta.Metadata, error) {
	mint := common.PublicKeyFromString(mintAddr)
	metadataAccount, err := tokenmeta.GetTokenMetaPubkey(mint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get metadata account")
	}

	// get data which stored in metadataAccount
	accountInfo, err := c.solana.GetAccountInfo(context.Background(), metadataAccount.ToBase58())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account info")
	}

	// parse it
	metadata, err := tokenmeta.MetadataDeserialize(accountInfo.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse meta account")
	}

	return &metadata, nil
}

func loadArweaveNFTMetadata(uri string) (*lib_solana.ArweaveNFTMetadata, error) {
	httpClient := http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := httpClient.Get(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get uri %v", uri)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read body of uri %v", uri)
	}

	log.Printf("uri: %s\nbody: %v\n", uri, string(body))

	var arweaveNFTMetadata lib_solana.ArweaveNFTMetadata
	if err := json.Unmarshal(body, &arweaveNFTMetadata); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal arweave nft metadata")
	}

	return &arweaveNFTMetadata, nil
}
