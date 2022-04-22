package coingecko

import "github.com/superoo7/go-gecko/v3/types"

//go:generate mockgen -destination=mock_client.go -package=coingecko github.com/SatorNetwork/sator-api/lib/coingecko Interface
type Interface interface {
	SimpleSinglePrice(id string, vsCurrency string) (*types.SimpleSinglePrice, error)
	SimplePrice(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error)
}
