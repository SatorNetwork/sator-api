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
	if q.addNFTStmt, err = db.PrepareContext(ctx, addNFT); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFT: %w", err)
	}
	if q.addNFTPackStmt, err = db.PrepareContext(ctx, addNFTPack); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFTPack: %w", err)
	}
	if q.addNewPlayerStmt, err = db.PrepareContext(ctx, addNewPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query AddNewPlayer: %w", err)
	}
	if q.addSettingStmt, err = db.PrepareContext(ctx, addSetting); err != nil {
		return nil, fmt.Errorf("error preparing query AddSetting: %w", err)
	}
	if q.craftNFTsStmt, err = db.PrepareContext(ctx, craftNFTs); err != nil {
		return nil, fmt.Errorf("error preparing query CraftNFTs: %w", err)
	}
	if q.deleteNFTStmt, err = db.PrepareContext(ctx, deleteNFT); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFT: %w", err)
	}
	if q.deleteNFTPackStmt, err = db.PrepareContext(ctx, deleteNFTPack); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFTPack: %w", err)
	}
	if q.deleteSettingStmt, err = db.PrepareContext(ctx, deleteSetting); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSetting: %w", err)
	}
	if q.finishGameStmt, err = db.PrepareContext(ctx, finishGame); err != nil {
		return nil, fmt.Errorf("error preparing query FinishGame: %w", err)
	}
	if q.getCurrentGameStmt, err = db.PrepareContext(ctx, getCurrentGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetCurrentGame: %w", err)
	}
	if q.getNFTStmt, err = db.PrepareContext(ctx, getNFT); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFT: %w", err)
	}
	if q.getNFTPackStmt, err = db.PrepareContext(ctx, getNFTPack); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTPack: %w", err)
	}
	if q.getNFTPacksListStmt, err = db.PrepareContext(ctx, getNFTPacksList); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTPacksList: %w", err)
	}
	if q.getPlayerStmt, err = db.PrepareContext(ctx, getPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayer: %w", err)
	}
	if q.getSettingByKeyStmt, err = db.PrepareContext(ctx, getSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettingByKey: %w", err)
	}
	if q.getSettingsStmt, err = db.PrepareContext(ctx, getSettings); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettings: %w", err)
	}
	if q.getUserNFTStmt, err = db.PrepareContext(ctx, getUserNFT); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserNFT: %w", err)
	}
	if q.getUserNFTByIDsStmt, err = db.PrepareContext(ctx, getUserNFTByIDs); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserNFTByIDs: %w", err)
	}
	if q.getUserNFTsStmt, err = db.PrepareContext(ctx, getUserNFTs); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserNFTs: %w", err)
	}
	if q.getUserRewardsStmt, err = db.PrepareContext(ctx, getUserRewards); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRewards: %w", err)
	}
	if q.getUserRewardsDepositedStmt, err = db.PrepareContext(ctx, getUserRewardsDeposited); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRewardsDeposited: %w", err)
	}
	if q.getUserRewardsWithdrawnStmt, err = db.PrepareContext(ctx, getUserRewardsWithdrawn); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRewardsWithdrawn: %w", err)
	}
	if q.refillEnergyOfPlayerStmt, err = db.PrepareContext(ctx, refillEnergyOfPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query RefillEnergyOfPlayer: %w", err)
	}
	if q.rewardsDepositStmt, err = db.PrepareContext(ctx, rewardsDeposit); err != nil {
		return nil, fmt.Errorf("error preparing query RewardsDeposit: %w", err)
	}
	if q.rewardsWithdrawStmt, err = db.PrepareContext(ctx, rewardsWithdraw); err != nil {
		return nil, fmt.Errorf("error preparing query RewardsWithdraw: %w", err)
	}
	if q.softDeleteNFTPackStmt, err = db.PrepareContext(ctx, softDeleteNFTPack); err != nil {
		return nil, fmt.Errorf("error preparing query SoftDeleteNFTPack: %w", err)
	}
	if q.startGameStmt, err = db.PrepareContext(ctx, startGame); err != nil {
		return nil, fmt.Errorf("error preparing query StartGame: %w", err)
	}
	if q.storeSelectedNFTStmt, err = db.PrepareContext(ctx, storeSelectedNFT); err != nil {
		return nil, fmt.Errorf("error preparing query StoreSelectedNFT: %w", err)
	}
	if q.updateNFTPackStmt, err = db.PrepareContext(ctx, updateNFTPack); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFTPack: %w", err)
	}
	if q.updateSettingStmt, err = db.PrepareContext(ctx, updateSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSetting: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addNFTStmt != nil {
		if cerr := q.addNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTStmt: %w", cerr)
		}
	}
	if q.addNFTPackStmt != nil {
		if cerr := q.addNFTPackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTPackStmt: %w", cerr)
		}
	}
	if q.addNewPlayerStmt != nil {
		if cerr := q.addNewPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNewPlayerStmt: %w", cerr)
		}
	}
	if q.addSettingStmt != nil {
		if cerr := q.addSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addSettingStmt: %w", cerr)
		}
	}
	if q.craftNFTsStmt != nil {
		if cerr := q.craftNFTsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing craftNFTsStmt: %w", cerr)
		}
	}
	if q.deleteNFTStmt != nil {
		if cerr := q.deleteNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFTStmt: %w", cerr)
		}
	}
	if q.deleteNFTPackStmt != nil {
		if cerr := q.deleteNFTPackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFTPackStmt: %w", cerr)
		}
	}
	if q.deleteSettingStmt != nil {
		if cerr := q.deleteSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSettingStmt: %w", cerr)
		}
	}
	if q.finishGameStmt != nil {
		if cerr := q.finishGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing finishGameStmt: %w", cerr)
		}
	}
	if q.getCurrentGameStmt != nil {
		if cerr := q.getCurrentGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCurrentGameStmt: %w", cerr)
		}
	}
	if q.getNFTStmt != nil {
		if cerr := q.getNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTStmt: %w", cerr)
		}
	}
	if q.getNFTPackStmt != nil {
		if cerr := q.getNFTPackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTPackStmt: %w", cerr)
		}
	}
	if q.getNFTPacksListStmt != nil {
		if cerr := q.getNFTPacksListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTPacksListStmt: %w", cerr)
		}
	}
	if q.getPlayerStmt != nil {
		if cerr := q.getPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStmt: %w", cerr)
		}
	}
	if q.getSettingByKeyStmt != nil {
		if cerr := q.getSettingByKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingByKeyStmt: %w", cerr)
		}
	}
	if q.getSettingsStmt != nil {
		if cerr := q.getSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingsStmt: %w", cerr)
		}
	}
	if q.getUserNFTStmt != nil {
		if cerr := q.getUserNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserNFTStmt: %w", cerr)
		}
	}
	if q.getUserNFTByIDsStmt != nil {
		if cerr := q.getUserNFTByIDsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserNFTByIDsStmt: %w", cerr)
		}
	}
	if q.getUserNFTsStmt != nil {
		if cerr := q.getUserNFTsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserNFTsStmt: %w", cerr)
		}
	}
	if q.getUserRewardsStmt != nil {
		if cerr := q.getUserRewardsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRewardsStmt: %w", cerr)
		}
	}
	if q.getUserRewardsDepositedStmt != nil {
		if cerr := q.getUserRewardsDepositedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRewardsDepositedStmt: %w", cerr)
		}
	}
	if q.getUserRewardsWithdrawnStmt != nil {
		if cerr := q.getUserRewardsWithdrawnStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRewardsWithdrawnStmt: %w", cerr)
		}
	}
	if q.refillEnergyOfPlayerStmt != nil {
		if cerr := q.refillEnergyOfPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing refillEnergyOfPlayerStmt: %w", cerr)
		}
	}
	if q.rewardsDepositStmt != nil {
		if cerr := q.rewardsDepositStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing rewardsDepositStmt: %w", cerr)
		}
	}
	if q.rewardsWithdrawStmt != nil {
		if cerr := q.rewardsWithdrawStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing rewardsWithdrawStmt: %w", cerr)
		}
	}
	if q.softDeleteNFTPackStmt != nil {
		if cerr := q.softDeleteNFTPackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing softDeleteNFTPackStmt: %w", cerr)
		}
	}
	if q.startGameStmt != nil {
		if cerr := q.startGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing startGameStmt: %w", cerr)
		}
	}
	if q.storeSelectedNFTStmt != nil {
		if cerr := q.storeSelectedNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing storeSelectedNFTStmt: %w", cerr)
		}
	}
	if q.updateNFTPackStmt != nil {
		if cerr := q.updateNFTPackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFTPackStmt: %w", cerr)
		}
	}
	if q.updateSettingStmt != nil {
		if cerr := q.updateSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSettingStmt: %w", cerr)
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
	db                          DBTX
	tx                          *sql.Tx
	addNFTStmt                  *sql.Stmt
	addNFTPackStmt              *sql.Stmt
	addNewPlayerStmt            *sql.Stmt
	addSettingStmt              *sql.Stmt
	craftNFTsStmt               *sql.Stmt
	deleteNFTStmt               *sql.Stmt
	deleteNFTPackStmt           *sql.Stmt
	deleteSettingStmt           *sql.Stmt
	finishGameStmt              *sql.Stmt
	getCurrentGameStmt          *sql.Stmt
	getNFTStmt                  *sql.Stmt
	getNFTPackStmt              *sql.Stmt
	getNFTPacksListStmt         *sql.Stmt
	getPlayerStmt               *sql.Stmt
	getSettingByKeyStmt         *sql.Stmt
	getSettingsStmt             *sql.Stmt
	getUserNFTStmt              *sql.Stmt
	getUserNFTByIDsStmt         *sql.Stmt
	getUserNFTsStmt             *sql.Stmt
	getUserRewardsStmt          *sql.Stmt
	getUserRewardsDepositedStmt *sql.Stmt
	getUserRewardsWithdrawnStmt *sql.Stmt
	refillEnergyOfPlayerStmt    *sql.Stmt
	rewardsDepositStmt          *sql.Stmt
	rewardsWithdrawStmt         *sql.Stmt
	softDeleteNFTPackStmt       *sql.Stmt
	startGameStmt               *sql.Stmt
	storeSelectedNFTStmt        *sql.Stmt
	updateNFTPackStmt           *sql.Stmt
	updateSettingStmt           *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                          tx,
		tx:                          tx,
		addNFTStmt:                  q.addNFTStmt,
		addNFTPackStmt:              q.addNFTPackStmt,
		addNewPlayerStmt:            q.addNewPlayerStmt,
		addSettingStmt:              q.addSettingStmt,
		craftNFTsStmt:               q.craftNFTsStmt,
		deleteNFTStmt:               q.deleteNFTStmt,
		deleteNFTPackStmt:           q.deleteNFTPackStmt,
		deleteSettingStmt:           q.deleteSettingStmt,
		finishGameStmt:              q.finishGameStmt,
		getCurrentGameStmt:          q.getCurrentGameStmt,
		getNFTStmt:                  q.getNFTStmt,
		getNFTPackStmt:              q.getNFTPackStmt,
		getNFTPacksListStmt:         q.getNFTPacksListStmt,
		getPlayerStmt:               q.getPlayerStmt,
		getSettingByKeyStmt:         q.getSettingByKeyStmt,
		getSettingsStmt:             q.getSettingsStmt,
		getUserNFTStmt:              q.getUserNFTStmt,
		getUserNFTByIDsStmt:         q.getUserNFTByIDsStmt,
		getUserNFTsStmt:             q.getUserNFTsStmt,
		getUserRewardsStmt:          q.getUserRewardsStmt,
		getUserRewardsDepositedStmt: q.getUserRewardsDepositedStmt,
		getUserRewardsWithdrawnStmt: q.getUserRewardsWithdrawnStmt,
		refillEnergyOfPlayerStmt:    q.refillEnergyOfPlayerStmt,
		rewardsDepositStmt:          q.rewardsDepositStmt,
		rewardsWithdrawStmt:         q.rewardsWithdrawStmt,
		softDeleteNFTPackStmt:       q.softDeleteNFTPackStmt,
		startGameStmt:               q.startGameStmt,
		storeSelectedNFTStmt:        q.storeSelectedNFTStmt,
		updateNFTPackStmt:           q.updateNFTPackStmt,
		updateSettingStmt:           q.updateSettingStmt,
	}
}
