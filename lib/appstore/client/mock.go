//go:build mock_appstore

package client

import (
	"fmt"

	"github.com/golang/mock/gomock"

	lib_appstore "github.com/SatorNetwork/sator-api/lib/appstore"
	"github.com/SatorNetwork/sator-api/test/mock"
)

type stub struct{}

func (l *stub) Errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Errorf stub is called!"+format, args...))
}

func (l *stub) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Fatalf stub is called!"+format, args...))
}

func New() lib_appstore.Interface {
	m := mock.GetMockObject(mock.AppStoreProvider)
	if m == nil {
		m = lib_appstore.NewMockInterface(gomock.NewController(&stub{}))
		m.(*lib_appstore.MockInterface).ExpectVerifyAny()
	}
	return m.(lib_appstore.Interface)
}
