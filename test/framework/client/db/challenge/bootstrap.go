package challenge

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"

	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

func (db *DB) Bootstrap(ctx context.Context) error {
	{
		challenge, err := db.challengeRepository.AddChallenge(ctx, defaultChallenge)
		if err != nil {
			return err
		}

		questionsWithOptions := getDefaultQuestionsWithOptions(challenge.ID)
		for _, q := range questionsWithOptions {
			addQuestionResp, err := db.challengeRepository.AddQuestion(ctx, q.question)
			if err != nil {
				return err
			}

			for _, option := range q.options {
				option.QuestionID = addQuestionResp.ID
				if _, err := db.challengeRepository.AddQuestionOption(ctx, option); err != nil {
					return err
				}
			}
		}

		show, err := db.addShow(ctx)
		if err != nil {
			return err
		}
		err = db.addEpisode(ctx, show.ID, challenge.ID)
		if err != nil {
			return err
		}
	}

	{
		defaultChallengeCopy := deepcopy.Copy(defaultChallenge).(challengeRepo.AddChallengeParams)
		defaultChallengeCopy.Title = "custom1"
		defaultChallengeCopy.PlayersToStart = 3
		defaultChallengeCopy.QuestionsPerGame = 1
		challenge, err := db.challengeRepository.AddChallenge(ctx, defaultChallengeCopy)
		if err != nil {
			return errors.Wrapf(err, "can't add %v challenge", defaultChallengeCopy.Title)
		}

		questionsWithOptions := getDefaultQuestionsWithOptions(challenge.ID)[:1]
		for _, q := range questionsWithOptions {
			addQuestionResp, err := db.challengeRepository.AddQuestion(ctx, q.question)
			if err != nil {
				return errors.Wrap(err, "can't add question")
			}

			for _, option := range q.options {
				option.QuestionID = addQuestionResp.ID
				if _, err := db.challengeRepository.AddQuestionOption(ctx, option); err != nil {
					return errors.Wrap(err, "can't add question option")
				}
			}
		}

		show, err := db.addShow(ctx)
		if err != nil {
			return err
		}
		err = db.addEpisode(ctx, show.ID, challenge.ID)
		if err != nil {
			return err
		}
	}

	{
		defaultChallengeCopy := deepcopy.Copy(defaultChallenge).(challengeRepo.AddChallengeParams)
		defaultChallengeCopy.Title = "custom2"
		defaultChallengeCopy.PlayersToStart = 10
		defaultChallengeCopy.QuestionsPerGame = 1
		challenge, err := db.challengeRepository.AddChallenge(ctx, defaultChallengeCopy)
		if err != nil {
			return errors.Wrapf(err, "can't add %v challenge", defaultChallengeCopy.Title)
		}

		questionsWithOptions := getDefaultQuestionsWithOptions(challenge.ID)[:1]
		for _, q := range questionsWithOptions {
			addQuestionResp, err := db.challengeRepository.AddQuestion(ctx, q.question)
			if err != nil {
				return errors.Wrap(err, "can't add question")
			}

			for _, option := range q.options {
				option.QuestionID = addQuestionResp.ID
				if _, err := db.challengeRepository.AddQuestionOption(ctx, option); err != nil {
					return errors.Wrap(err, "can't add question option")
				}
			}
		}

		show, err := db.addShow(ctx)
		if err != nil {
			return err
		}
		err = db.addEpisode(ctx, show.ID, challenge.ID)
		if err != nil {
			return err
		}
	}

	{
		defaultChallengeCopy := deepcopy.Copy(defaultChallenge).(challengeRepo.AddChallengeParams)
		defaultChallengeCopy.Title = "new-prize-pool-logic"
		defaultChallengeCopy.PercentForQuiz = 5
		challenge, err := db.challengeRepository.AddChallenge(ctx, defaultChallengeCopy)
		if err != nil {
			return errors.Wrapf(err, "can't add %v challenge", defaultChallengeCopy.Title)
		}

		questionsWithOptions := getDefaultQuestionsWithOptions(challenge.ID)
		for _, q := range questionsWithOptions {
			addQuestionResp, err := db.challengeRepository.AddQuestion(ctx, q.question)
			if err != nil {
				return err
			}

			for _, option := range q.options {
				option.QuestionID = addQuestionResp.ID
				if _, err := db.challengeRepository.AddQuestionOption(ctx, option); err != nil {
					return err
				}
			}
		}

		show, err := db.addShow(ctx)
		if err != nil {
			return err
		}
		err = db.addEpisode(ctx, show.ID, challenge.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) addShow(ctx context.Context) (*shows_repository.Show, error) {
	show, err := db.showsRepository.AddShow(ctx, shows_repository.AddShowParams{
		Title:          "",
		Cover:          "",
		HasNewEpisode:  false,
		Category:       sql.NullString{},
		Description:    sql.NullString{},
		RealmsTitle:    sql.NullString{},
		RealmsSubtitle: sql.NullString{},
		Watch:          sql.NullString{},
	})
	if err != nil {
		return nil, err
	}

	return &show, nil
}

func (db *DB) addEpisode(ctx context.Context, showID, challengeID uuid.UUID) error {
	_, err := db.showsRepository.AddEpisode(ctx, shows_repository.AddEpisodeParams{
		ShowID:        showID,
		SeasonID:      uuid.NullUUID{},
		EpisodeNumber: 0,
		Cover:         sql.NullString{},
		Title:         "",
		Description:   sql.NullString{},
		ReleaseDate:   sql.NullTime{},
		ChallengeID: uuid.NullUUID{
			UUID:  challengeID,
			Valid: true,
		},
		VerificationChallengeID: uuid.NullUUID{},
		HintText:                sql.NullString{},
		Watch:                   sql.NullString{},
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DefaultChallengeID(ctx context.Context) (uuid.UUID, error) {
	challenge, err := db.challengeRepository.GetChallengeByTitle(ctx, "Miscellaneous")
	if err != nil {
		return uuid.UUID{}, err
	}

	return challenge.ID, nil
}

func (db *DB) ChallengeIDByTitle(ctx context.Context, title string) (uuid.UUID, error) {
	challenge, err := db.challengeRepository.GetChallengeByTitle(ctx, title)
	if err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "can't get challenge by %v title", title)
	}

	return challenge.ID, nil
}
