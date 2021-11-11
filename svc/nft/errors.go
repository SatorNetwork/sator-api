package nft

import "errors"

// Predefined package errors
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrAlreadySold      = errors.New("this nft sold out")
	ErrAlreadyBought    = errors.New("you already bought this NFT")
	ErrAlreadyMinted    = errors.New("this nft minted")
)
