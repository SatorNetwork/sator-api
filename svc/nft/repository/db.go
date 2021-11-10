// Code generated by sqlc. DO NOT EDIT.

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
	if q.addNFTCategoryStmt, err = db.PrepareContext(ctx, addNFTCategory); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFTCategory: %w", err)
	}
	if q.addNFTItemStmt, err = db.PrepareContext(ctx, addNFTItem); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFTItem: %w", err)
	}
	if q.addNFTItemOwnerStmt, err = db.PrepareContext(ctx, addNFTItemOwner); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFTItemOwner: %w", err)
	}
	if q.addNFTRelationStmt, err = db.PrepareContext(ctx, addNFTRelation); err != nil {
		return nil, fmt.Errorf("error preparing query AddNFTRelation: %w", err)
	}
	if q.deleteNFTCategoryByIDStmt, err = db.PrepareContext(ctx, deleteNFTCategoryByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFTCategoryByID: %w", err)
	}
	if q.deleteNFTRelationStmt, err = db.PrepareContext(ctx, deleteNFTRelation); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFTRelation: %w", err)
	}
	if q.doesUserOwnNFTStmt, err = db.PrepareContext(ctx, doesUserOwnNFT); err != nil {
		return nil, fmt.Errorf("error preparing query DoesUserOwnNFT: %w", err)
	}
	if q.getMainNFTCategoryStmt, err = db.PrepareContext(ctx, getMainNFTCategory); err != nil {
		return nil, fmt.Errorf("error preparing query GetMainNFTCategory: %w", err)
	}
	if q.getNFTCategoriesListStmt, err = db.PrepareContext(ctx, getNFTCategoriesList); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTCategoriesList: %w", err)
	}
	if q.getNFTCategoryByIDStmt, err = db.PrepareContext(ctx, getNFTCategoryByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTCategoryByID: %w", err)
	}
	if q.getNFTItemByIDStmt, err = db.PrepareContext(ctx, getNFTItemByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTItemByID: %w", err)
	}
	if q.getNFTItemsListStmt, err = db.PrepareContext(ctx, getNFTItemsList); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTItemsList: %w", err)
	}
	if q.getNFTItemsListByOwnerIDStmt, err = db.PrepareContext(ctx, getNFTItemsListByOwnerID); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTItemsListByOwnerID: %w", err)
	}
	if q.getNFTItemsListByRelationIDStmt, err = db.PrepareContext(ctx, getNFTItemsListByRelationID); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFTItemsListByRelationID: %w", err)
	}
	if q.resetMainNFTCategoryStmt, err = db.PrepareContext(ctx, resetMainNFTCategory); err != nil {
		return nil, fmt.Errorf("error preparing query ResetMainNFTCategory: %w", err)
	}
	if q.updateNFTCategoryStmt, err = db.PrepareContext(ctx, updateNFTCategory); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFTCategory: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addNFTCategoryStmt != nil {
		if cerr := q.addNFTCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTCategoryStmt: %w", cerr)
		}
	}
	if q.addNFTItemStmt != nil {
		if cerr := q.addNFTItemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTItemStmt: %w", cerr)
		}
	}
	if q.addNFTItemOwnerStmt != nil {
		if cerr := q.addNFTItemOwnerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTItemOwnerStmt: %w", cerr)
		}
	}
	if q.addNFTRelationStmt != nil {
		if cerr := q.addNFTRelationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNFTRelationStmt: %w", cerr)
		}
	}
	if q.deleteNFTCategoryByIDStmt != nil {
		if cerr := q.deleteNFTCategoryByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFTCategoryByIDStmt: %w", cerr)
		}
	}
	if q.deleteNFTRelationStmt != nil {
		if cerr := q.deleteNFTRelationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFTRelationStmt: %w", cerr)
		}
	}
	if q.doesUserOwnNFTStmt != nil {
		if cerr := q.doesUserOwnNFTStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing doesUserOwnNFTStmt: %w", cerr)
		}
	}
	if q.getMainNFTCategoryStmt != nil {
		if cerr := q.getMainNFTCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMainNFTCategoryStmt: %w", cerr)
		}
	}
	if q.getNFTCategoriesListStmt != nil {
		if cerr := q.getNFTCategoriesListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTCategoriesListStmt: %w", cerr)
		}
	}
	if q.getNFTCategoryByIDStmt != nil {
		if cerr := q.getNFTCategoryByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTCategoryByIDStmt: %w", cerr)
		}
	}
	if q.getNFTItemByIDStmt != nil {
		if cerr := q.getNFTItemByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTItemByIDStmt: %w", cerr)
		}
	}
	if q.getNFTItemsListStmt != nil {
		if cerr := q.getNFTItemsListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTItemsListStmt: %w", cerr)
		}
	}
	if q.getNFTItemsListByOwnerIDStmt != nil {
		if cerr := q.getNFTItemsListByOwnerIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTItemsListByOwnerIDStmt: %w", cerr)
		}
	}
	if q.getNFTItemsListByRelationIDStmt != nil {
		if cerr := q.getNFTItemsListByRelationIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFTItemsListByRelationIDStmt: %w", cerr)
		}
	}
	if q.resetMainNFTCategoryStmt != nil {
		if cerr := q.resetMainNFTCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing resetMainNFTCategoryStmt: %w", cerr)
		}
	}
	if q.updateNFTCategoryStmt != nil {
		if cerr := q.updateNFTCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFTCategoryStmt: %w", cerr)
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
	db                              DBTX
	tx                              *sql.Tx
	addNFTCategoryStmt              *sql.Stmt
	addNFTItemStmt                  *sql.Stmt
	addNFTItemOwnerStmt             *sql.Stmt
	addNFTRelationStmt              *sql.Stmt
	deleteNFTCategoryByIDStmt       *sql.Stmt
	deleteNFTRelationStmt           *sql.Stmt
	doesUserOwnNFTStmt              *sql.Stmt
	getMainNFTCategoryStmt          *sql.Stmt
	getNFTCategoriesListStmt        *sql.Stmt
	getNFTCategoryByIDStmt          *sql.Stmt
	getNFTItemByIDStmt              *sql.Stmt
	getNFTItemsListStmt             *sql.Stmt
	getNFTItemsListByOwnerIDStmt    *sql.Stmt
	getNFTItemsListByRelationIDStmt *sql.Stmt
	resetMainNFTCategoryStmt        *sql.Stmt
	updateNFTCategoryStmt           *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                              tx,
		tx:                              tx,
		addNFTCategoryStmt:              q.addNFTCategoryStmt,
		addNFTItemStmt:                  q.addNFTItemStmt,
		addNFTItemOwnerStmt:             q.addNFTItemOwnerStmt,
		addNFTRelationStmt:              q.addNFTRelationStmt,
		deleteNFTCategoryByIDStmt:       q.deleteNFTCategoryByIDStmt,
		deleteNFTRelationStmt:           q.deleteNFTRelationStmt,
		doesUserOwnNFTStmt:              q.doesUserOwnNFTStmt,
		getMainNFTCategoryStmt:          q.getMainNFTCategoryStmt,
		getNFTCategoriesListStmt:        q.getNFTCategoriesListStmt,
		getNFTCategoryByIDStmt:          q.getNFTCategoryByIDStmt,
		getNFTItemByIDStmt:              q.getNFTItemByIDStmt,
		getNFTItemsListStmt:             q.getNFTItemsListStmt,
		getNFTItemsListByOwnerIDStmt:    q.getNFTItemsListByOwnerIDStmt,
		getNFTItemsListByRelationIDStmt: q.getNFTItemsListByRelationIDStmt,
		resetMainNFTCategoryStmt:        q.resetMainNFTCategoryStmt,
		updateNFTCategoryStmt:           q.updateNFTCategoryStmt,
	}
}
