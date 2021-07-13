package questions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Predefined request query keys
const (
	pageParam         = "page"
	itemsPerPageParam = "items_per_page"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Post("/", httptransport.NewServer(
		e.AddQuestion,
		decodeAddQuestionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{question_id}/answers", httptransport.NewServer(
		e.AddQuestionOption,
		decodeAddQuestionOptionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{question_id}", httptransport.NewServer(
		e.DeleteQuestionByID,
		decodeDeleteQuestionByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{question_id}/answers/{answer_id}", httptransport.NewServer(
		e.DeleteAnswerByID,
		decodeDeleteAnswerByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{question_id}", httptransport.NewServer(
		e.UpdateQuestion,
		decodeUpdateQuestionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{question_id}/answers/{answer_id}", httptransport.NewServer(
		e.UpdateAnswer,
		decodeUpdateAnswerRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{question_id}", httptransport.NewServer(
		e.GetQuestionByID,
		decodeGetQuestionByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/challenge/{challenge_id}", httptransport.NewServer(
		e.GetQuestionsByChallengeID,
		decodeGetQuestionsByChallengeIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeAddQuestionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeAddQuestionOptionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AnswerOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	questionID := chi.URLParam(r, "question_id")
	if questionID == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}
	req.QuestionID = questionID

	return req, nil
}

func decodeDeleteQuestionByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "question_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}

	return id, nil

}

func decodeDeleteAnswerByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	questionID := chi.URLParam(r, "question_id")
	if questionID == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}

	answerID := chi.URLParam(r, "answer_id")
	if answerID == "" {
		return nil, fmt.Errorf("%w: missed answer id", ErrInvalidParameter)
	}

	return DeleteAnswerByIDRequest{
		AnswerID:   answerID,
		QuestionID: questionID,
	}, nil

}

func decodeUpdateQuestionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	questionID := chi.URLParam(r, "question_id")
	if questionID == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}
	req.ID = questionID

	return req, nil
}

func decodeUpdateAnswerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	questionID := chi.URLParam(r, "question_id")
	if questionID == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}
	req.QuestionID = questionID

	answerID := chi.URLParam(r, "answer_id")
	if answerID == "" {
		return nil, fmt.Errorf("%w: missed answer id", ErrInvalidParameter)
	}
	req.ID = answerID

	return req, nil
}

func decodeGetQuestionByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "question_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed question id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetQuestionsByChallengeIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "challenge_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}

	return id, nil
}

func castStrToInt32(source string) int32 {
	res, err := strconv.Atoi(source)
	if err != nil {
		return 0
	}
	return int32(res)
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}
	return httpencoder.CodeAndMessageFrom(err)
}
