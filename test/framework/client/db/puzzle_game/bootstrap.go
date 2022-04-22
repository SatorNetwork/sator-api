package puzzle_game

import (
	"database/sql"
	"fmt"
	"strings"

	files_repository "github.com/SatorNetwork/sator-api/svc/files/repository"
	puzzle_game_repository "github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var createGameUnlockOptionQuery = `INSERT INTO public.puzzle_game_unlock_options (id, steps, amount, disabled, updated_at, created_at, locked) VALUES
    ('899fcefd-9b4b-4c67-905f-b1e580fcaf78', 32, 0, false, null, '2022-04-21 05:27:55.628799', false);`

func (db *DB) Bootstrap(ctx context.Context) error {
	var show shows_repository.Show
	err := db.dbClient.QueryRowContext(ctx, "SELECT id FROM public.shows WHERE title = $1", "test-title-show").Scan(&show.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			if _, err := db.dbClient.Exec("DELETE FROM public.shows WHERE title = $1", "test-title-show"); err != nil {
				return err
			}
			show, err = db.showsRepository.AddShow(ctx, shows_repository.AddShowParams{
				Title:         "test-title-show",
				Cover:         "test-cover",
				HasNewEpisode: false,
				Category: sql.NullString{
					String: "test-category-pg",
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
		} else {
			return err
		}
	}

	var episode shows_repository.Episode
	if err := db.dbClient.QueryRowContext(ctx, "SELECT id FROM public.episodes WHERE title = $1 AND show_id = $2", "test-episode-ep", show.ID).Scan(&episode.ID); err != nil {
		if err == sql.ErrNoRows {
			if _, err := db.dbClient.Exec("DELETE FROM public.episodes WHERE title = $1", "test-episode-ep"); err != nil {
				return err
			}
			episode, err = db.showsRepository.AddEpisode(ctx, shows_repository.AddEpisodeParams{
				ShowID:        show.ID,
				EpisodeNumber: 1,
				Title:         "test-episode-ep",
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var pg puzzle_game_repository.PuzzleGame
	if err = db.dbClient.QueryRowContext(ctx, "SELECT id FROM public.puzzle_games WHERE episode_id=$1 AND prize_pool = 100 AND parts_x = 4", episode.ID).
		Scan(&pg.ID); err != nil {
		if err == sql.ErrNoRows {
			pg, err = db.puzzleGameRepository.CreatePuzzleGame(ctx, puzzle_game_repository.CreatePuzzleGameParams{
				EpisodeID: episode.ID,
				PrizePool: 100,
				PartsX:    4,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	fileIDs := []string{
		"'b7684f9f-91d6-4683-8489-f16612f28aa5'",
		"'3320cf8c-f6d9-42f1-b8d6-cea843910c43'",
		"'e260ec9c-13a7-42b7-a945-40ef2dbc0303'",
		"'9fe63234-90ba-4aa6-9fb2-62ca738f597d'",
		"'ac2acc14-ce99-4c28-9596-5a26d574e965'",
		"'998baee6-74d4-4898-8a94-32a4e867749e'",
		"'49062388-3ea3-4833-b23a-7d3e6081739d'",
		"'a851b737-3c87-4bc4-8d04-39bcc3c3f34a'",
		"'c8195e1d-83d1-4701-8ebc-0b237040101b'",
		"'194ad696-d7d9-4b97-90ab-f52f514027e5'",
		"'e3693fe0-c157-4742-9781-a96bccc35c35'",
		"'681b7f1d-870a-4367-9b60-e3df60955ead'",
		"'29e7a54b-c86c-4021-ad17-3c27d7cd6b9c'",
		"'311cfc0e-f1bb-4c54-baac-bb82e9964c6e'",
		"'08e18925-bc61-40db-90fe-d60f363c5487'",
		"'a706620b-a179-41f8-9217-830327dea1f5'",
	}

	fileUIDs := make([]uuid.UUID, 16)
	for i, id := range fileIDs {
		fileUIDs[i] = uuid.MustParse(id)
	}

	fileIDsStr := fmt.Sprintf("(%s)", strings.Join(fileIDs, ","))
	var count int
	err = db.dbClient.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(id) FROM public.files WHERE id IN %s", fileIDsStr)).Scan(&count)
	if err != nil {
		return err
	}
	if count != 16 {
		if _, err := db.dbClient.Exec(fmt.Sprintf("DELETE FROM public.files WHERE id IN %s", fileIDsStr)); err != nil {
			return err
		}
		if _, err := db.dbClient.Exec(fmt.Sprintf("DELETE FROM public.puzzle_games_to_images WHERE file_id IN %s", fileIDsStr)); err != nil {
			return err
		}

		for i := 0; i < 16; i++ {
			_, err := db.filesRepository.AddFile(ctx, files_repository.AddFileParams{
				ID:       fileUIDs[i],
				FileName: "test-file-name-pg",
				FilePath: "test-path",
				FileUrl:  "test-url",
			})
			if err != nil {
				return err
			}
			err = db.puzzleGameRepository.LinkImageToPuzzleGame(ctx, puzzle_game_repository.LinkImageToPuzzleGameParams{
				FileID:       fileUIDs[i],
				PuzzleGameID: pg.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	err = db.dbClient.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(file_id) FROM public.puzzle_games_to_images WHERE file_id IN %s AND puzzle_game_id = '%s'", fileIDsStr, pg.ID)).Scan(&count)
	if err != nil {
		return err
	}
	if count != 16 {
		if _, err := db.dbClient.Exec(fmt.Sprintf("DELETE FROM public.puzzle_games_to_images WHERE file_id IN %s", fileIDsStr)); err != nil {
			return err
		}
		for i := 0; i < 16; i++ {
			err = db.puzzleGameRepository.LinkImageToPuzzleGame(ctx, puzzle_game_repository.LinkImageToPuzzleGameParams{
				FileID:       fileUIDs[i],
				PuzzleGameID: pg.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	var optionID string
	err = db.dbClient.QueryRowContext(ctx, "SELECT id FROM public.puzzle_game_unlock_options WHERE id = $1", "899fcefd-9b4b-4c67-905f-b1e580fcaf78").Scan(&optionID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.dbClient.Exec(createGameUnlockOptionQuery)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
