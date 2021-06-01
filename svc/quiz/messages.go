package quiz

import "time"

// Predefined message types
const (
	UserConnectedMessage   = "player_connected"
	UserDisonnectedMessage = "player_disconnected"
	CountdownMessage       = "countdown"
	QuestionMessage        = "question"
	AnswerMessage          = "answer"
	QuestionResultMessage  = "question_result"
	ChallengeResultMessage = "challenge_result"
)

type (
	// Websocket message structure
	Message struct {
		Type    string      `json:"type"`
		SentAt  time.Time   `json:"sent_at"`
		Payload interface{} `json:"payload"`
	}

	User struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
	}

	Countdown struct {
		Countdown int `json:"countdown"`
	}

	Question struct {
		QuestionID     string         `json:"question_id"`
		QuestionText   string         `json:"question_text"`
		TimeForAnswer  int            `json:"time_for_answer"`
		TotalQuestions int            `json:"total_questions"`
		QuestionNumber int            `json:"question_number"`
		AnswerOptions  []AnswerOption `json:"answer_options"`
		correctID      string
	}

	AnswerOption struct {
		AnswerID   string `json:"answer_id"`
		AnswerText string `json:"answer_text"`
	}

	Answer struct {
		QuestionID string `json:"question_id"`
		AnswerID   string `json:"answer_id"`
	}

	MessageAnswer struct {
		Type    string    `json:"type"`
		SentAt  time.Time `json:"sent_at"`
		Payload Answer    `json:"payload"`
	}

	QuestionResult struct {
		QuestionID      string `json:"question_id"`
		Result          bool   `json:"result"`
		Rate            int    `json:"rate"`
		CorrectAnswerID string `json:"correct_answer_id"`
		QuestionsLeft   int    `json:"questions_left"`
		AdditionalPts   int    `json:"additional_pts"`
	}

	ChallengeResult struct {
		ChallengeID string   `json:"challenge_id"`
		PrizePool   string   `json:"prize_pool"`
		Winners     []Winner `json:"winners"`
	}

	Winner struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Prize    string `json:"prize"`
	}
)
