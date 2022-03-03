package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/dmitrymomot/go-env"

	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
)

var (
	// Solana
	solanaApiBaseUrl     = env.MustString("SOLANA_API_BASE_URL")
	solanaSystemProgram  = env.MustString("SOLANA_SYSTEM_PROGRAM")
	solanaSysvarRent     = env.MustString("SOLANA_SYSVAR_RENT")
	solanaSysvarClock    = env.MustString("SOLANA_SYSVAR_CLOCK")
	solanaSplToken       = env.MustString("SOLANA_SPL_TOKEN")
	solanaStakeProgramID = env.MustString("SOLANA_STAKE_PROGRAM_ID")
	feePayerTestnet      = "74gL9GZyyHZQtAfGmnBckeAkhxdeQvngzaL2qAoHSQLg"
)

func main() {
	c := solana_client.New(solanaApiBaseUrl, solana_client.Config{
		SystemProgram:  solanaSystemProgram,
		SysvarRent:     solanaSysvarRent,
		SysvarClock:    solanaSysvarClock,
		SplToken:       solanaSplToken,
		StakeProgramID: solanaStakeProgramID,
	})

	for i := 0; i < 66; i++ {
		tx, err := c.RequestAirdrop(context.TODO(), feePayerTestnet, 1)
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else {
			log.Printf("airdrop transaction: %s", tx)
		}

		time.Sleep(time.Second * time.Duration(rand.Int63n(120)+30))
	}

}
