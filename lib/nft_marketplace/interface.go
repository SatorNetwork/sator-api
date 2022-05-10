package nft_marketplace

//go:generate mockgen -destination=mock_client.go -package=nft_marketplace github.com/SatorNetwork/sator-api/lib/nft_marketplace Interface
type (
	Interface interface {
		PrepareBuyTx(req *PrepareBuyTxRequest) (*PrepareBuyTxResponse, error)
		SendPreparedBuyTx(req *SendPreparedBuyTxRequest) (*SendPreparedBuyTxResponse, error)
	}

	PrepareBuyTxRequest struct {
		MintAddress      string `json:"mint_address,omitempty"`
		ChargeTokensFrom string `json:"charge_tokens_from,omitempty"`
	}

	PrepareBuyTxResponse struct {
		EncodedTx       string `json:"encoded_tx,omitempty"`
		EncodedMessage  string `json:"encoded_message,omitempty"`
		PreparedBuyTxId string `json:"prepared_buy_tx_id,omitempty"`
	}

	SendPreparedBuyTxRequest struct {
		PreparedBuyTxId string `json:"prepared_buy_tx_id,omitempty"`
		BuyerSignature  string `json:"buyer_signature,omitempty"`
	}

	SendPreparedBuyTxResponse struct {
		TxHash string `json:"tx_hash,omitempty"`
	}
)
