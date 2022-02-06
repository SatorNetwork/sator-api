package challenge

import (
	"database/sql"

	"github.com/google/uuid"

	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
)

type questionWithOptions struct {
	question challengeRepo.AddQuestionParams
	options  []challengeRepo.AddQuestionOptionParams
}

var defaultChallenge = challengeRepo.AddChallengeParams{
	ShowID: uuid.MustParse("0899a736-e40d-4ba1-b17e-94857e0c8ff0"),
	Title:  "Miscellaneous",
	Description: sql.NullString{
		String: "Could you be the biggest Friends fan?",
		Valid:  true,
	},
	PrizePool:      250,
	PlayersToStart: 2,
	TimePerQuestion: sql.NullInt32{
		Int32: 1,
		Valid: true,
	},
	Kind:            0,
	UserMaxAttempts: 2,
	MaxWinners: sql.NullInt32{
		Int32: 2,
		Valid: true,
	},
	QuestionsPerGame:  5,
	MinCorrectAnswers: 1,
}

func getDefaultQuestionsWithOptions(challengeID uuid.UUID) []questionWithOptions {
	return []questionWithOptions{
		{
			question: challengeRepo.AddQuestionParams{
				ChallengeID: challengeID,
				Question:    "Joey played Dr. Drake Ramoray on which soap opera show?",
			},
			options: []challengeRepo.AddQuestionOptionParams{
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Santa Barbara",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Days of Our Lives",
					IsCorrect: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Neighbours",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "General Hospital",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
			},
		},
		{
			question: challengeRepo.AddQuestionParams{
				ChallengeID: challengeID,
				Question:    "What store does Phoebe hate?",
			},
			options: []challengeRepo.AddQuestionOptionParams{
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Walmart",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Pottery Barn",
					IsCorrect: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Amazon",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Costco",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
			},
		},
		{
			question: challengeRepo.AddQuestionParams{
				ChallengeID: challengeID,
				Question:    "Rachel got a job with which company in Paris?",
			},
			options: []challengeRepo.AddQuestionOptionParams{
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Louis Vuitton",
					IsCorrect: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Cartier",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Gucci",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Zara",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
			},
		},
		{
			question: challengeRepo.AddQuestionParams{
				ChallengeID: challengeID,
				Question:    "Phoebe’s scientist boyfriend David worked in what city?",
			},
			options: []challengeRepo.AddQuestionOptionParams{
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Minsk",
					IsCorrect: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Berlin",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Moscow",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Kyiv",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
			},
		},
		{
			question: challengeRepo.AddQuestionParams{
				ChallengeID: challengeID,
				Question:    "Monica dated an ophthalmologist named?",
			},
			options: []challengeRepo.AddQuestionOptionParams{
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Rafael",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Chandler",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Robert",
					IsCorrect: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				{
					QuestionID:   uuid.UUID{},
					AnswerOption: "Richard",
					IsCorrect: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
			},
		},
	}
}
