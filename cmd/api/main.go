package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SatorNetwork/sator-api/internal/sumsub"

	db_internal "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/ethereum"
	"github.com/SatorNetwork/sator-api/internal/firebase"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/mail"
	"github.com/SatorNetwork/sator-api/internal/resizer"
	"github.com/SatorNetwork/sator-api/internal/solana"
	storage "github.com/SatorNetwork/sator-api/internal/storage"
	"github.com/SatorNetwork/sator-api/svc/auth"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/balance"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	challengeClient "github.com/SatorNetwork/sator-api/svc/challenge/client"
	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/SatorNetwork/sator-api/svc/files"
	filesRepo "github.com/SatorNetwork/sator-api/svc/files/repository"
	"github.com/SatorNetwork/sator-api/svc/invitations"
	invitationsClient "github.com/SatorNetwork/sator-api/svc/invitations/client"
	invitationsRepo "github.com/SatorNetwork/sator-api/svc/invitations/repository"
	"github.com/SatorNetwork/sator-api/svc/nft"
	nftRepo "github.com/SatorNetwork/sator-api/svc/nft/repository"
	"github.com/SatorNetwork/sator-api/svc/profile"
	profileRepo "github.com/SatorNetwork/sator-api/svc/profile/repository"
	"github.com/SatorNetwork/sator-api/svc/qrcodes"
	qrcodesRepo "github.com/SatorNetwork/sator-api/svc/qrcodes/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz"
	quizRepo "github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2"
	"github.com/SatorNetwork/sator-api/svc/referrals"
	referralsRepo "github.com/SatorNetwork/sator-api/svc/referrals/repository"
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

var buildTag string

// Application environment variables
var (
	// Build tag is set up while deployment
	buildTagDO = env.GetString("COMMIT_HASH", "")

	// General
	appPort            = env.MustInt("APP_PORT")
	appBaseURL         = env.MustString("APP_BASE_URL")
	httpRequestTimeout = env.GetDuration("HTTP_REQUEST_TIMEOUT", 30*time.Second)

	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 20)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 2)

	// JWT
	jwtSigningKey = env.MustString("JWT_SIGNING_KEY")
	jwtTTL        = env.GetDuration("JWT_TTL", 24*time.Hour)

	// Auth
	otpLength     = env.GetInt("OTP_LENGTH", 5)
	masterOTPHash = env.GetString("MASTER_OTP_HASH", "")

	// Quiz
	quizWsConnURL   = env.MustString("QUIZ_WS_CONN_URL")
	quizBotsTimeout = env.GetDuration("QUIZ_BOTS_TIMEOUT", 5*time.Second)

	// Solana
	solanaEnv                   = env.GetString("SOLANA_ENV", "devnet")
	solanaApiBaseUrl            = env.MustString("SOLANA_API_BASE_URL")
	solanaAssetAddr             = env.MustString("SOLANA_ASSET_ADDR")
	solanaFeePayerAddr          = env.MustString("SOLANA_FEE_PAYER_ADDR")
	solanaFeePayerPrivateKey    = env.MustString("SOLANA_FEE_PAYER_PRIVATE_KEY")
	solanaTokenHolderAddr       = env.MustString("SOLANA_TOKEN_HOLDER_ADDR")
	solanaTokenHolderPrivateKey = env.MustString("SOLANA_TOKEN_HOLDER_PRIVATE_KEY")
	tokenCirculatingSupply      = env.GetFloat("TOKEN_CIRCULATING_SUPPLY", 11839844)

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

	// Rewards
	holdRewardsPeriod = env.GetDuration("HOLD_REWARDS_PERIOD", 0)

	// Invitation
	invitationReward = env.GetFloat("INVITATION_REWARD", 0)
	invitationURL    = env.GetString("INVITATION_URL", "https://sator.io")

	// File Storage
	fileStorageKey            = env.MustString("STORAGE_KEY")
	fileStorageSecret         = env.MustString("STORAGE_SECRET")
	fileStorageEndpoint       = env.MustString("STORAGE_ENDPOINT")
	fileStorageRegion         = env.MustString("STORAGE_REGION")
	fileStorageBucket         = env.MustString("STORAGE_BUCKET")
	fileStorageUrl            = env.MustString("STORAGE_URL")
	fileStorageDisableSsl     = env.GetBool("STORAGE_DISABLE_SSL", true)
	fileStorageForcePathStyle = env.GetBool("STORAGE_FORCE_PATH_STYLE", false)

	// firebase
	baseFirebaseURL    = env.MustString("FIREBASE_BASE_URL")
	fbWebAPIKey        = env.MustString("FIREBASE_WEB_API_KEY")
	mainSiteLink       = env.MustString("FIREBASE_MAIN_SITE_LINK")
	androidPackageName = env.MustString("FIREBASE_ANDROID_PACKAGE_NAME")
	iosBundleId        = env.MustString("FIREBASE_IOS_BUNDLE_ID")
	suffixOption       = env.MustString("FIREBASE_SUFFIX_OPTION")

	// Min amounts
	minAmountToTransfer = env.GetFloat("MIN_AMOUNT_TO_TRANSFER", 0)
	minAmountToClaim    = env.GetFloat("MIN_AMOUNT_TO_CLAIM", 0)

	// KYC
	appToken  = env.MustString("KYC_APP_TOKEN")
	appSecret = env.MustString("KYC_APP_SECRET")
	baseURL   = env.MustString("KYC_APP_BASE_URL")
	ttl       = env.GetInt("KYC_APP_TTL", 1200)

	// NATS
	natsURL   = env.MustString("NATS_URL")
	natsWSURL = env.MustString("NATS_WS_URL")
)

