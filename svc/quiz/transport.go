package quiz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Predefined request query keys
const (
	token       = "token"
	userid      = "userid"
	quizid      = "quizid"
	questionidi = "questionidi"
	answerid    = "answerid"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger, quizWsHandler http.HandlerFunc) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/{challenge_id}/play", httptransport.NewServer(
		e.GetQuizLink,
		decodeGetQuizLinkRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{challenge_id}/play/{token}", quizWsHandler)

	r.Get("/quizzes/token", httptransport.NewServer(
		e.ParseQuizToken,
		decodeParseQuizTokenRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/quizzes/answer", httptransport.NewServer(
		e.StoreAnswer,
		decodeStoreAnswerRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/quizzes/{quiz_id}", httptransport.NewServer(
		e.DeleteQuizByID,
		decodeDeleteQuizByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetQuizLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "challenge_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}
	return id, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

func decodeParseQuizTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r.URL.Query().Get(token), nil // TODO: Do I need to get from Header?
}

func decodeDeleteQuizByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "quiz_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed quiz_id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeStoreAnswerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return StoreAnswerRequest{
		UserID:     r.URL.Query().Get(userid),
		QuizID:     r.URL.Query().Get(quizid),
		QuestionID: r.URL.Query().Get(questionidi),
		AnswerID:   r.URL.Query().Get(answerid),
	}, nil
}
