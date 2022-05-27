package gapi

import "errors"

// Predefined package errors
var (
	ErrCouldNotSignResponse          = errors.New("could not sign response")
	ErrCouldNotVerifySignature       = errors.New("could not verify signature")
	ErrNotAllNftsToCraftWereFound    = errors.New("not all nfts to craft were found")
	ErrNotEnoughNFTsToCraft          = errors.New("not enough nfts to craft")
	ErrNFTsToCraftHaveDifferentTypes = errors.New("nfts to craft have different types")
	ErrNFTTypeLegendCannotBeCrafted  = errors.New("nft type legend cannot be crafted")
	ErrCouldNotCraftNFT              = errors.New("could not craft nft")
	ErrNotEnoughElectricity          = errors.New("not enough electricity")
)
