//go:build !mock_solana

package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/tokenprog"
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
	resp, err := c.solanaRpc.GetProgramAccountsWithConfig(ctx, common.TokenProgramID.ToBase58(), rpc.GetProgramAccountsConfig{
		Encoding: rpc.GetProgramAccountsConfigEncodingJsonParsed,
		Filters: []rpc.GetProgramAccountsConfigFilter{
			{
				DataSize: tokenprog.TokenAccountSize,
			},
			{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 32,
					Bytes:  walletAddr,
				},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get program accounts with config")
	}

	mintAddr := make([]string, 0)
	for _, account := range resp.Result {
		dataInJSON, err := json.Marshal(account.Account.Data)
		if err != nil {
			return nil, errors.Wrap(err, "can't marshal token account data")
		}

		var tokenAccountData TokenAccountData
		if err := json.Unmarshal(dataInJSON, &tokenAccountData); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal token account data")
		}

		if !isNFT(tokenAccountData.Parsed.Info.TokenAmount) {
			log.Printf("token %s is not nft\ndata: %+v\n", tokenAccountData.Parsed.Info.Mint, tokenAccountData.Parsed.Info.TokenAmount.Amount)
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
		return nil, err
	}

	var arweaveNFTMetadata lib_solana.ArweaveNFTMetadata
	if err := json.Unmarshal(body, &arweaveNFTMetadata); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal arweave nft metadata")
	}

	return &arweaveNFTMetadata, nil
}
