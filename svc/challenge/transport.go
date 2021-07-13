package challenge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
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

	r.Get("/{episode_id}/validation-question", httptransport.NewServer(
		e.GetVerificationQuestionByEpisodeID,
		decodeGetValidationQuestionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{question_id}/check-answer/{answer_id}", httptransport.NewServer(
		e.CheckVerificationQuestionAnswer,
		decodeCheckAnswerRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{episode_id}/is-activated", httptransport.NewServer(
		e.VerifyUserAccessToEpisode,
		decodeIsEpisodeActivatedRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{id}", httptransport.NewServer(
		e.GetChallengeById,
		decodeGetChallengeByIdRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/", httptransport.NewServer(
		e.AddChallenge,
		decodeAddChallengeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{id}", httptransport.NewServer(
		e.DeleteChallengeByID,
		decodeDeleteChallengeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{id}", httptransport.NewServer(
		e.UpdateChallenge,
		decodeUpdateChallengeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetValidationQuestionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "episode_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode_id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeCheckAnswerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	qid := chi.URLParam(r, "question_id")
	if qid == "" {
		return nil, fmt.Errorf("%w: missed question_id", ErrInvalidParameter)
	}

	aid := chi.URLParam(r, "answer_id")
	if aid == "" {
		return nil, fmt.Errorf("%w: missed answer_id", ErrInvalidParameter)
	}

	return CheckAnswerRequest{
		QuestionID: qid,
		AnswerID:   aid,
	}, nil
}

func decodeIsEpisodeActivatedRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "episode_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode_id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetChallengeByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}

	return id, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeAddChallengeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeDeleteChallengeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}

	return id, nil

}

func decodeUpdateChallengeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	challengeID := chi.URLParam(r, "id")
	if challengeID == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}
	req.ID = challengeID

	return req, nil
}
