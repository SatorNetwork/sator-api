package mail

import "github.com/golang/mock/gomock"

func (m *MockInterface) ExpectSendVerificationCodeAny() *gomock.Call {
	return m.EXPECT().
		SendVerificationCode(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}

func (m *MockInterface) ExpectSendResetPasswordCodeAny() *gomock.Call {
	return m.EXPECT().
		SendResetPasswordCode(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}

func (m *MockInterface) ExpectSendDestroyAccountCodeAny() *gomock.Call {
	return m.EXPECT().
		SendDestroyAccountCode(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}

func (m *MockInterface) ExpectSendInvitationAny() *gomock.Call {
	return m.EXPECT().
		SendInvitation(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
}
