package appstore

import "github.com/golang/mock/gomock"

func (m *MockInterface) ExpectVerifyAny() *gomock.Call {
	return m.EXPECT().
		Verify(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}
