package solana

import (
	"github.com/golang/mock/gomock"
	"github.com/portto/solana-go-sdk/types"
)

func (m *MockInterface) ExpectCheckPrivateKeyAny() *gomock.Call {
	return m.EXPECT().
		CheckPrivateKey(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}

func (m *MockInterface) ExpectNewAccountAny() *gomock.Call {
	return m.EXPECT().
		NewAccount().
		Return(types.NewAccount()).
		AnyTimes()
}

func (m *MockInterface) ExpectAccountFromPrivateKeyBytesAny() *gomock.Call {
	return m.EXPECT().
		AccountFromPrivateKeyBytes(gomock.Any()).
		DoAndReturn(func(pk []byte) (types.Account, error) {
			return types.AccountFromBytes(pk)
		}).
		AnyTimes()
}

func (m *MockInterface) ExpectTransactionDeserializeAny() *gomock.Call {
	return m.EXPECT().
		TransactionDeserialize(gomock.Any()).
		Return(types.Transaction{}, nil).
		AnyTimes()
}

func (m *MockInterface) ExpectSerializeTxMessageAny() *gomock.Call {
	return m.EXPECT().
		SerializeTxMessage(gomock.Any()).
		Return([]byte{}, nil).
		AnyTimes()
}
