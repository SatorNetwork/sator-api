package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/dmitrymomot/go-env"
)

var (
	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")
	feePayerTestnet  = "74gL9GZyyHZQtAfGmnBckeAkhxdeQvngzaL2qAoHSQLg"
	systemProgram    = env.MustString("SOLANA_SYSTEM_PROGRAM")
	sysvarRent       = env.MustString("SOLANA_SYSVAR_RENT")
	sysvarClock      = env.MustString("SOLANA_SYSVAR_CLOCK")
	splToken         = env.MustString("SOLANA_SPL_TOKEN")
	stakeProgramID   = env.MustString("SOLANA_STAKE_PROGRAM_ID")
	rewardProgramID  = env.MustString("SOLANA_REWARD_PROGRAM_ID")
)

func main() {
	c := solana.New(solanaApiBaseUrl, solana.Config{
		SystemProgram:   systemProgram,
		SysvarRent:      sysvarRent,
		SysvarClock:     sysvarClock,
		SplToken:        splToken,
		StakeProgramID:  stakeProgramID,
		RewardProgramID: rewardProgramID,
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
