package challenge

import (
	"context"

	"github.com/google/uuid"
)

func (db *DB) Bootstrap(ctx context.Context) error {
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

	return nil
}

func (db *DB) DefaultChallengeID(ctx context.Context) (uuid.UUID, error) {
	challenge, err := db.challengeRepository.GetChallengeByTitle(ctx, "Miscellaneous")
	if err != nil {
		return uuid.UUID{}, err
	}

	return challenge.ID, nil
}
