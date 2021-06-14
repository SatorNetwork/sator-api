package repository

const (
	VerifyConfirmAccount = (iota + 1) << 2
	VerifyChangeEmail
	VerifyResetPassword
	VerifyDestroyAccount
)
