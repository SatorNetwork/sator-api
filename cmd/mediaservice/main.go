package mediaservice

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SatorNetwork/sator-api/internal/mediaservice/handler"
	"github.com/SatorNetwork/sator-api/internal/mediaservice/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	_ "github.com/lib/pq" // init pg driver

	"github.com/TV4/graceful"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

// Build tag is set up while compiling
var buildTag string

func init() {
	// load environment config
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
}

func main() {
	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)
	db, err := sql.Open(os.Getenv("DB_DRIVER"), dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(3)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := mediaservice.Migrate(db, os.Getenv("DB_DRIVER"), "/sql/migrations"); err != nil {
		log.Fatal(err)
	}

	query := repository.New(db)

	disableSSL, _ := strconv.ParseBool(os.Getenv("STORAGE_DISABLE_SSL"))
	forcePathStyle, _ := strconv.ParseBool(os.Getenv("STORAGE_FORCE_PATH_STYLE"))
	opt := storage.Options{
		Key:            os.Getenv("STORAGE_KEY"),
		Secret:         os.Getenv("STORAGE_SECRET"),
		Endpoint:       os.Getenv("STORAGE_ENDPOINT"),
		Region:         os.Getenv("STORAGE_REGION"),
		Bucket:         os.Getenv("STORAGE_BUCKET"),
		URL:            os.Getenv("STORAGE_URL"),
		DisableSSL:     disableSSL,
		ForcePathStyle: forcePathStyle,
	}
	stor := storage.New(storage.NewS3Client(opt), opt)

	h := handler.New(db, query, stor)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("build_tag: %s", buildTag)))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Mount(fmt.Sprintf("/%s", os.Getenv("API_VERSION")), handler.Router(h))

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", os.Getenv("API_PORT")),
	}
	graceful.LogListenAndServe(s)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": http.StatusText(http.StatusNotFound)})
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": http.StatusText(http.StatusMethodNotAllowed)})
}
