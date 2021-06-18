package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	"github.com/SatorNetwork/sator-api/svc/questions"
	questionsClient "github.com/SatorNetwork/sator-api/svc/questions/client"
	questionsRepo "github.com/SatorNetwork/sator-api/svc/questions/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz"
	quizRepo "github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/SatorNetwork/sator-api/svc/rewards"
	rewardsClient "github.com/SatorNetwork/sator-api/svc/rewards/client"
	rewardsRepo "github.com/SatorNetwork/sator-api/svc/rewards/repository"
	"github.com/SatorNetwork/sator-api/svc/shows"
	showsRepo "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/dmitrymomot/distlock"
	"github.com/dmitrymomot/distlock/inmem"
	"github.com/dmitrymomot/go-env"
	signature "github.com/dmitrymomot/go-signature"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kitlog "github.com/go-kit/kit/log"
	_ "github.com/lib/pq" // init pg driver
	"github.com/oklog/run"
	"github.com/rs/cors"
	"github.com/zeebo/errs"
)

// Build tag is set up while compiling
var buildTag string

// Application environment variables
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

	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Llongfile)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mutex := distlock.New(
		distlock.WithStorageDrivers(inmem.New()),
		distlock.WithTries(10),
	)

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

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

	// rewards repo
	// TODO: needs refactoring
	rewardRepo, err := rewardsRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("rewardsRepo error: %v", err)
	}

	// Wallet service
	wRepo, err := walletRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("walletRepo error: %v", err)
	}
	walletService := wallet.NewService(wRepo, solana.New(solanaApiBaseUrl), rewardRepo)
	r.Mount("/wallet", wallet.MakeHTTPHandler(
		wallet.MakeEndpoints(walletService, jwtMdw),
		logger,
	))

	// Auth service
	{
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
				os.Getenv("MASTER_OTP_HASH"),
				auth.WithCustomOTPLength(otpLength),
				// auth.WithMailService(/** encapsulate mail service */),
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
	challengeRepo, err := challengeRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("challengeRepo error: %v", err)
	}
	challengeSvc := challenge.NewService(
		challengeRepo,
		challenge.DefaultPlayURLGenerator(
			fmt.Sprintf("%s/challenges", strings.TrimSuffix(appBaseURL, "/")),
		),
	)
	challengeClient := challengeClient.New(challengeSvc)
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
		shows.MakeEndpoints(shows.NewService(showRepo, challengeClient), jwtMdw),
		logger,
	))

	// Questions service
	questRepo, err := questionsRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("questionsRepo error: %v", err)
	}
	questClient := questionsClient.New(questions.NewService(questRepo))

	// Rewards service
	rewardSvc := rewards.NewService(rewardRepo, walletService)
	rewardClient := rewardsClient.New(rewardSvc)
	r.Mount("/rewards", rewards.MakeHTTPHandler(
		rewards.MakeEndpoints(rewardSvc, jwtMdw),
		logger,
	))

	// Quiz service
	quizRepo, err := quizRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("quizRepo error: %v", err)
	}
	quizSvc := quiz.NewService(
		mutex,
		quizRepo,
		questClient,
		rewardClient,
		challengeClient,
		quizWsConnURL,
		quiz.WithCustomTokenGenerateFunction(signature.NewTemporary),
		quiz.WithCustomTokenParseFunction(signature.Parse),
	)
	r.Mount("/quiz", quiz.MakeHTTPHandler(
		quiz.MakeEndpoints(quizSvc, jwtMdw),
		logger,
		quiz.QuizWsHandler(quizSvc),
	))

	var g run.Group
	{
		// run quiz service
		g.Add(func() error {
			return quizSvc.Serve(ctx)
		}, func(err error) {
			log.Fatalf("quiz service: %v", err)
		})

		// Init and run http server
		httpServer := &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", appPort),
		}
		g.Add(func() error {
			log.Printf("[http-server] start listening on :%d...\n", appPort)
			err := httpServer.ListenAndServe()
			if err != nil {
				fmt.Println("[http-server] stopped listening with error:", err)
			}
			return err
		}, func(err error) {
			fmt.Println("[http-server] terminating because of error:", err)
			_ = httpServer.Shutdown(context.Background())
		})

		g.Add(func() error {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			c := <-sigChan
			return fmt.Errorf("terminated with sig %q", c)
		}, func(err error) {})
	}

	// Init and run http server
	// httpServer := &http.Server{
	// 	Handler: r,
	// 	Addr:    fmt.Sprintf(":%d", appPort),
	// }
	// httpServer.RegisterOnShutdown(cancel)
	// graceful.LogListenAndServe(httpServer, log.Default())

	if err := g.Run(); err != nil {
		log.Println("API terminated with error:", err)
	}
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
