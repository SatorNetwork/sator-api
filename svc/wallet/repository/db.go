// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addEthereumAccountStmt, err = db.PrepareContext(ctx, addEthereumAccount); err != nil {
		return nil, fmt.Errorf("error preparing query AddEthereumAccount: %w", err)
	}
	if q.addSolanaAccountStmt, err = db.PrepareContext(ctx, addSolanaAccount); err != nil {
		return nil, fmt.Errorf("error preparing query AddSolanaAccount: %w", err)
	}
	if q.addStakeStmt, err = db.PrepareContext(ctx, addStake); err != nil {
		return nil, fmt.Errorf("error preparing query AddStake: %w", err)
	}
	if q.addStakeLevelStmt, err = db.PrepareContext(ctx, addStakeLevel); err != nil {
		return nil, fmt.Errorf("error preparing query AddStakeLevel: %w", err)
	}
	if q.addTokenTransferStmt, err = db.PrepareContext(ctx, addTokenTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query AddTokenTransfer: %w", err)
	}
	if q.checkRecipientAddressStmt, err = db.PrepareContext(ctx, checkRecipientAddress); err != nil {
		return nil, fmt.Errorf("error preparing query CheckRecipientAddress: %w", err)
	}
	if q.createWalletStmt, err = db.PrepareContext(ctx, createWallet); err != nil {
		return nil, fmt.Errorf("error preparing query CreateWallet: %w", err)
	}
	if q.deleteStakeByUserIDStmt, err = db.PrepareContext(ctx, deleteStakeByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteStakeByUserID: %w", err)
	}
	if q.deleteWalletByIDStmt, err = db.PrepareContext(ctx, deleteWalletByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteWalletByID: %w", err)
	}
	if q.doesUserHaveFraudulentTransfersStmt, err = db.PrepareContext(ctx, doesUserHaveFraudulentTransfers); err != nil {
		return nil, fmt.Errorf("error preparing query DoesUserHaveFraudulentTransfers: %w", err)
	}
	if q.doesUserMakeTransferForLastMinuteStmt, err = db.PrepareContext(ctx, doesUserMakeTransferForLastMinute); err != nil {
		return nil, fmt.Errorf("error preparing query DoesUserMakeTransferForLastMinute: %w", err)
	}
	if q.getAllEnabledStakeLevelsStmt, err = db.PrepareContext(ctx, getAllEnabledStakeLevels); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllEnabledStakeLevels: %w", err)
	}
	if q.getAllStakeLevelsStmt, err = db.PrepareContext(ctx, getAllStakeLevels); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllStakeLevels: %w", err)
	}
	if q.getAllStakesStmt, err = db.PrepareContext(ctx, getAllStakes); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllStakes: %w", err)
	}
	if q.getEthereumAccountByIDStmt, err = db.PrepareContext(ctx, getEthereumAccountByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEthereumAccountByID: %w", err)
	}
	if q.getEthereumAccountByUserIDAndTypeStmt, err = db.PrepareContext(ctx, getEthereumAccountByUserIDAndType); err != nil {
		return nil, fmt.Errorf("error preparing query GetEthereumAccountByUserIDAndType: %w", err)
	}
	if q.getMinimalStakeLevelStmt, err = db.PrepareContext(ctx, getMinimalStakeLevel); err != nil {
		return nil, fmt.Errorf("error preparing query GetMinimalStakeLevel: %w", err)
	}
	if q.getSolanaAccountByIDStmt, err = db.PrepareContext(ctx, getSolanaAccountByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetSolanaAccountByID: %w", err)
	}
	if q.getSolanaAccountByTypeStmt, err = db.PrepareContext(ctx, getSolanaAccountByType); err != nil {
		return nil, fmt.Errorf("error preparing query GetSolanaAccountByType: %w", err)
	}
	if q.getSolanaAccountByUserIDAndTypeStmt, err = db.PrepareContext(ctx, getSolanaAccountByUserIDAndType); err != nil {
		return nil, fmt.Errorf("error preparing query GetSolanaAccountByUserIDAndType: %w", err)
	}
	if q.getSolanaAccountTypeByPublicKeyStmt, err = db.PrepareContext(ctx, getSolanaAccountTypeByPublicKey); err != nil {
		return nil, fmt.Errorf("error preparing query GetSolanaAccountTypeByPublicKey: %w", err)
	}
	if q.getStakeByUserIDStmt, err = db.PrepareContext(ctx, getStakeByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetStakeByUserID: %w", err)
	}
	if q.getStakeLevelByAmountStmt, err = db.PrepareContext(ctx, getStakeLevelByAmount); err != nil {
		return nil, fmt.Errorf("error preparing query GetStakeLevelByAmount: %w", err)
	}
	if q.getStakeLevelByIDStmt, err = db.PrepareContext(ctx, getStakeLevelByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetStakeLevelByID: %w", err)
	}
	if q.getTotalStakeStmt, err = db.PrepareContext(ctx, getTotalStake); err != nil {
		return nil, fmt.Errorf("error preparing query GetTotalStake: %w", err)
	}
	if q.getWalletByEthereumAccountIDStmt, err = db.PrepareContext(ctx, getWalletByEthereumAccountID); err != nil {
		return nil, fmt.Errorf("error preparing query GetWalletByEthereumAccountID: %w", err)
	}
	if q.getWalletByIDStmt, err = db.PrepareContext(ctx, getWalletByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetWalletByID: %w", err)
	}
	if q.getWalletBySolanaAccountIDStmt, err = db.PrepareContext(ctx, getWalletBySolanaAccountID); err != nil {
		return nil, fmt.Errorf("error preparing query GetWalletBySolanaAccountID: %w", err)
	}
	if q.getWalletByUserIDAndTypeStmt, err = db.PrepareContext(ctx, getWalletByUserIDAndType); err != nil {
		return nil, fmt.Errorf("error preparing query GetWalletByUserIDAndType: %w", err)
	}
	if q.getWalletsByUserIDStmt, err = db.PrepareContext(ctx, getWalletsByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetWalletsByUserID: %w", err)
	}
	if q.updateStakeStmt, err = db.PrepareContext(ctx, updateStake); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateStake: %w", err)
	}
	if q.updateStakeLevelStmt, err = db.PrepareContext(ctx, updateStakeLevel); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateStakeLevel: %w", err)
	}
	if q.updateTokenTransferStmt, err = db.PrepareContext(ctx, updateTokenTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTokenTransfer: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addEthereumAccountStmt != nil {
		if cerr := q.addEthereumAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addEthereumAccountStmt: %w", cerr)
		}
	}
	if q.addSolanaAccountStmt != nil {
		if cerr := q.addSolanaAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addSolanaAccountStmt: %w", cerr)
		}
	}
	if q.addStakeStmt != nil {
		if cerr := q.addStakeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addStakeStmt: %w", cerr)
		}
	}
	if q.addStakeLevelStmt != nil {
		if cerr := q.addStakeLevelStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addStakeLevelStmt: %w", cerr)
		}
	}
	if q.addTokenTransferStmt != nil {
		if cerr := q.addTokenTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addTokenTransferStmt: %w", cerr)
		}
	}
	if q.checkRecipientAddressStmt != nil {
		if cerr := q.checkRecipientAddressStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkRecipientAddressStmt: %w", cerr)
		}
	}
	if q.createWalletStmt != nil {
		if cerr := q.createWalletStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createWalletStmt: %w", cerr)
		}
	}
	if q.deleteStakeByUserIDStmt != nil {
		if cerr := q.deleteStakeByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteStakeByUserIDStmt: %w", cerr)
		}
	}
	if q.deleteWalletByIDStmt != nil {
		if cerr := q.deleteWalletByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteWalletByIDStmt: %w", cerr)
		}
	}
	if q.doesUserHaveFraudulentTransfersStmt != nil {
		if cerr := q.doesUserHaveFraudulentTransfersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing doesUserHaveFraudulentTransfersStmt: %w", cerr)
		}
	}
	if q.doesUserMakeTransferForLastMinuteStmt != nil {
		if cerr := q.doesUserMakeTransferForLastMinuteStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing doesUserMakeTransferForLastMinuteStmt: %w", cerr)
		}
	}
	if q.getAllEnabledStakeLevelsStmt != nil {
		if cerr := q.getAllEnabledStakeLevelsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllEnabledStakeLevelsStmt: %w", cerr)
		}
	}
	if q.getAllStakeLevelsStmt != nil {
		if cerr := q.getAllStakeLevelsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllStakeLevelsStmt: %w", cerr)
		}
	}
	if q.getAllStakesStmt != nil {
		if cerr := q.getAllStakesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllStakesStmt: %w", cerr)
		}
	}
	if q.getEthereumAccountByIDStmt != nil {
		if cerr := q.getEthereumAccountByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEthereumAccountByIDStmt: %w", cerr)
		}
	}
	if q.getEthereumAccountByUserIDAndTypeStmt != nil {
		if cerr := q.getEthereumAccountByUserIDAndTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEthereumAccountByUserIDAndTypeStmt: %w", cerr)
		}
	}
	if q.getMinimalStakeLevelStmt != nil {
		if cerr := q.getMinimalStakeLevelStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMinimalStakeLevelStmt: %w", cerr)
		}
	}
	if q.getSolanaAccountByIDStmt != nil {
		if cerr := q.getSolanaAccountByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSolanaAccountByIDStmt: %w", cerr)
		}
	}
	if q.getSolanaAccountByTypeStmt != nil {
		if cerr := q.getSolanaAccountByTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSolanaAccountByTypeStmt: %w", cerr)
		}
	}
	if q.getSolanaAccountByUserIDAndTypeStmt != nil {
		if cerr := q.getSolanaAccountByUserIDAndTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSolanaAccountByUserIDAndTypeStmt: %w", cerr)
		}
	}
	if q.getSolanaAccountTypeByPublicKeyStmt != nil {
		if cerr := q.getSolanaAccountTypeByPublicKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSolanaAccountTypeByPublicKeyStmt: %w", cerr)
		}
	}
	if q.getStakeByUserIDStmt != nil {
		if cerr := q.getStakeByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStakeByUserIDStmt: %w", cerr)
		}
	}
	if q.getStakeLevelByAmountStmt != nil {
		if cerr := q.getStakeLevelByAmountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStakeLevelByAmountStmt: %w", cerr)
		}
	}
	if q.getStakeLevelByIDStmt != nil {
		if cerr := q.getStakeLevelByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStakeLevelByIDStmt: %w", cerr)
		}
	}
	if q.getTotalStakeStmt != nil {
		if cerr := q.getTotalStakeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTotalStakeStmt: %w", cerr)
		}
	}
	if q.getWalletByEthereumAccountIDStmt != nil {
		if cerr := q.getWalletByEthereumAccountIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWalletByEthereumAccountIDStmt: %w", cerr)
		}
	}
	if q.getWalletByIDStmt != nil {
		if cerr := q.getWalletByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWalletByIDStmt: %w", cerr)
		}
	}
	if q.getWalletBySolanaAccountIDStmt != nil {
		if cerr := q.getWalletBySolanaAccountIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWalletBySolanaAccountIDStmt: %w", cerr)
		}
	}
	if q.getWalletByUserIDAndTypeStmt != nil {
		if cerr := q.getWalletByUserIDAndTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWalletByUserIDAndTypeStmt: %w", cerr)
		}
	}
	if q.getWalletsByUserIDStmt != nil {
		if cerr := q.getWalletsByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWalletsByUserIDStmt: %w", cerr)
		}
	}
	if q.updateStakeStmt != nil {
		if cerr := q.updateStakeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateStakeStmt: %w", cerr)
		}
	}
	if q.updateStakeLevelStmt != nil {
		if cerr := q.updateStakeLevelStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateStakeLevelStmt: %w", cerr)
		}
	}
	if q.updateTokenTransferStmt != nil {
		if cerr := q.updateTokenTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTokenTransferStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                    DBTX
	tx                                    *sql.Tx
	addEthereumAccountStmt                *sql.Stmt
	addSolanaAccountStmt                  *sql.Stmt
	addStakeStmt                          *sql.Stmt
	addStakeLevelStmt                     *sql.Stmt
	addTokenTransferStmt                  *sql.Stmt
	checkRecipientAddressStmt             *sql.Stmt
	createWalletStmt                      *sql.Stmt
	deleteStakeByUserIDStmt               *sql.Stmt
	deleteWalletByIDStmt                  *sql.Stmt
	doesUserHaveFraudulentTransfersStmt   *sql.Stmt
	doesUserMakeTransferForLastMinuteStmt *sql.Stmt
	getAllEnabledStakeLevelsStmt          *sql.Stmt
	getAllStakeLevelsStmt                 *sql.Stmt
	getAllStakesStmt                      *sql.Stmt
	getEthereumAccountByIDStmt            *sql.Stmt
	getEthereumAccountByUserIDAndTypeStmt *sql.Stmt
	getMinimalStakeLevelStmt              *sql.Stmt
	getSolanaAccountByIDStmt              *sql.Stmt
	getSolanaAccountByTypeStmt            *sql.Stmt
	getSolanaAccountByUserIDAndTypeStmt   *sql.Stmt
	getSolanaAccountTypeByPublicKeyStmt   *sql.Stmt
	getStakeByUserIDStmt                  *sql.Stmt
	getStakeLevelByAmountStmt             *sql.Stmt
	getStakeLevelByIDStmt                 *sql.Stmt
	getTotalStakeStmt                     *sql.Stmt
	getWalletByEthereumAccountIDStmt      *sql.Stmt
	getWalletByIDStmt                     *sql.Stmt
	getWalletBySolanaAccountIDStmt        *sql.Stmt
	getWalletByUserIDAndTypeStmt          *sql.Stmt
	getWalletsByUserIDStmt                *sql.Stmt
	updateStakeStmt                       *sql.Stmt
	updateStakeLevelStmt                  *sql.Stmt
	updateTokenTransferStmt               *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                    tx,
		tx:                                    tx,
		addEthereumAccountStmt:                q.addEthereumAccountStmt,
		addSolanaAccountStmt:                  q.addSolanaAccountStmt,
		addStakeStmt:                          q.addStakeStmt,
		addStakeLevelStmt:                     q.addStakeLevelStmt,
		addTokenTransferStmt:                  q.addTokenTransferStmt,
		checkRecipientAddressStmt:             q.checkRecipientAddressStmt,
		createWalletStmt:                      q.createWalletStmt,
		deleteStakeByUserIDStmt:               q.deleteStakeByUserIDStmt,
		deleteWalletByIDStmt:                  q.deleteWalletByIDStmt,
		doesUserHaveFraudulentTransfersStmt:   q.doesUserHaveFraudulentTransfersStmt,
		doesUserMakeTransferForLastMinuteStmt: q.doesUserMakeTransferForLastMinuteStmt,
		getAllEnabledStakeLevelsStmt:          q.getAllEnabledStakeLevelsStmt,
		getAllStakeLevelsStmt:                 q.getAllStakeLevelsStmt,
		getAllStakesStmt:                      q.getAllStakesStmt,
		getEthereumAccountByIDStmt:            q.getEthereumAccountByIDStmt,
		getEthereumAccountByUserIDAndTypeStmt: q.getEthereumAccountByUserIDAndTypeStmt,
		getMinimalStakeLevelStmt:              q.getMinimalStakeLevelStmt,
		getSolanaAccountByIDStmt:              q.getSolanaAccountByIDStmt,
		getSolanaAccountByTypeStmt:            q.getSolanaAccountByTypeStmt,
		getSolanaAccountByUserIDAndTypeStmt:   q.getSolanaAccountByUserIDAndTypeStmt,
		getSolanaAccountTypeByPublicKeyStmt:   q.getSolanaAccountTypeByPublicKeyStmt,
		getStakeByUserIDStmt:                  q.getStakeByUserIDStmt,
		getStakeLevelByAmountStmt:             q.getStakeLevelByAmountStmt,
		getStakeLevelByIDStmt:                 q.getStakeLevelByIDStmt,
		getTotalStakeStmt:                     q.getTotalStakeStmt,
		getWalletByEthereumAccountIDStmt:      q.getWalletByEthereumAccountIDStmt,
		getWalletByIDStmt:                     q.getWalletByIDStmt,
		getWalletBySolanaAccountIDStmt:        q.getWalletBySolanaAccountIDStmt,
		getWalletByUserIDAndTypeStmt:          q.getWalletByUserIDAndTypeStmt,
		getWalletsByUserIDStmt:                q.getWalletsByUserIDStmt,
		updateStakeStmt:                       q.updateStakeStmt,
		updateStakeLevelStmt:                  q.updateStakeLevelStmt,
		updateTokenTransferStmt:               q.updateTokenTransferStmt,
	}
}
