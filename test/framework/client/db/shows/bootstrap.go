package shows

import (
	"database/sql"

	"golang.org/x/net/context"

	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

func (db *DB) Bootstrap(ctx context.Context) error {
	const showNum = 10

	for idx := 0; idx < showNum; idx++ {
		_, err := db.showsRepository.AddShow(ctx, shows_repository.AddShowParams{
			Title:         "test-title",
			Cover:         "test-cover",
			HasNewEpisode: false,
			Category: sql.NullString{
				String: "test-category",
				Valid:  true,
			},
			Description: sql.NullString{
				String: "test-description",
				Valid:  true,
			},
			RealmsTitle: sql.NullString{
				String: "test-realms-title",
				Valid:  true,
			},
			RealmsSubtitle: sql.NullString{
				String: "test-realms-subtitle",
				Valid:  true,
			},
			Watch: sql.NullString{
				String: "test-watch",
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
