package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/svc/auth"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/profile"
	profileRepo "github.com/SatorNetwork/sator-api/svc/profile/repository"
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
	}

	// Init JWT parser middleware
	// not depends on transport
	jwtMdw := jwt.NewParser(jwtSigningKey)
	jwtInteractor := jwt.NewInteractor(jwtSigningKey, jwtTTL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Auth service
	{
		repo, err := authRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("authRepo error: %v", err)
		}
		r.Mount("/auth", auth.MakeHTTPHandler(
			auth.MakeEndpoints(auth.NewService(
				jwtInteractor,
				repo,
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
