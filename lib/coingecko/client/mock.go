//go:build mock_coingecko

package client

import (
	"fmt"

	"github.com/golang/mock/gomock"

	lib_coingecko "github.com/SatorNetwork/sator-api/lib/coingecko"
	"github.com/SatorNetwork/sator-api/test/mock"
)

type stub struct{}

func (l *stub) Errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Errorf stub is called!"+format, args...))
}

func (l *stub) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Fatalf stub is called!"+format, args...))
}

func NewCoingeckoClient() lib_coingecko.Interface {
	m := mock.GetMockObject(mock.CoingeckoProvider)
	if m == nil {
		m = lib_coingecko.NewMockInterface(gomock.NewController(&stub{}))
		m.(*lib_coingecko.MockInterface).ExpectSimplePriceAny()
	}
	return m.(lib_coingecko.Interface)
}
