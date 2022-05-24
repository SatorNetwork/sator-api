// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SatorNetwork/sator-api/lib/solana (interfaces: Interface)

// Package solana is a generated GoMock package.
package solana

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/portto/solana-go-sdk/common"
	types "github.com/portto/solana-go-sdk/types"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// AccountFromPrivateKeyBytes mocks base method.
func (m *MockInterface) AccountFromPrivateKeyBytes(arg0 []byte) (types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountFromPrivateKeyBytes", arg0)
	ret0, _ := ret[0].(types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountFromPrivateKeyBytes indicates an expected call of AccountFromPrivateKeyBytes.
func (mr *MockInterfaceMockRecorder) AccountFromPrivateKeyBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountFromPrivateKeyBytes", reflect.TypeOf((*MockInterface)(nil).AccountFromPrivateKeyBytes), arg0)
}

// CheckPrivateKey mocks base method.
func (m *MockInterface) CheckPrivateKey(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPrivateKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckPrivateKey indicates an expected call of CheckPrivateKey.
func (mr *MockInterfaceMockRecorder) CheckPrivateKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPrivateKey", reflect.TypeOf((*MockInterface)(nil).CheckPrivateKey), arg0, arg1)
}

// CreateAccountWithATA mocks base method.
func (m *MockInterface) CreateAccountWithATA(arg0 context.Context, arg1, arg2 string, arg3 types.Account) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccountWithATA", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccountWithATA indicates an expected call of CreateAccountWithATA.
func (mr *MockInterfaceMockRecorder) CreateAccountWithATA(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccountWithATA", reflect.TypeOf((*MockInterface)(nil).CreateAccountWithATA), arg0, arg1, arg2, arg3)
}

// CreateAsset mocks base method.
func (m *MockInterface) CreateAsset(arg0 context.Context, arg1, arg2, arg3 types.Account) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAsset", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAsset indicates an expected call of CreateAsset.
func (mr *MockInterfaceMockRecorder) CreateAsset(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAsset", reflect.TypeOf((*MockInterface)(nil).CreateAsset), arg0, arg1, arg2, arg3)
}

// Endpoint mocks base method.
func (m *MockInterface) Endpoint() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Endpoint")
	ret0, _ := ret[0].(string)
	return ret0
}

// Endpoint indicates an expected call of Endpoint.
func (mr *MockInterfaceMockRecorder) Endpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Endpoint", reflect.TypeOf((*MockInterface)(nil).Endpoint))
}

// FeeAccumulatorAddress mocks base method.
func (m *MockInterface) FeeAccumulatorAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FeeAccumulatorAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// FeeAccumulatorAddress indicates an expected call of FeeAccumulatorAddress.
func (mr *MockInterfaceMockRecorder) FeeAccumulatorAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FeeAccumulatorAddress", reflect.TypeOf((*MockInterface)(nil).FeeAccumulatorAddress))
}

// GetAccountBalanceSOL mocks base method.
func (m *MockInterface) GetAccountBalanceSOL(arg0 context.Context, arg1 string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountBalanceSOL", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountBalanceSOL indicates an expected call of GetAccountBalanceSOL.
func (mr *MockInterfaceMockRecorder) GetAccountBalanceSOL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountBalanceSOL", reflect.TypeOf((*MockInterface)(nil).GetAccountBalanceSOL), arg0, arg1)
}

