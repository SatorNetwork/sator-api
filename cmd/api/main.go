package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/auth"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	challengeClient "github.com/SatorNetwork/sator-api/svc/challenge/client"
	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/SatorNetwork/sator-api/svc/profile"
	profileRepo "github.com/SatorNetwork/sator-api/svc/profile/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz"
	quizRepo "github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/SatorNetwork/sator-api/svc/shows"
	showsRepo "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
	signature "github.com/dmitrymomot/go-signature"

	"github.com/TV4/graceful"
	"github.com/dmitrymomot/go-env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kitlog "github.com/go-kit/kit/log"
	_ "github.com/lib/pq" // init pg driver
	"github.com/rs/cors"
)

// Build tag is set up while compiling
var buildTag string

// Application environemnt variables
var (
	// General
	appPort            = env.MustInt("APP_PORT")
	appBaseURL         = env.MustString("APP_BASE_URL")
	httpRequestTimeout = env.GetDuration("HTTP_REQUEST_TIMEOUT", 30*time.Second)

	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 10)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// JWT
	jwtSigningKey = env.MustString("JWT_SIGNING_KEY")
	jwtTTL        = env.GetDuration("JWT_TTL", 24*time.Hour)

	// Auth
	otpLength = env.GetInt("OTP_LENGTH", 5)

	// Quiz
	quizWsConnURL = env.MustString("QUIZ_WS_CONN_URL")
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Llongfile)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	if err := db.Ping(); err != nil {
		log.Fatalf("db pinng error: %v", err)
	}

	r := chi.NewRouter()
	{
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(httpRequestTimeout))
		r.Use(cors.AllowAll().Handler)

		r.NotFound(notFoundHandler)
		r.MethodNotAllowed(methodNotAllowedHandler)

		r.Get("/", rootHandler)
		r.Get("/health", healthCheckHandler)
		r.Get("/ws", testWsHandler)
	}

	// Init JWT parser middleware
	// not depends on transport
	jwtMdw := jwt.NewParser(jwtSigningKey)
	jwtInteractor := jwt.NewInteractor(jwtSigningKey, jwtTTL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	solanaClient := solana.New()

	// Auth service
	{
		// Wallet service
		wRepo, err := walletRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("walletRepo error: %v", err)
		}
		walletService := wallet.NewService(wRepo, solanaClient)
		r.Mount("/wallet", wallet.MakeHTTPHandler(
			wallet.MakeEndpoints(walletService, jwtMdw),
			logger,
		))

		// auth
		repo, err := authRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("authRepo error: %v", err)
		}
		r.Mount("/auth", auth.MakeHTTPHandler(
			auth.MakeEndpoints(auth.NewService(
				jwtInteractor,
				repo,
				walletService,
				auth.WithCustomOTPLength(otpLength),
				// auth.WithMailService(/** incapsulate mail service */),
			), jwtMdw),
			logger,
		))
	}

	// Profile service
	{
		repo, err := profileRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("profileRepo error: %v", err)
		}
		r.Mount("/profile", profile.MakeHTTPHandler(
			profile.MakeEndpoints(profile.NewService(repo), jwtMdw),
			logger,
		))
	}

	// Challenges service
	{
		repo, err := challengeRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("challengeRepo error: %v", err)
		}
		challengeSvc := challenge.NewService(
			repo,
			challenge.DefaultPlayURLGenerator(
				fmt.Sprintf("%s/challenges", strings.TrimSuffix(appBaseURL, "/")),
			),
		)
		r.Mount("/challenges", challenge.MakeHTTPHandler(
			challenge.MakeEndpoints(challengeSvc, jwtMdw),
			logger,
		))

		// Shows service
		showRepo, err := showsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("showsRepo error: %v", err)
		}
		r.Mount("/shows", shows.MakeHTTPHandler(
			shows.MakeEndpoints(shows.NewService(showRepo, challengeClient.New(challengeSvc)), jwtMdw),
			logger,
		))
	}

	// Quiz service
	{
		repo, err := quizRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("quizRepo error: %v", err)
		}
		quizSvc := quiz.NewService(
			repo,
			quizWsConnURL,
			quiz.WithCustomTokenGenerateFunction(signature.NewTemporary),
			quiz.WithCustomTokenParseFunction(signature.Parse),
		)
		r.Mount("/quiz", quiz.MakeHTTPHandler(
			quiz.MakeEndpoints(quizSvc, jwtMdw),
			logger,
			quiz.QuizWsHandler(quizSvc.ParseQuizToken),
		))
	}

	// Init and run http server
	httpServer := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", appPort),
	}
	httpServer.RegisterOnShutdown(cancel)
	graceful.LogListenAndServe(httpServer, log.Default())
}

// returns current build tag
func rootHandler(w http.ResponseWriter, _ *http.Request) {
	defaultResponse(w, http.StatusOK, map[string]interface{}{"build_tag": buildTag})
}

// returns html page to test websocket
func testWsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./cmd/api/index.html")
}

// returns 204 HTTP status without content
func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// returns 404 HTTP status with payload
func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	defaultResponse(w, http.StatusNotFound, map[string]interface{}{
		"error": http.StatusText(http.StatusNotFound),
	})
}

// returns 405 HTTP status with payload
func methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	defaultResponse(w, http.StatusMethodNotAllowed, map[string]interface{}{
		"error": http.StatusText(http.StatusMethodNotAllowed),
	})
}

// helper to send response as a json data
func defaultResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
