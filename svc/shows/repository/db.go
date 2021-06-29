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
	if q.addEpisodeStmt, err = db.PrepareContext(ctx, addEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query AddEpisode: %w", err)
	}
	if q.addShowStmt, err = db.PrepareContext(ctx, addShow); err != nil {
		return nil, fmt.Errorf("error preparing query AddShow: %w", err)
	}
	if q.deleteEpisodeByIDStmt, err = db.PrepareContext(ctx, deleteEpisodeByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEpisodeByID: %w", err)
	}
	if q.deleteEpisodeByShowIDStmt, err = db.PrepareContext(ctx, deleteEpisodeByShowID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEpisodeByShowID: %w", err)
	}
	if q.deleteShowByIDStmt, err = db.PrepareContext(ctx, deleteShowByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteShowByID: %w", err)
	}
	if q.getEpisodeByIDStmt, err = db.PrepareContext(ctx, getEpisodeByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeByID: %w", err)
	}
	if q.getEpisodesByShowIDStmt, err = db.PrepareContext(ctx, getEpisodesByShowID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodesByShowID: %w", err)
	}
	if q.getShowByIDStmt, err = db.PrepareContext(ctx, getShowByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetShowByID: %w", err)
	}
	if q.getShowsStmt, err = db.PrepareContext(ctx, getShows); err != nil {
		return nil, fmt.Errorf("error preparing query GetShows: %w", err)
	}
	if q.getShowsByCategoryStmt, err = db.PrepareContext(ctx, getShowsByCategory); err != nil {
		return nil, fmt.Errorf("error preparing query GetShowsByCategory: %w", err)
	}
	if q.updateEpisodeStmt, err = db.PrepareContext(ctx, updateEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateEpisode: %w", err)
	}
	if q.updateShowStmt, err = db.PrepareContext(ctx, updateShow); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateShow: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addEpisodeStmt != nil {
		if cerr := q.addEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addEpisodeStmt: %w", cerr)
		}
	}
	if q.addShowStmt != nil {
		if cerr := q.addShowStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addShowStmt: %w", cerr)
		}
	}
	if q.deleteEpisodeByIDStmt != nil {
		if cerr := q.deleteEpisodeByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEpisodeByIDStmt: %w", cerr)
		}
	}
	if q.deleteEpisodeByShowIDStmt != nil {
		if cerr := q.deleteEpisodeByShowIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEpisodeByShowIDStmt: %w", cerr)
		}
	}
	if q.deleteShowByIDStmt != nil {
		if cerr := q.deleteShowByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteShowByIDStmt: %w", cerr)
		}
	}
	if q.getEpisodeByIDStmt != nil {
		if cerr := q.getEpisodeByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeByIDStmt: %w", cerr)
		}
	}
	if q.getEpisodesByShowIDStmt != nil {
		if cerr := q.getEpisodesByShowIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodesByShowIDStmt: %w", cerr)
		}
	}
	if q.getShowByIDStmt != nil {
		if cerr := q.getShowByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getShowByIDStmt: %w", cerr)
		}
	}
	if q.getShowsStmt != nil {
		if cerr := q.getShowsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getShowsStmt: %w", cerr)
		}
	}
	if q.getShowsByCategoryStmt != nil {
		if cerr := q.getShowsByCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getShowsByCategoryStmt: %w", cerr)
		}
	}
	if q.updateEpisodeStmt != nil {
		if cerr := q.updateEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateEpisodeStmt: %w", cerr)
		}
	}
	if q.updateShowStmt != nil {
		if cerr := q.updateShowStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateShowStmt: %w", cerr)
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
	db                        DBTX
	tx                        *sql.Tx
	addEpisodeStmt            *sql.Stmt
	addShowStmt               *sql.Stmt
	deleteEpisodeByIDStmt     *sql.Stmt
	deleteEpisodeByShowIDStmt *sql.Stmt
	deleteShowByIDStmt        *sql.Stmt
	getEpisodeByIDStmt        *sql.Stmt
	getEpisodesByShowIDStmt   *sql.Stmt
	getShowByIDStmt           *sql.Stmt
	getShowsStmt              *sql.Stmt
	getShowsByCategoryStmt    *sql.Stmt
	updateEpisodeStmt         *sql.Stmt
	updateShowStmt            *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                        tx,
		tx:                        tx,
		addEpisodeStmt:            q.addEpisodeStmt,
		addShowStmt:               q.addShowStmt,
		deleteEpisodeByIDStmt:     q.deleteEpisodeByIDStmt,
		deleteEpisodeByShowIDStmt: q.deleteEpisodeByShowIDStmt,
		deleteShowByIDStmt:        q.deleteShowByIDStmt,
		getEpisodeByIDStmt:        q.getEpisodeByIDStmt,
		getEpisodesByShowIDStmt:   q.getEpisodesByShowIDStmt,
		getShowByIDStmt:           q.getShowByIDStmt,
		getShowsStmt:              q.getShowsStmt,
		getShowsByCategoryStmt:    q.getShowsByCategoryStmt,
		updateEpisodeStmt:         q.updateEpisodeStmt,
		updateShowStmt:            q.updateShowStmt,
	}
}
