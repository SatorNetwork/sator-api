package wallet

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

type DB struct {
	walletRepository *walletRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	walletRepository, err := walletRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare wallet repository")
	}

	return &DB{
		walletRepository: walletRepository,
	}, nil
}

func (db *DB) AddSolanaAccount(ctx context.Context, accountType, publicKey string, privateKey []byte) (walletRepo.SolanaAccount, error) {
	 return db.walletRepository.AddSolanaAccount(ctx, walletRepo.AddSolanaAccountParams{
		AccountType: accountType,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
	})
}

// 12
// Бердибек Садугаш Галымкызы
// Томирис 87774179033
