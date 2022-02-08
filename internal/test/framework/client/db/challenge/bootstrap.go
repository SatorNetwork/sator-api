package challenge

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"

	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
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
	}

	{
		defaultChallengeCopy := deepcopy.Copy(defaultChallenge).(challengeRepo.AddChallengeParams)
		defaultChallengeCopy.Title = "custom1"
		defaultChallengeCopy.PlayersToStart = 3
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
