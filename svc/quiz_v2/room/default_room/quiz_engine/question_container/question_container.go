package question_container

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	quiz_v2_challenge "github.com/SatorNetwork/sator-api/svc/quiz_v2/challenge"
)

type QuestionContainer interface {
	GetChallenge() *challenge.RawChallenge
	GetNumberOfQuestions() int
	GetQuestions() []challenge.Question
	GetQuestionByID(questionID uuid.UUID) (challenge.Question, error)
	GetQuestionNumByID(questionID uuid.UUID) (int, error)
	CheckAnswer(questionID, answerID uuid.UUID) (bool, error)
}

type questionContainer struct {
	challenge *challenge.RawChallenge
	questions []challenge.Question
}

func New(challengeID string, challengesSvc quiz_v2_challenge.ChallengesService) (*questionContainer, error) {
	challenge, err := loadChallenge(challengeID, challengesSvc)
	if err != nil {
		return nil, err
	}
	questions, err := loadQuestions(challengeID, challengesSvc)
	if err != nil {
		return nil, err
	}

	return &questionContainer{
		challenge: challenge,
		questions: questions,
	}, nil
}

func loadChallenge(challengeID string, challengesSvc quiz_v2_challenge.ChallengesService) (*challenge.RawChallenge, error) {
	ctx := context.Background()
	challengeUID, err := uuid.Parse(challengeID)
	if err != nil {
		return nil, err
	}
	challenge, err := challengesSvc.GetRawChallengeByID(ctx, challengeUID)
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}

func loadQuestions(challengeID string, challengesSvc quiz_v2_challenge.ChallengesService) ([]challenge.Question, error) {
	ctx := context.Background()
	challengeUID, err := uuid.Parse(challengeID)
	if err != nil {
		return nil, err
	}
	questions, err := challengesSvc.GetQuestionsByChallengeID(ctx, challengeUID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get questions by challenge id")
	}

	return questions, nil
}

func (e *questionContainer) GetChallenge() *challenge.RawChallenge {
	return e.challenge
}

func (e *questionContainer) GetNumberOfQuestions() int {
	return len(e.questions)
}

func (e *questionContainer) GetQuestions() []challenge.Question {
	return e.questions
}

func (e *questionContainer) GetQuestionByID(questionID uuid.UUID) (challenge.Question, error) {
	for _, q := range e.questions {
		if q.ID == questionID {
			return q, nil
		}
	}

	return challenge.Question{}, errors.Errorf("question not found")
}

func (e *questionContainer) GetQuestionNumByID(questionID uuid.UUID) (int, error) {
	question, err := e.GetQuestionByID(questionID)
	if err != nil {
		return 0, err
	}

	return int(question.Order) - 1, nil
}

func (e *questionContainer) CheckAnswer(questionID, answerID uuid.UUID) (bool, error) {
	question, err := e.GetQuestionByID(questionID)
	if err != nil {
		return false, err
	}

	answer, err := e.getAnswerByID(question, answerID)
	if err != nil {
		return false, err
	}

	return answer.IsCorrect, nil
}

func (e *questionContainer) getAnswerByID(question challenge.Question, answerID uuid.UUID) (challenge.AnswerOption, error) {
	for _, answer := range question.AnswerOptions {
		if answer.ID == answerID {
			return answer, nil
		}
	}

	return challenge.AnswerOption{}, errors.Errorf("answer not found")
}