// GetConfirmedTransaction mocks base method.
func (m *MockInterface) GetConfirmedTransaction(arg0 context.Context, arg1 string) (GetConfirmedTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfirmedTransaction", arg0, arg1)
	ret0, _ := ret[0].(GetConfirmedTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfirmedTransaction indicates an expected call of GetConfirmedTransaction.
func (mr *MockInterfaceMockRecorder) GetConfirmedTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfirmedTransaction", reflect.TypeOf((*MockInterface)(nil).GetConfirmedTransaction), arg0, arg1)
}

// GetConfirmedTransactionForAccount mocks base method.
func (m *MockInterface) GetConfirmedTransactionForAccount(arg0 context.Context, arg1, arg2, arg3 string) (ConfirmedTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfirmedTransactionForAccount", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(ConfirmedTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfirmedTransactionForAccount indicates an expected call of GetConfirmedTransactionForAccount.
func (mr *MockInterfaceMockRecorder) GetConfirmedTransactionForAccount(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfirmedTransactionForAccount", reflect.TypeOf((*MockInterface)(nil).GetConfirmedTransactionForAccount), arg0, arg1, arg2, arg3)
}

// GetNFTMetadata mocks base method.
func (m *MockInterface) GetNFTMetadata(arg0 string) (*ArweaveNFTMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNFTMetadata", arg0)
	ret0, _ := ret[0].(*ArweaveNFTMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNFTMetadata indicates an expected call of GetNFTMetadata.
func (mr *MockInterfaceMockRecorder) GetNFTMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNFTMetadata", reflect.TypeOf((*MockInterface)(nil).GetNFTMetadata), arg0)
}

// GetNFTMintAddrs mocks base method.
func (m *MockInterface) GetNFTMintAddrs(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNFTMintAddrs", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNFTMintAddrs indicates an expected call of GetNFTMintAddrs.
func (mr *MockInterfaceMockRecorder) GetNFTMintAddrs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNFTMintAddrs", reflect.TypeOf((*MockInterface)(nil).GetNFTMintAddrs), arg0, arg1)
}

// GetNFTsByWalletAddress mocks base method.
func (m *MockInterface) GetNFTsByWalletAddress(arg0 context.Context, arg1 string) ([]*ArweaveNFTMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNFTsByWalletAddress", arg0, arg1)
	ret0, _ := ret[0].([]*ArweaveNFTMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNFTsByWalletAddress indicates an expected call of GetNFTsByWalletAddress.
func (mr *MockInterfaceMockRecorder) GetNFTsByWalletAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNFTsByWalletAddress", reflect.TypeOf((*MockInterface)(nil).GetNFTsByWalletAddress), arg0, arg1)
}

// GetTokenAccountBalance mocks base method.
func (m *MockInterface) GetTokenAccountBalance(arg0 context.Context, arg1 string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenAccountBalance", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenAccountBalance indicates an expected call of GetTokenAccountBalance.
func (mr *MockInterfaceMockRecorder) GetTokenAccountBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenAccountBalance", reflect.TypeOf((*MockInterface)(nil).GetTokenAccountBalance), arg0, arg1)
}

// GetTokenAccountBalanceWithAutoDerive mocks base method.
func (m *MockInterface) GetTokenAccountBalanceWithAutoDerive(arg0 context.Context, arg1, arg2 string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenAccountBalanceWithAutoDerive", arg0, arg1, arg2)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenAccountBalanceWithAutoDerive indicates an expected call of GetTokenAccountBalanceWithAutoDerive.
func (mr *MockInterfaceMockRecorder) GetTokenAccountBalanceWithAutoDerive(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenAccountBalanceWithAutoDerive", reflect.TypeOf((*MockInterface)(nil).GetTokenAccountBalanceWithAutoDerive), arg0, arg1, arg2)
}

// GetTransactions mocks base method.
func (m *MockInterface) GetTransactions(arg0 context.Context, arg1, arg2, arg3 string) ([]ConfirmedTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]ConfirmedTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockInterfaceMockRecorder) GetTransactions(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockInterface)(nil).GetTransactions), arg0, arg1, arg2, arg3)
}

// GetTransactionsWithAutoDerive mocks base method.
func (m *MockInterface) GetTransactionsWithAutoDerive(arg0 context.Context, arg1, arg2 string) ([]ConfirmedTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsWithAutoDerive", arg0, arg1, arg2)
	ret0, _ := ret[0].([]ConfirmedTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsWithAutoDerive indicates an expected call of GetTransactionsWithAutoDerive.
func (mr *MockInterfaceMockRecorder) GetTransactionsWithAutoDerive(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsWithAutoDerive", reflect.TypeOf((*MockInterface)(nil).GetTransactionsWithAutoDerive), arg0, arg1, arg2)
}

// GiveAssetsWithAutoDerive mocks base method.
func (m *MockInterface) GiveAssetsWithAutoDerive(arg0 context.Context, arg1 string, arg2, arg3 types.Account, arg4 string, arg5 float64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GiveAssetsWithAutoDerive", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GiveAssetsWithAutoDerive indicates an expected call of GiveAssetsWithAutoDerive.
func (mr *MockInterfaceMockRecorder) GiveAssetsWithAutoDerive(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GiveAssetsWithAutoDerive", reflect.TypeOf((*MockInterface)(nil).GiveAssetsWithAutoDerive), arg0, arg1, arg2, arg3, arg4, arg5)
}

// InitAccountToUseAsset mocks base method.
func (m *MockInterface) InitAccountToUseAsset(arg0 context.Context, arg1, arg2, arg3, arg4 types.Account) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitAccountToUseAsset", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InitAccountToUseAsset indicates an expected call of InitAccountToUseAsset.
func (mr *MockInterfaceMockRecorder) InitAccountToUseAsset(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitAccountToUseAsset", reflect.TypeOf((*MockInterface)(nil).InitAccountToUseAsset), arg0, arg1, arg2, arg3, arg4)
}

// InitializeStakePool mocks base method.
func (m *MockInterface) InitializeStakePool(arg0 context.Context, arg1, arg2 types.Account, arg3 common.PublicKey) (string, types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeStakePool", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(types.Account)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InitializeStakePool indicates an expected call of InitializeStakePool.
func (mr *MockInterfaceMockRecorder) InitializeStakePool(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeStakePool", reflect.TypeOf((*MockInterface)(nil).InitializeStakePool), arg0, arg1, arg2, arg3)
}

// IssueAsset mocks base method.
func (m *MockInterface) IssueAsset(arg0 context.Context, arg1, arg2, arg3 types.Account, arg4 common.PublicKey, arg5 float64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueAsset", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IssueAsset indicates an expected call of IssueAsset.
func (mr *MockInterfaceMockRecorder) IssueAsset(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueAsset", reflect.TypeOf((*MockInterface)(nil).IssueAsset), arg0, arg1, arg2, arg3, arg4, arg5)
}

// NewAccount mocks base method.
func (m *MockInterface) NewAccount() types.Account {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAccount")
	ret0, _ := ret[0].(types.Account)
	return ret0
}

// NewAccount indicates an expected call of NewAccount.
func (mr *MockInterfaceMockRecorder) NewAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAccount", reflect.TypeOf((*MockInterface)(nil).NewAccount))
}

// PrepareSendAssetsTx mocks base method.
func (m *MockInterface) PrepareSendAssetsTx(arg0 context.Context, arg1 string, arg2, arg3 types.Account, arg4 string, arg5 float64, arg6 *SendAssetsConfig) (*PrepareTxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareSendAssetsTx", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(*PrepareTxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareSendAssetsTx indicates an expected call of PrepareSendAssetsTx.
func (mr *MockInterfaceMockRecorder) PrepareSendAssetsTx(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareSendAssetsTx", reflect.TypeOf((*MockInterface)(nil).PrepareSendAssetsTx), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// PublicKeyFromString mocks base method.
func (m *MockInterface) PublicKeyFromString(arg0 string) common.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKeyFromString", arg0)
	ret0, _ := ret[0].(common.PublicKey)
	return ret0
}

// PublicKeyFromString indicates an expected call of PublicKeyFromString.
func (mr *MockInterfaceMockRecorder) PublicKeyFromString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKeyFromString", reflect.TypeOf((*MockInterface)(nil).PublicKeyFromString), arg0)
}

// RequestAirdrop mocks base method.
func (m *MockInterface) RequestAirdrop(arg0 context.Context, arg1 string, arg2 float64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestAirdrop", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestAirdrop indicates an expected call of RequestAirdrop.
func (mr *MockInterfaceMockRecorder) RequestAirdrop(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestAirdrop", reflect.TypeOf((*MockInterface)(nil).RequestAirdrop), arg0, arg1, arg2)
}

// SendAssetsWithAutoDerive mocks base method.
func (m *MockInterface) SendAssetsWithAutoDerive(arg0 context.Context, arg1 string, arg2, arg3 types.Account, arg4 string, arg5 float64, arg6 *SendAssetsConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAssetsWithAutoDerive", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendAssetsWithAutoDerive indicates an expected call of SendAssetsWithAutoDerive.
func (mr *MockInterfaceMockRecorder) SendAssetsWithAutoDerive(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAssetsWithAutoDerive", reflect.TypeOf((*MockInterface)(nil).SendAssetsWithAutoDerive), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// SendTransaction mocks base method.
func (m *MockInterface) SendTransaction(arg0 context.Context, arg1, arg2 types.Account, arg3 ...types.Instruction) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendTransaction", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTransaction indicates an expected call of SendTransaction.
func (mr *MockInterfaceMockRecorder) SendTransaction(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockInterface)(nil).SendTransaction), varargs...)
}

// SerializeTxMessage mocks base method.
func (m *MockInterface) SerializeTxMessage(arg0 types.Message) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SerializeTxMessage", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SerializeTxMessage indicates an expected call of SerializeTxMessage.
func (mr *MockInterfaceMockRecorder) SerializeTxMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SerializeTxMessage", reflect.TypeOf((*MockInterface)(nil).SerializeTxMessage), arg0)
}

// Stake mocks base method.
func (m *MockInterface) Stake(arg0 context.Context, arg1, arg2 types.Account, arg3, arg4 common.PublicKey, arg5 int64, arg6 uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stake", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stake indicates an expected call of Stake.
func (mr *MockInterfaceMockRecorder) Stake(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stake", reflect.TypeOf((*MockInterface)(nil).Stake), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// TransactionDeserialize mocks base method.
func (m *MockInterface) TransactionDeserialize(arg0 []byte) (types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionDeserialize", arg0)
	ret0, _ := ret[0].(types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionDeserialize indicates an expected call of TransactionDeserialize.
func (mr *MockInterfaceMockRecorder) TransactionDeserialize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionDeserialize", reflect.TypeOf((*MockInterface)(nil).TransactionDeserialize), arg0)
}

// Unstake mocks base method.
func (m *MockInterface) Unstake(arg0 context.Context, arg1, arg2 types.Account, arg3, arg4 common.PublicKey) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unstake", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unstake indicates an expected call of Unstake.
func (mr *MockInterfaceMockRecorder) Unstake(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unstake", reflect.TypeOf((*MockInterface)(nil).Unstake), arg0, arg1, arg2, arg3, arg4)
}