var circulatingSupply float64 = 0

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

		r.Use(testingMdw)

		r.NotFound(notFoundHandler)
		r.MethodNotAllowed(methodNotAllowedHandler)

		r.Get("/", rootHandler)
		r.Get("/health", healthCheckHandler)
		r.Get("/supply", supplyHandler)
		r.Get("/ws", testWsHandler)
	}

	// auth repo
	authRepository, err := authRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("authRepo error: %v", err)
	}

	// Init JWT parser middleware
	// not depends on transport
	jwtMdw := jwt.NewParser(jwtSigningKey, jwt.CheckUser(authRepository.IsUserDisabled))
	jwtInteractor := jwt.NewInteractor(jwtSigningKey, jwtTTL)

	ethereumClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatalf("failed to init eth client: %v", err)
	}

	// KYC middleware
	kycMdw := sumsub.KYCStatusMdw(authRepository.GetKYCStatus)

	var walletSvcClient *walletClient.Client
	// Wallet service
	{
		walletRepository, err := walletRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("walletRepo error: %v", err)
		}

		feePayerPk, err := base64.StdEncoding.DecodeString(solanaFeePayerPrivateKey)
		if err != nil {
			log.Fatalf("feePayerPk base64 decoding error: %v", err)
		}
		tokenHolderPk, err := base64.StdEncoding.DecodeString(solanaTokenHolderPrivateKey)
		if err != nil {
			log.Fatalf("tokenHolderPk base64 decoding error: %v", err)
		}

		solanaClient := solana.New(solanaApiBaseUrl)
		if err := solanaClient.CheckPrivateKey(solanaFeePayerAddr, feePayerPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: fee payer: %v", err)
		}
		if err := solanaClient.CheckPrivateKey(solanaTokenHolderAddr, tokenHolderPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: token holder: %v", err)
		}

		walletService := wallet.NewService(
			walletRepository,
			solanaClient,
			ethereumClient,
			wallet.WithAssetSolanaAddress(solanaAssetAddr),
			wallet.WithSolanaFeePayer(solanaFeePayerAddr, feePayerPk),
			wallet.WithSolanaTokenHolder(solanaTokenHolderAddr, tokenHolderPk),
			wallet.WithMinAmountToTransfer(minAmountToTransfer),
		)
		walletSvcClient = walletClient.New(walletService)
		r.Mount("/wallets", wallet.MakeHTTPHandler(
			wallet.MakeEndpoints(walletService, kycMdw, jwtMdw),
			logger,
		))
	}

	// Rewards service
	var rewardsSvcClient *rewardsClient.Client

	rewardsRepository, err := rewardsRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("rewardsRepo error: %v", err)
	}
	rewardService := rewards.NewService(
		rewardsRepository,
		walletSvcClient,
		db_internal.NewAdvisoryLocks(db),
		rewards.WithExplorerURLTmpl("https://explorer.solana.com/tx/%s?cluster="+solanaEnv),
		rewards.WithHoldRewardsPeriod(holdRewardsPeriod),
		rewards.WithMinAmountToClaim(minAmountToClaim),
	)
	rewardsSvcClient = rewardsClient.New(rewardService)
	r.Mount("/rewards", rewards.MakeHTTPHandler(
		rewards.MakeEndpoints(rewardService, kycMdw, jwtMdw),
		logger,
	))

	// Invitation service
	invitationsRepository, err := invitationsRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("invitationsRepo error: %v", err)
	}
	invitationsService := invitations.NewService(invitationsRepository, mailer, rewardsSvcClient, invitations.Config{
		InvitationReward: invitationReward,
		InvitationURL:    invitationURL,
	})
	invitationsClient := invitationsClient.New(invitationsService)
	r.Mount("/invitations", invitations.MakeHTTPHandler(
		invitations.MakeEndpoints(invitationsService, jwtMdw),
		logger,
	))

	{
		// KYC
		kycService := sumsub.New(appToken, appSecret, baseURL, ttl)
		kycClient := sumsub.NewClient(kycService)

		// Auth service
		{
			r.Mount("/auth", auth.MakeHTTPHandler(
				auth.MakeEndpoints(auth.NewService(
					jwtInteractor,
					authRepository,
					walletSvcClient,
					invitationsClient,
					kycClient,
					auth.WithMasterOTPCode(masterOTPHash),
					auth.WithCustomOTPLength(otpLength),
					auth.WithMailService(mailer),
				), jwtMdw),
				logger,
			))
		}
	}

	// Profile service
	profileRepository, err := profileRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("profileRepo error: %v", err)
	}
	profileSvc := profile.NewService(profileRepository)
	r.Mount("/profile", profile.MakeHTTPHandler(
		profile.MakeEndpoints(profileSvc, jwtMdw),
		logger,
	))

	{
		// firebase connection
		var httpClient http.Client
		var fbClient firebase.FBClient
		fb := firebase.New(fbClient, httpClient, firebase.Config{
			BaseFirebaseURL:    baseFirebaseURL,
			WebAPIKey:          fbWebAPIKey,
			MainSiteLink:       mainSiteLink,
			AndroidPackageName: androidPackageName,
			IosBundleId:        iosBundleId,
			SuffixOption:       suffixOption,
		})

		// Referrals service
		referralRepository, err := referralsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("referralRepo error: %v", err)
		}
		r.Mount("/ref", referrals.MakeHTTPHandler(
			referrals.MakeEndpoints(referrals.NewService(referralRepository, fb, firebase.Config{
				BaseFirebaseURL:    baseFirebaseURL,
				WebAPIKey:          fbWebAPIKey,
				MainSiteLink:       mainSiteLink,
				AndroidPackageName: androidPackageName,
				IosBundleId:        iosBundleId,
				SuffixOption:       suffixOption,
			}), jwtMdw),
			logger,
		))
	}

	// Challenge client instance
	var challengeSvcClient *challengeClient.Client

	// Shows service
	{

		// Show repo
		// FIXME: remove it when the app will be fixed
		showRepo, err := showsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("showsRepo error: %v", err)
		}

		// Challenges service
		challengeRepository, err := challengeRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("challengeRepo error: %v", err)
		}
		challengeSvc := challenge.NewService(
			challengeRepository,
			showRepo,
			challenge.DefaultPlayURLGenerator(
				fmt.Sprintf("%s/challenges", strings.TrimSuffix(appBaseURL, "/")),
			),
			challenge.WithChargeForUnlockFunc(walletSvcClient.PayForService),
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
		r.Mount("/shows", shows.MakeHTTPHandler(
			shows.MakeEndpoints(shows.NewService(showRepo, challengeSvcClient), jwtMdw),
			logger,
		))
	}

	// files service
	{
		opt := storage.Options{
			Key:            fileStorageKey,
			Secret:         fileStorageSecret,
			Endpoint:       fileStorageEndpoint,
			Region:         fileStorageRegion,
			Bucket:         fileStorageBucket,
			URL:            fileStorageUrl,
			DisableSSL:     fileStorageDisableSsl,
			ForcePathStyle: fileStorageForcePathStyle,
		}
		stor := storage.New(storage.NewS3Client(opt), opt)

		mediaServiceRepo, err := filesRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("mediaServiceRepo error: %v", err)
		}
		r.Mount("/files", files.MakeHTTPHandler(
			files.MakeEndpoints(files.NewService(mediaServiceRepo, stor, resizer.Resize), jwtMdw),
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
			quiz.QuizWsHandler(
				quizSvc,
				invitationsService.SendReward(rewardService.AddTransaction),
				challengeSvcClient,
				profileSvc,
				quizBotsTimeout,
			),
		))

		// run quiz service
		g.Add(func() error {
			return quizSvc.Serve(ctx)
		}, func(err error) {
			log.Fatalf("quiz service: %v", err)
		})
	}

	{
		quizV2Svc := quiz_v2.NewService(natsURL, natsWSURL, challengeSvcClient)
		r.Mount("/quiz_v2", quiz_v2.MakeHTTPHandler(
			quiz_v2.MakeEndpoints(quizV2Svc, jwtMdw),
			logger,
		))

		go quizV2Svc.StartEngine()
		// TODO(evg): gracefully shutdown the engine
	}

	{
		// NFT service
		nftRepository, err := nftRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("nftRepo error: %v", err)
		}
		nftService := nft.NewService(nftRepository, walletSvcClient.PayForService)
		r.Mount("/nft", nft.MakeHTTPHandler(
			nft.MakeEndpoints(nftService, jwtMdw),
			logger,
		))
	}

	{
		// Init and run http server
		httpServer := &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", appPort),
		}
		g.Add(func() error {
			log.Printf("[http-server] start listening on :%d...\n", appPort)
			if err := httpServer.ListenAndServe(); err != nil {
				fmt.Println("[http-server] stopped listening with error:", err)
				return err
			}
			return nil
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

		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		tickerDone := make(chan bool)
		defer close(tickerDone)

		g.Add(func() error {
			circulatingSupply = tokenCirculatingSupply
			for {
				select {
				case <-tickerDone:
					return nil
				case <-ticker.C:
					circulatingSupply++
				}
			}
		}, func(err error) {
			tickerDone <- true
		})
	}

	if err := g.Run(); err != nil {
		log.Println("API terminated with error:", err)
	}
}

// returns current build tag
func rootHandler(w http.ResponseWriter, _ *http.Request) {
	if buildTag == "" {
		buildTag = buildTagDO
	}
	defaultResponse(w, http.StatusOK, map[string]interface{}{"build_tag": buildTag})
}

// returns token circulating supply
func supplyHandler(w http.ResponseWriter, _ *http.Request) {
	defaultResponse(w, http.StatusOK, map[string]interface{}{
		"supply": circulatingSupply,
	})
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

// Testing middleware
// Helps to test any HTTP error
// Pass must_err query parameter with code you want get
// E.g.: /shows?must_err=403
func testingMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if errCodeStr := r.URL.Query().Get("must_err"); len(errCodeStr) == 3 {
			if errCode, err := strconv.Atoi(errCodeStr); err == nil && errCode >= 400 && errCode < 600 {
				w.WriteHeader(errCode)
				w.Write([]byte(http.StatusText(errCode)))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
