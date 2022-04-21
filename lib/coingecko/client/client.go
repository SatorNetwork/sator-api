//go:build !mock_coingecko

package client

import (
	"net/http"
	"time"

	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"

	lib_coingecko "github.com/SatorNetwork/sator-api/lib/coingecko"
)

type coingeckoClient struct {
	client *coingecko.Client
}

func NewCoingeckoClient() lib_coingecko.Interface {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	client := coingecko.NewClient(httpClient)

	return &coingeckoClient{
		client: client,
	}
}

func (cg *coingeckoClient) SimpleSinglePrice(id string, vsCurrency string) (*types.SimpleSinglePrice, error) {
	return cg.client.SimpleSinglePrice(id, vsCurrency)
}

func (cg *coingeckoClient) SimplePrice(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
	return cg.client.SimplePrice(ids, vsCurrencies)
}
