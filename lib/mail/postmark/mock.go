//go:build mock_postmark

package postmark

import (
	"fmt"

	"github.com/golang/mock/gomock"

	lib_mail "github.com/SatorNetwork/sator-api/lib/mail"
	"github.com/SatorNetwork/sator-api/test/mock"
)

type stub struct{}

func (l *stub) Errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Errorf stub is called!"+format, args...))
}

func (l *stub) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Fatalf stub is called!"+format, args...))
}

func New(pc postmarkClient, cnf Config) lib_mail.Interface {
	m := mock.GetMockObject(mock.PostMarkProvider)
	if m == nil {
		m = lib_mail.NewMockInterface(gomock.NewController(&stub{}))
		m.(*lib_mail.MockInterface).ExpectSendVerificationCodeAny()
		m.(*lib_mail.MockInterface).ExpectSendResetPasswordCodeAny()
		m.(*lib_mail.MockInterface).ExpectSendDestroyAccountCodeAny()
		m.(*lib_mail.MockInterface).ExpectSendInvitationAny()
	}
	return m.(lib_mail.Interface)
}
