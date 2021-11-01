package ethereum_test

import (
	"testing"

	"github.com/SatorNetwork/sator-api/internal/ethereum"
)

func TestNew(t *testing.T) {
	err := ethereum.MintNFT()
	print(err)
}
