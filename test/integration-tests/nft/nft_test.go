package iap

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"

	lib_nft_marketplace "github.com/SatorNetwork/sator-api/lib/nft_marketplace"
	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestBuyNFTViaMarketplace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	nftMarketplaceMock := lib_nft_marketplace.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.NftMarketplaceProvider, nftMarketplaceMock)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.ExpectTransactionDeserializeAny()
	solanaMock.ExpectSerializeTxMessageAny()

	nftMarketplaceMock.EXPECT().
		PrepareBuyTx(gomock.Any()).
		Return(&lib_nft_marketplace.PrepareBuyTxResponse{}, nil).
		Times(1)
	nftMarketplaceMock.EXPECT().
		SendPreparedBuyTx(gomock.Any()).
		Return(&lib_nft_marketplace.SendPreparedBuyTxResponse{}, nil).
		Times(1)

	defer app_config.RunAndWait()()

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	user := user.NewInitializedUser(signUpRequest, t)

	nftMintAddress := types.NewAccount().PublicKey.ToBase58()
	_, err := c.NftClient.BuyNFTViaMarketplace(user.AccessToken(), nftMintAddress)
	require.NoError(t, err)
}
