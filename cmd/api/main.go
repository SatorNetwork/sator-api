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
	"github.com/SatorNetwork/sator-api/internal/mail"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/auth"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/balance"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	challengeClient "github.com/SatorNetwork/sator-api/svc/challenge/client"
	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/SatorNetwork/sator-api/svc/profile"
	profileRepo "github.com/SatorNetwork/sator-api/svc/profile/repository"
	"github.com/SatorNetwork/sator-api/svc/qrcodes"
	qrcodesRepo "github.com/SatorNetwork/sator-api/svc/qrcodes/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz"
	quizRepo "github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/SatorNetwork/sator-api/svc/rewards"
	rewardsClient "github.com/SatorNetwork/sator-api/svc/rewards/client"
	rewardsRepo "github.com/SatorNetwork/sator-api/svc/rewards/repository"
	"github.com/SatorNetwork/sator-api/svc/shows"
	showsRepo "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	walletClient "github.com/SatorNetwork/sator-api/svc/wallet/client"
	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/dmitrymomot/distlock"
	"github.com/dmitrymomot/distlock/inmem"
	"github.com/dmitrymomot/go-env"
	signature "github.com/dmitrymomot/go-signature"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kitlog "github.com/go-kit/kit/log"
	"github.com/keighl/postmark"
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
	otpLength     = env.GetInt("OTP_LENGTH", 5)
	masterOTPHash = env.GetString("MASTER_OTP_HASH", "")

	// Quiz
	quizWsConnURL = env.MustString("QUIZ_WS_CONN_URL")

	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")

	// Mailer
	postmarkServerToken   = env.MustString("POSTMARK_SERVER_TOKEN")
	postmarkAccountToken  = env.MustString("POSTMARK_ACCOUNT_TOKEN")
	notificationFromName  = env.GetString("NOTIFICATION_FROM_NAME", "Sator.io")
	notificationFromEmail = env.GetString("NOTIFICATION_FROM_EMAIL", "notifications@sator.io")

	// Product
	productName    = env.GetString("PRODUCT_NAME", "Sator.io")
	productURL     = env.GetString("PRODUCT_URL", "https://sator.io")
	supportURL     = env.GetString("SUPPORT_URL", "https://sator.io")
	supportEmail   = env.GetString("SUPPORT_EMAIL", "support@sator.io")
	companyName    = env.GetString("COMPANY_NAME", "Sator")
	companyAddress = env.GetString("COMPANY_ADDRESS", "New York")
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

	// runtime group
	var g run.Group

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

	// Init mail service
	mailer := mail.New(postmark.NewClient(postmarkServerToken, postmarkAccountToken), mail.Config{
		ProductName:    productName,
		ProductURL:     productURL,
		SupportURL:     supportURL,
		SupportEmail:   supportEmail,
		CompanyName:    companyName,
		CompanyAddress: companyAddress,
		FromEmail:      notificationFromEmail,
		FromName:       notificationFromName,
	})

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

	var walletSvcClient *walletClient.Client
	// Wallet service
	{
		walletRepository, err := walletRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("walletRepo error: %v", err)
		}
		walletService := wallet.NewService(walletRepository, solana.New(solanaApiBaseUrl))
		walletSvcClient = walletClient.New(walletService)
		r.Mount("/wallets", wallet.MakeHTTPHandler(
			wallet.MakeEndpoints(walletService, jwtMdw),
			logger,
		))
	}

	// Auth service
	{
		// auth
		authRepository, err := authRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("authRepo error: %v", err)
		}
		r.Mount("/auth", auth.MakeHTTPHandler(
			auth.MakeEndpoints(auth.NewService(
				jwtInteractor,
				authRepository,
				walletSvcClient,
				auth.WithMasterOTPCode(masterOTPHash),
				auth.WithCustomOTPLength(otpLength),
				auth.WithMailService(mailer),
			), jwtMdw),
			logger,
		))
	}

	// Profile service
	{
		profileRepository, err := profileRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("profileRepo error: %v", err)
		}
		r.Mount("/profile", profile.MakeHTTPHandler(
			profile.MakeEndpoints(profile.NewService(profileRepository), jwtMdw),
			logger,
		))
	}

	// Challenge client instance
	var challengeSvcClient *challengeClient.Client

	// Shows service
	{
		// Challenges service
		challengeRepository, err := challengeRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("challengeRepo error: %v", err)
		}
		challengeSvc := challenge.NewService(
			challengeRepository,
			challenge.DefaultPlayURLGenerator(
				fmt.Sprintf("%s/challenges", strings.TrimSuffix(appBaseURL, "/")),
			),
		)
		challengeSvcClient = challengeClient.New(challengeSvc)
		r.Mount("/challenges", challenge.MakeHTTPHandlerChallenges(
			challenge.MakeEndpoints(challengeSvc, jwtMdw),
			logger,
		))

		r.Mount("/questions", challenge.MakeHTTPHandlerQuestions(
			challenge.MakeEndpoints(challengeSvc, jwtMdw),
			logger,
		))

		// Shows service
		showRepo, err := showsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("showsRepo error: %v", err)
		}
		r.Mount("/shows", shows.MakeHTTPHandler(
			shows.MakeEndpoints(shows.NewService(showRepo, challengeSvcClient), jwtMdw),
			logger,
		))
	}

	var rewardsSvcClient *rewardsClient.Client
	// Rewards service
	{
		rewardsRepository, err := rewardsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("rewardsRepo error: %v", err)
		}
		rewardService := rewards.NewService(rewardsRepository, walletSvcClient)
		rewardsSvcClient = rewardsClient.New(rewardService)
		r.Mount("/rewards", rewards.MakeHTTPHandler(
			rewards.MakeEndpoints(rewardService, jwtMdw),
			logger,
		))
	}

	// Balance service
	{
		r.Mount("/balance", balance.MakeHTTPHandler(
			balance.MakeEndpoints(balance.NewService(walletSvcClient, rewardsSvcClient), jwtMdw),
			logger,
		))
	}

	// QR-codes service
	{
		qrcodesRepository, err := qrcodesRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("qrcodesRepo error: %v", err)
		}
		r.Mount("/qrcodes", qrcodes.MakeHTTPHandler(
			qrcodes.MakeEndpoints(qrcodes.NewService(qrcodesRepository, rewardsSvcClient), jwtMdw),
			logger,
		))
	}

	// Quiz service
	{
		// Quiz service
		quizRepository, err := quizRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("quizRepo error: %v", err)
		}
		quizSvc := quiz.NewService(
			mutex,
			quizRepository,
			rewardsSvcClient,
			challengeSvcClient,
			quizWsConnURL,
			quiz.WithCustomTokenGenerateFunction(signature.NewTemporary),
			quiz.WithCustomTokenParseFunction(signature.Parse),
		)
		r.Mount("/quiz", quiz.MakeHTTPHandler(
			quiz.MakeEndpoints(quizSvc, jwtMdw),
			logger,
			quiz.QuizWsHandler(quizSvc),
		))

		// run quiz service
		g.Add(func() error {
			return quizSvc.Serve(ctx)
		}, func(err error) {
			log.Fatalf("quiz service: %v", err)
		})
	}

	{
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
