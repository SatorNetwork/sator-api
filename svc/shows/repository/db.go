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
	if q.addClapForShowStmt, err = db.PrepareContext(ctx, addClapForShow); err != nil {
		return nil, fmt.Errorf("error preparing query AddClapForShow: %w", err)
	}
	if q.addEpisodeStmt, err = db.PrepareContext(ctx, addEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query AddEpisode: %w", err)
	}
	if q.addSeasonStmt, err = db.PrepareContext(ctx, addSeason); err != nil {
		return nil, fmt.Errorf("error preparing query AddSeason: %w", err)
	}
	if q.addShowStmt, err = db.PrepareContext(ctx, addShow); err != nil {
		return nil, fmt.Errorf("error preparing query AddShow: %w", err)
	}
	if q.countUserClapsStmt, err = db.PrepareContext(ctx, countUserClaps); err != nil {
		return nil, fmt.Errorf("error preparing query CountUserClaps: %w", err)
	}
	if q.deleteEpisodeByIDStmt, err = db.PrepareContext(ctx, deleteEpisodeByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEpisodeByID: %w", err)
	}
	if q.deleteReviewStmt, err = db.PrepareContext(ctx, deleteReview); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteReview: %w", err)
	}
	if q.deleteSeasonByIDStmt, err = db.PrepareContext(ctx, deleteSeasonByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSeasonByID: %w", err)
	}
	if q.deleteShowByIDStmt, err = db.PrepareContext(ctx, deleteShowByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteShowByID: %w", err)
	}
	if q.didUserRateEpisodeStmt, err = db.PrepareContext(ctx, didUserRateEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query DidUserRateEpisode: %w", err)
	}
	if q.didUserReviewEpisodeStmt, err = db.PrepareContext(ctx, didUserReviewEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query DidUserReviewEpisode: %w", err)
	}
	if q.getEpisodeByIDStmt, err = db.PrepareContext(ctx, getEpisodeByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeByID: %w", err)
	}
	if q.getEpisodeIDByVerificationChallengeIDStmt, err = db.PrepareContext(ctx, getEpisodeIDByVerificationChallengeID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeIDByVerificationChallengeID: %w", err)
	}
	if q.getEpisodeRatingByIDStmt, err = db.PrepareContext(ctx, getEpisodeRatingByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeRatingByID: %w", err)
	}
	if q.getEpisodesByShowIDStmt, err = db.PrepareContext(ctx, getEpisodesByShowID); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodesByShowID: %w", err)
	}
	if q.getListEpisodesByIDsStmt, err = db.PrepareContext(ctx, getListEpisodesByIDs); err != nil {
		return nil, fmt.Errorf("error preparing query GetListEpisodesByIDs: %w", err)
	}
	if q.getReviewByIDStmt, err = db.PrepareContext(ctx, getReviewByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetReviewByID: %w", err)
	}
	if q.getSeasonByIDStmt, err = db.PrepareContext(ctx, getSeasonByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetSeasonByID: %w", err)
	}
	if q.getSeasonsByShowIDStmt, err = db.PrepareContext(ctx, getSeasonsByShowID); err != nil {
		return nil, fmt.Errorf("error preparing query GetSeasonsByShowID: %w", err)
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
	if q.rateEpisodeStmt, err = db.PrepareContext(ctx, rateEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query RateEpisode: %w", err)
	}
	if q.reviewEpisodeStmt, err = db.PrepareContext(ctx, reviewEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query ReviewEpisode: %w", err)
	}
	if q.reviewsListStmt, err = db.PrepareContext(ctx, reviewsList); err != nil {
		return nil, fmt.Errorf("error preparing query ReviewsList: %w", err)
	}
	if q.reviewsListByUserIDStmt, err = db.PrepareContext(ctx, reviewsListByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query ReviewsListByUserID: %w", err)
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
	if q.addClapForShowStmt != nil {
		if cerr := q.addClapForShowStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addClapForShowStmt: %w", cerr)
		}
	}
	if q.addEpisodeStmt != nil {
		if cerr := q.addEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addEpisodeStmt: %w", cerr)
		}
	}
	if q.addSeasonStmt != nil {
		if cerr := q.addSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addSeasonStmt: %w", cerr)
		}
	}
	if q.addShowStmt != nil {
		if cerr := q.addShowStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addShowStmt: %w", cerr)
		}
	}
	if q.countUserClapsStmt != nil {
		if cerr := q.countUserClapsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countUserClapsStmt: %w", cerr)
		}
	}
	if q.deleteEpisodeByIDStmt != nil {
		if cerr := q.deleteEpisodeByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEpisodeByIDStmt: %w", cerr)
		}
	}
	if q.deleteReviewStmt != nil {
		if cerr := q.deleteReviewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteReviewStmt: %w", cerr)
		}
	}
	if q.deleteSeasonByIDStmt != nil {
		if cerr := q.deleteSeasonByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSeasonByIDStmt: %w", cerr)
		}
	}
	if q.deleteShowByIDStmt != nil {
		if cerr := q.deleteShowByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteShowByIDStmt: %w", cerr)
		}
	}
	if q.didUserRateEpisodeStmt != nil {
		if cerr := q.didUserRateEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing didUserRateEpisodeStmt: %w", cerr)
		}
	}
	if q.didUserReviewEpisodeStmt != nil {
		if cerr := q.didUserReviewEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing didUserReviewEpisodeStmt: %w", cerr)
		}
	}
	if q.getEpisodeByIDStmt != nil {
		if cerr := q.getEpisodeByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeByIDStmt: %w", cerr)
		}
	}
	if q.getEpisodeIDByVerificationChallengeIDStmt != nil {
		if cerr := q.getEpisodeIDByVerificationChallengeIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeIDByVerificationChallengeIDStmt: %w", cerr)
		}
	}
	if q.getEpisodeRatingByIDStmt != nil {
		if cerr := q.getEpisodeRatingByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeRatingByIDStmt: %w", cerr)
		}
	}
	if q.getEpisodesByShowIDStmt != nil {
		if cerr := q.getEpisodesByShowIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodesByShowIDStmt: %w", cerr)
		}
	}
	if q.getListEpisodesByIDsStmt != nil {
		if cerr := q.getListEpisodesByIDsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getListEpisodesByIDsStmt: %w", cerr)
		}
	}
	if q.getReviewByIDStmt != nil {
		if cerr := q.getReviewByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getReviewByIDStmt: %w", cerr)
		}
	}
	if q.getSeasonByIDStmt != nil {
		if cerr := q.getSeasonByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSeasonByIDStmt: %w", cerr)
		}
	}
	if q.getSeasonsByShowIDStmt != nil {
		if cerr := q.getSeasonsByShowIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSeasonsByShowIDStmt: %w", cerr)
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
	if q.rateEpisodeStmt != nil {
		if cerr := q.rateEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing rateEpisodeStmt: %w", cerr)
		}
	}
	if q.reviewEpisodeStmt != nil {
		if cerr := q.reviewEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing reviewEpisodeStmt: %w", cerr)
		}
	}
	if q.reviewsListStmt != nil {
		if cerr := q.reviewsListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing reviewsListStmt: %w", cerr)
		}
	}
	if q.reviewsListByUserIDStmt != nil {
		if cerr := q.reviewsListByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing reviewsListByUserIDStmt: %w", cerr)
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
	db                                        DBTX
	tx                                        *sql.Tx
	addClapForShowStmt                        *sql.Stmt
	addEpisodeStmt                            *sql.Stmt
	addSeasonStmt                             *sql.Stmt
	addShowStmt                               *sql.Stmt
	countUserClapsStmt                        *sql.Stmt
	deleteEpisodeByIDStmt                     *sql.Stmt
	deleteReviewStmt                          *sql.Stmt
	deleteSeasonByIDStmt                      *sql.Stmt
	deleteShowByIDStmt                        *sql.Stmt
	didUserRateEpisodeStmt                    *sql.Stmt
	didUserReviewEpisodeStmt                  *sql.Stmt
	getEpisodeByIDStmt                        *sql.Stmt
	getEpisodeIDByVerificationChallengeIDStmt *sql.Stmt
	getEpisodeRatingByIDStmt                  *sql.Stmt
	getEpisodesByShowIDStmt                   *sql.Stmt
	getListEpisodesByIDsStmt                  *sql.Stmt
	getReviewByIDStmt                         *sql.Stmt
	getSeasonByIDStmt                         *sql.Stmt
	getSeasonsByShowIDStmt                    *sql.Stmt
	getShowByIDStmt                           *sql.Stmt
	getShowsStmt                              *sql.Stmt
	getShowsByCategoryStmt                    *sql.Stmt
	rateEpisodeStmt                           *sql.Stmt
	reviewEpisodeStmt                         *sql.Stmt
	reviewsListStmt                           *sql.Stmt
	reviewsListByUserIDStmt                   *sql.Stmt
	updateEpisodeStmt                         *sql.Stmt
	updateShowStmt                            *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                       tx,
		tx:                       tx,
		addClapForShowStmt:       q.addClapForShowStmt,
		addEpisodeStmt:           q.addEpisodeStmt,
		addSeasonStmt:            q.addSeasonStmt,
		addShowStmt:              q.addShowStmt,
		countUserClapsStmt:       q.countUserClapsStmt,
		deleteEpisodeByIDStmt:    q.deleteEpisodeByIDStmt,
		deleteReviewStmt:         q.deleteReviewStmt,
		deleteSeasonByIDStmt:     q.deleteSeasonByIDStmt,
		deleteShowByIDStmt:       q.deleteShowByIDStmt,
		didUserRateEpisodeStmt:   q.didUserRateEpisodeStmt,
		didUserReviewEpisodeStmt: q.didUserReviewEpisodeStmt,
		getEpisodeByIDStmt:       q.getEpisodeByIDStmt,
		getEpisodeIDByVerificationChallengeIDStmt: q.getEpisodeIDByVerificationChallengeIDStmt,
		getEpisodeRatingByIDStmt:                  q.getEpisodeRatingByIDStmt,
		getEpisodesByShowIDStmt:                   q.getEpisodesByShowIDStmt,
		getListEpisodesByIDsStmt:                  q.getListEpisodesByIDsStmt,
		getReviewByIDStmt:                         q.getReviewByIDStmt,
		getSeasonByIDStmt:                         q.getSeasonByIDStmt,
		getSeasonsByShowIDStmt:                    q.getSeasonsByShowIDStmt,
		getShowByIDStmt:                           q.getShowByIDStmt,
		getShowsStmt:                              q.getShowsStmt,
		getShowsByCategoryStmt:                    q.getShowsByCategoryStmt,
		rateEpisodeStmt:                           q.rateEpisodeStmt,
		reviewEpisodeStmt:                         q.reviewEpisodeStmt,
		reviewsListStmt:                           q.reviewsListStmt,
		reviewsListByUserIDStmt:                   q.reviewsListByUserIDStmt,
		updateEpisodeStmt:                         q.updateEpisodeStmt,
		updateShowStmt:                            q.updateShowStmt,
	}
}
