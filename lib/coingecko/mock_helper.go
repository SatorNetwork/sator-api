package coingecko

import (
	"github.com/golang/mock/gomock"
)

func (m *MockInterface) ExpectSimplePriceAny() *gomock.Call {
	simplePriceCallback := func(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
		priceMap := map[string]map[string]float32{
			"solana":  {"usd": 1},
			"sator":   {"usd": 2},
			"arweave": {"usd": 3},
		}
		return &priceMap, nil
	}
	return m.EXPECT().
		SimplePrice([]string{"solana", "sator", "arweave"}, []string{"usd"}).
		DoAndReturn(simplePriceCallback).
		AnyTimes()
}
