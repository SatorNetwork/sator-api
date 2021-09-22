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
	feePayerTestnet = "74gL9GZyyHZQtAfGmnBckeAkhxdeQvngzaL2qAoHSQLg"
)

func main() {
	c := solana.New(solanaApiBaseUrl)

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