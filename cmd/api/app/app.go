package app

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

	"github.com/dmitrymomot/go-env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kitlog "github.com/go-kit/kit/log"
	"github.com/keighl/postmark"
	_ "github.com/lib/pq" // init pg driver
	"github.com/oklog/run"
	"github.com/rs/cors"
	"github.com/zeebo/errs"

	db_internal "github.com/SatorNetwork/sator-api/lib/db"
	internal_rsa "github.com/SatorNetwork/sator-api/lib/encryption/rsa"
	"github.com/SatorNetwork/sator-api/lib/ethereum"
	"github.com/SatorNetwork/sator-api/lib/firebase"
	"github.com/SatorNetwork/sator-api/lib/jwt"
	lib_postmark "github.com/SatorNetwork/sator-api/lib/mail/postmark"
	"github.com/SatorNetwork/sator-api/lib/resizer"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	storage "github.com/SatorNetwork/sator-api/lib/storage"
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	"github.com/SatorNetwork/sator-api/svc/auth"
	authc "github.com/SatorNetwork/sator-api/svc/auth/client"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/balance"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	challengeClient "github.com/SatorNetwork/sator-api/svc/challenge/client"
	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	exchange_rates_repository "github.com/SatorNetwork/sator-api/svc/exchange_rates/repository"
	"github.com/SatorNetwork/sator-api/svc/files"
	filesRepo "github.com/SatorNetwork/sator-api/svc/files/repository"
	"github.com/SatorNetwork/sator-api/svc/invitations"
	invitationsClient "github.com/SatorNetwork/sator-api/svc/invitations/client"
	invitationsRepo "github.com/SatorNetwork/sator-api/svc/invitations/repository"
	"github.com/SatorNetwork/sator-api/svc/nft"
	nftC "github.com/SatorNetwork/sator-api/svc/nft/client"
	nftRepo "github.com/SatorNetwork/sator-api/svc/nft/repository"
	"github.com/SatorNetwork/sator-api/svc/profile"
	profileRepo "github.com/SatorNetwork/sator-api/svc/profile/repository"
	"github.com/SatorNetwork/sator-api/svc/puzzle_game"
	puzzleGameRepo "github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
	"github.com/SatorNetwork/sator-api/svc/qrcodes"
	qrcodesRepo "github.com/SatorNetwork/sator-api/svc/qrcodes/repository"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2"
	quizV2Repo "github.com/SatorNetwork/sator-api/svc/quiz_v2/repository"
	"github.com/SatorNetwork/sator-api/svc/referrals"
	referralsRepo "github.com/SatorNetwork/sator-api/svc/referrals/repository"
	"github.com/SatorNetwork/sator-api/svc/rewards"
	rewardsClient "github.com/SatorNetwork/sator-api/svc/rewards/client"
	rewardsRepo "github.com/SatorNetwork/sator-api/svc/rewards/repository"
	"github.com/SatorNetwork/sator-api/svc/shows"
	"github.com/SatorNetwork/sator-api/svc/shows/private"
	showsRepo "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/svc/trading_platforms"
	tradingPlatformsRepo "github.com/SatorNetwork/sator-api/svc/trading_platforms/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	walletClient "github.com/SatorNetwork/sator-api/svc/wallet/client"
	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

type Config struct {
	BuildTagDO                  string
	AppPort                     int
	AppBaseURL                  string
	HttpRequestTimeout          time.Duration
	DBConnString                string
	DBMaxOpenConns              int
	DBMaxIdleConns              int
	JwtSigningKey               string
	JwtTTL                      time.Duration
	OtpLength                   int
	MasterOTPHash               string
	QuizWsConnURL               string
	QuizBotsTimeout             time.Duration
	QuizLobbyLatency            time.Duration
	TokenCirculatingSupply      float64
	SolanaEnv                   string
	SolanaApiBaseUrl            string
	SolanaAssetAddr             string
	SolanaFeePayerAddr          string
	SolanaFeePayerPrivateKey    string
	SolanaTokenHolderAddr       string
	SolanaTokenHolderPrivateKey string
	SolanaStakePoolAddr         string
	SolanaSystemProgram         string
	SolanaSysvarRent            string
	SolanaSysvarClock           string
	SolanaSplToken              string
	SolanaStakeProgramID        string
	PostmarkServerToken         string
	PostmarkAccountToken        string
	NotificationFromName        string
	NotificationFromEmail       string
	ProductName                 string
	ProductURL                  string
	SupportURL                  string
	SupportEmail                string
	CompanyName                 string
	CompanyAddress              string
	HoldRewardsPeriod           time.Duration
	InvitationReward            float64
	InvitationURL               string
	FileStorageKey              string
	FileStorageSecret           string
	FileStorageEndpoint         string
	FileStorageRegion           string
	FileStorageBucket           string
	FileStorageUrl              string
	FileStorageDisableSsl       bool
	FileStorageForcePathStyle   bool
	BaseFirebaseURL             string
	FBWebAPIKey                 string
	MainSiteLink                string
	AndroidPackageName          string
	IOSBundleId                 string
	SuffixOption                string
	MinAmountToTransfer         float64
	MinAmountToClaim            float64
	KycAppToken                 string
	KycAppSecret                string
	KycAppBaseURL               string
	KycAppTTL                   int
	KycSkip                     bool
	NatsURL                     string
	NatsWSURL                   string
	QuizV2ShuffleQuestions      bool
	ServerRSAPrivateKey         string
	TipsPercent                 float64
	SatorAPIKey                 string
	WhitelistMode               bool
	BlacklistMode               bool
}

var buildTag string

// Application environment variables
func ConfigFromEnv() *Config {
	return &Config{
		// Build tag is set up while deployment
		BuildTagDO: env.GetString("COMMIT_HASH", ""),

		// General
		AppPort:            env.MustInt("APP_PORT"),
		AppBaseURL:         env.MustString("APP_BASE_URL"),
		HttpRequestTimeout: env.GetDuration("HTTP_REQUEST_TIMEOUT", 30*time.Second),

		// DB
		DBConnString:   env.MustString("DATABASE_URL"),
		DBMaxOpenConns: env.GetInt("DATABASE_MAX_OPEN_CONNS", 20),
		DBMaxIdleConns: env.GetInt("DATABASE_IDLE_CONNS", 2),

		// Auth
		OtpLength:     env.GetInt("OTP_LENGTH", 5),
		MasterOTPHash: env.GetString("MASTER_OTP_HASH", ""),
		WhitelistMode: env.GetBool("WHITELIST_MODE", true),
		BlacklistMode: env.GetBool("BLACKLIST_MODE", true),

		// JWT
		JwtSigningKey: env.MustString("JWT_SIGNING_KEY"),
		JwtTTL:        env.GetDuration("JWT_TTL", 24*time.Hour),

		// Quiz
		QuizWsConnURL:    env.MustString("QUIZ_WS_CONN_URL"),
		QuizBotsTimeout:  env.GetDuration("QUIZ_BOTS_TIMEOUT", 5*time.Second),
		QuizLobbyLatency: env.GetDuration("QUIZ_LOBBY_LATENCY", 5*time.Second),

		// Solana
		TokenCirculatingSupply:      env.GetFloat("TOKEN_CIRCULATING_SUPPLY", 11839844),
		SolanaEnv:                   env.GetString("SOLANA_ENV", "devnet"),
		SolanaApiBaseUrl:            env.MustString("SOLANA_API_BASE_URL"),
		SolanaAssetAddr:             env.MustString("SOLANA_ASSET_ADDR"),
		SolanaFeePayerAddr:          env.MustString("SOLANA_FEE_PAYER_ADDR"),
		SolanaFeePayerPrivateKey:    env.MustString("SOLANA_FEE_PAYER_PRIVATE_KEY"),
		SolanaTokenHolderAddr:       env.MustString("SOLANA_TOKEN_HOLDER_ADDR"),
		SolanaTokenHolderPrivateKey: env.MustString("SOLANA_TOKEN_HOLDER_PRIVATE_KEY"),

		// Tokens lock pool
		SolanaStakePoolAddr:  env.MustString("SOLANA_STAKE_POOL_ADDR"),
		SolanaSystemProgram:  env.MustString("SOLANA_SYSTEM_PROGRAM"),
		SolanaSysvarRent:     env.MustString("SOLANA_SYSVAR_RENT"),
		SolanaSysvarClock:    env.MustString("SOLANA_SYSVAR_CLOCK"),
		SolanaSplToken:       env.MustString("SOLANA_SPL_TOKEN"),
		SolanaStakeProgramID: env.MustString("SOLANA_STAKE_PROGRAM_ID"),

		// Mailer
		PostmarkServerToken:   env.MustString("POSTMARK_SERVER_TOKEN"),
		PostmarkAccountToken:  env.MustString("POSTMARK_ACCOUNT_TOKEN"),
		NotificationFromName:  env.GetString("NOTIFICATION_FROM_NAME", "Sator.io"),
		NotificationFromEmail: env.GetString("NOTIFICATION_FROM_EMAIL", "notifications@sator.io"),

		// Product
		ProductName:    env.GetString("PRODUCT_NAME", "Sator.io"),
		ProductURL:     env.GetString("PRODUCT_URL", "https://sator.io"),
		SupportURL:     env.GetString("SUPPORT_URL", "https://sator.io"),
		SupportEmail:   env.GetString("SUPPORT_EMAIL", "support@sator.io"),
		CompanyName:    env.GetString("COMPANY_NAME", "Sator"),
		CompanyAddress: env.GetString("COMPANY_ADDRESS", "New York"),

		// Rewards
		HoldRewardsPeriod: env.GetDuration("HOLD_REWARDS_PERIOD", 0),

		// Invitation
		InvitationReward: env.GetFloat("INVITATION_REWARD", 0),
		InvitationURL:    env.GetString("INVITATION_URL", "https://sator.io"),

		// File Storage
		FileStorageKey:            env.MustString("STORAGE_KEY"),
		FileStorageSecret:         env.MustString("STORAGE_SECRET"),
		FileStorageEndpoint:       env.MustString("STORAGE_ENDPOINT"),
		FileStorageRegion:         env.MustString("STORAGE_REGION"),
		FileStorageBucket:         env.MustString("STORAGE_BUCKET"),
		FileStorageUrl:            env.MustString("STORAGE_URL"),
		FileStorageDisableSsl:     env.GetBool("STORAGE_DISABLE_SSL", true),
		FileStorageForcePathStyle: env.GetBool("STORAGE_FORCE_PATH_STYLE", false),

		// firebase
		BaseFirebaseURL:    env.MustString("FIREBASE_BASE_URL"),
		FBWebAPIKey:        env.MustString("FIREBASE_WEB_API_KEY"),
		MainSiteLink:       env.MustString("FIREBASE_MAIN_SITE_LINK"),
		AndroidPackageName: env.MustString("FIREBASE_ANDROID_PACKAGE_NAME"),
		IOSBundleId:        env.MustString("FIREBASE_IOS_BUNDLE_ID"),
		SuffixOption:       env.MustString("FIREBASE_SUFFIX_OPTION"),

		// Min amounts
		MinAmountToTransfer: env.GetFloat("MIN_AMOUNT_TO_TRANSFER", 0),
		MinAmountToClaim:    env.GetFloat("MIN_AMOUNT_TO_CLAIM", 0),

		// KYC
		KycAppToken:   env.MustString("KYC_APP_TOKEN"),
		KycAppSecret:  env.MustString("KYC_APP_SECRET"),
		KycAppBaseURL: env.MustString("KYC_APP_BASE_URL"),
		KycAppTTL:     env.GetInt("KYC_APP_TTL", 1200),
		KycSkip:       env.GetBool("KYC_SKIP", false),

		// NATS
		NatsURL:   env.MustString("NATS_URL"),
		NatsWSURL: env.MustString("NATS_WS_URL"),

		QuizV2ShuffleQuestions: env.GetBool("QUIZ_V2_SHUFFLE_QUESTIONS", true),
		ServerRSAPrivateKey:    env.MustString("SERVER_RSA_PRIVATE_KEY"),

		TipsPercent: env.GetFloat("TIPS_PERCENT", 0.5),

		SatorAPIKey: env.MustString("SATOR_API_KEY"),
	}
}

var circulatingSupply float64 = 0

type app struct {
	cfg *Config

	shutdown bool
	done     chan struct{}
}

func New() (*app, error) {
	cfg := ConfigFromEnv()
	return WithConfig(cfg), nil
}

func WithConfig(cfg *Config) *app {
	return &app{
		cfg:  cfg,
		done: make(chan struct{}),
	}
}

func (a *app) Run() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Llongfile)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// runtime group
	var g run.Group

	// Init DB connection
	db, err := sql.Open("postgres", a.cfg.DBConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	db.SetMaxOpenConns(a.cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(a.cfg.DBMaxIdleConns)

	if err := db.Ping(); err != nil {
		log.Fatalf("db pinng error: %v", err)
	}

	// Init mail service
	mailer := lib_postmark.New(postmark.NewClient(a.cfg.PostmarkServerToken, a.cfg.PostmarkAccountToken), lib_postmark.Config{
		ProductName:    a.cfg.ProductName,
		ProductURL:     a.cfg.ProductURL,
		SupportURL:     a.cfg.SupportURL,
		SupportEmail:   a.cfg.SupportEmail,
		CompanyName:    a.cfg.CompanyName,
		CompanyAddress: a.cfg.CompanyAddress,
		FromEmail:      a.cfg.NotificationFromEmail,
		FromName:       a.cfg.NotificationFromName,
	})

	r := chi.NewRouter()
	{
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(a.cfg.HttpRequestTimeout))
		r.Use(cors.AllowAll().Handler)

		r.Use(testingMdw)

		r.NotFound(notFoundHandler)
		r.MethodNotAllowed(methodNotAllowedHandler)

		r.Get("/", mkRootHandler(a.cfg.BuildTagDO))
		r.Get("/health", healthCheckHandler)
		r.Get("/supply", supplyHandler)
		// r.Get("/ws", testWsHandler)
	}

	serverRSAPrivateKey, err := internal_rsa.BytesToPrivateKey([]byte(a.cfg.ServerRSAPrivateKey))
	if err != nil {
		log.Fatalf("can't decode server's RSA private key")
	}

	// auth repo
	authRepository, err := authRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("authRepo error: %v", err)
	}

	// Init JWT parser middleware
	// not depends on transport
	jwtMdw := jwt.NewParser(a.cfg.JwtSigningKey, jwt.CheckUser(authRepository.IsUserDisabled), authRepository)
	jwtInteractor := jwt.NewInteractor(a.cfg.JwtSigningKey, a.cfg.JwtTTL)

	ethereumClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatalf("failed to init eth client: %v", err)
	}

	// KYC middleware
	kycMdw := sumsub.KYCStatusMdw(authRepository.GetKYCStatus, func() bool {
		return a.cfg.KycSkip
	})

	var exchangeRatesClient *exchange_rates_client.Client
	{
		exchangeRatesRepository, err := exchange_rates_repository.Prepare(context.Background(), db)
		if err != nil {
			log.Fatalf("can't prepare exchange rates repository: %v", err)
		}

		exchangeRatesServer, err := exchange_rates.NewExchangeRatesServer(
			exchangeRatesRepository,
		)
		if err != nil {
			log.Fatalf("can't create exchange rates server: %v\n", err)
		}
		exchangeRatesClient = exchange_rates_client.New(exchangeRatesServer)
	}
	_ = exchangeRatesClient

	var walletSvcClient *walletClient.Client
	// Wallet service
	{
		walletRepository, err := walletRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("walletRepo error: %v", err)
		}

		feePayerPk, err := base64.StdEncoding.DecodeString(a.cfg.SolanaFeePayerPrivateKey)
		if err != nil {
			log.Fatalf("feePayerPk base64 decoding error: %v", err)
		}
		tokenHolderPk, err := base64.StdEncoding.DecodeString(a.cfg.SolanaTokenHolderPrivateKey)
		if err != nil {
			log.Fatalf("tokenHolderPk base64 decoding error: %v", err)
		}

		solanaClient := solana_client.New(a.cfg.SolanaApiBaseUrl, solana_client.Config{
			SystemProgram:   a.cfg.SolanaSystemProgram,
			SysvarRent:      a.cfg.SolanaSysvarRent,
			SysvarClock:     a.cfg.SolanaSysvarClock,
			SplToken:        a.cfg.SolanaSplToken,
			StakeProgramID:  a.cfg.SolanaStakeProgramID,
			TokenHolderAddr: a.cfg.SolanaTokenHolderAddr,
		})
		if err := solanaClient.CheckPrivateKey(a.cfg.SolanaFeePayerAddr, feePayerPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: fee payer: %v", err)
		}
		if err := solanaClient.CheckPrivateKey(a.cfg.SolanaTokenHolderAddr, tokenHolderPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: token holder: %v", err)
		}

		walletService := wallet.NewService(
			walletRepository,
			solanaClient,
			ethereumClient,
			wallet.WithAssetSolanaAddress(a.cfg.SolanaAssetAddr),
			wallet.WithSolanaFeePayer(a.cfg.SolanaFeePayerAddr, feePayerPk),
			wallet.WithSolanaTokenHolder(a.cfg.SolanaTokenHolderAddr, tokenHolderPk),
			wallet.WithMinAmountToTransfer(a.cfg.MinAmountToTransfer),
			wallet.WithStakePoolSolanaAddress(a.cfg.SolanaStakePoolAddr),
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
		rewards.WithExplorerURLTmpl("https://explorer.solana.com/tx/%s?cluster="+a.cfg.SolanaEnv),
		rewards.WithHoldRewardsPeriod(a.cfg.HoldRewardsPeriod),
		rewards.WithMinAmountToClaim(a.cfg.MinAmountToClaim),
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
		InvitationReward: a.cfg.InvitationReward,
		InvitationURL:    a.cfg.InvitationURL,
	})
	invitationsClient := invitationsClient.New(invitationsService)
	r.Mount("/invitations", invitations.MakeHTTPHandler(
		invitations.MakeEndpoints(invitationsService, jwtMdw),
		logger,
	))

	var authClient *authc.Client
	var nftClient *nftC.Client
	{
		// KYC
		kycService := sumsub.New(a.cfg.KycAppToken, a.cfg.KycAppSecret, a.cfg.KycAppBaseURL, a.cfg.KycAppTTL)
		kycClient := sumsub.NewClient(kycService)

		authService := auth.NewService(
			jwtInteractor,
			authRepository,
			walletSvcClient,
			invitationsClient,
			kycClient,
			auth.WithMasterOTPCode(a.cfg.MasterOTPHash),
			auth.WithCustomOTPLength(a.cfg.OtpLength),
			auth.WithMailService(mailer),
			auth.WithBlacklistMode(a.cfg.BlacklistMode),
			auth.WithWhitelistMode(a.cfg.WhitelistMode),
		)

		// Auth service
		{
			r.Mount("/auth", auth.MakeHTTPHandler(
				auth.MakeEndpoints(authService, jwtMdw),
				logger,
			))
		}

		authClient = authc.New(authService)
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
			BaseFirebaseURL:    a.cfg.BaseFirebaseURL,
			WebAPIKey:          a.cfg.FBWebAPIKey,
			MainSiteLink:       a.cfg.MainSiteLink,
			AndroidPackageName: a.cfg.AndroidPackageName,
			IosBundleId:        a.cfg.IOSBundleId,
			SuffixOption:       a.cfg.SuffixOption,
		})

		// Referrals service
		referralRepository, err := referralsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("referralRepo error: %v", err)
		}
		r.Mount("/ref", referrals.MakeHTTPHandler(
			referrals.MakeEndpoints(referrals.NewService(referralRepository, fb, firebase.Config{
				BaseFirebaseURL:    a.cfg.BaseFirebaseURL,
				WebAPIKey:          a.cfg.FBWebAPIKey,
				MainSiteLink:       a.cfg.MainSiteLink,
				AndroidPackageName: a.cfg.AndroidPackageName,
				IosBundleId:        a.cfg.IOSBundleId,
				SuffixOption:       a.cfg.SuffixOption,
			}), jwtMdw),
			logger,
		))
	}

	// Challenge client instance
	var challengeSvcClient *challengeClient.Client

	{
		// NFT service
		nftRepository, err := nftRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("nftRepo error: %v", err)
		}
		nftService := nft.NewService(nftRepository, walletSvcClient.PayForNFT)
		r.Mount("/nft", nft.MakeHTTPHandler(
			nft.MakeEndpoints(nftService, jwtMdw),
			logger,
		))
		nftClient = nftC.New(nftService)
	}

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
				fmt.Sprintf("%s/challenges", strings.TrimSuffix(a.cfg.AppBaseURL, "/")),
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

		showsService := shows.NewService(showRepo, challengeSvcClient, profileSvc, authClient, walletSvcClient.P2PTransfer, nftClient, a.cfg.TipsPercent)
		r.Mount("/shows", shows.MakeHTTPHandler(
			shows.MakeEndpoints(showsService, jwtMdw),
			logger,
		))
		r.Mount("/nft-marketplace/shows", private.MakeHTTPHandler(
			private.MakeEndpoints(showsService, jwt.NewAPIKeyMdw(a.cfg.SatorAPIKey)),
			logger,
		))
	}

	// files service
	var fileSvc *files.Service
	{
		opt := storage.Options{
			Key:            a.cfg.FileStorageKey,
			Secret:         a.cfg.FileStorageSecret,
			Endpoint:       a.cfg.FileStorageEndpoint,
			Region:         a.cfg.FileStorageRegion,
			Bucket:         a.cfg.FileStorageBucket,
			URL:            a.cfg.FileStorageUrl,
			DisableSSL:     a.cfg.FileStorageDisableSsl,
			ForcePathStyle: a.cfg.FileStorageForcePathStyle,
		}
		stor := storage.New(storage.NewS3Client(opt), opt)

		mediaServiceRepo, err := filesRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("mediaServiceRepo error: %v", err)
		}

		fileSvc = files.NewService(mediaServiceRepo, stor, resizer.Resize)

		r.Mount("/files", files.MakeHTTPHandler(
			files.MakeEndpoints(fileSvc, jwtMdw),
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

	{
		quizV2Repository, err := quizV2Repo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("can't prepare Quiz V2 repository: %v", err)
		}

		quizV2Svc := quiz_v2.NewService(
			a.cfg.NatsURL,
			a.cfg.NatsWSURL,
			challengeSvcClient,
			walletSvcClient,
			rewardsSvcClient,
			authClient,
			profileSvc,
			db,
			quizV2Repository,
			serverRSAPrivateKey,
			a.cfg.QuizV2ShuffleQuestions,
			a.cfg.QuizLobbyLatency,
		)
		r.Mount("/quiz_v2", quiz_v2.MakeHTTPHandler(
			quiz_v2.MakeEndpoints(quizV2Svc, jwtMdw),
			logger,
		))

		go quizV2Svc.StartEngine()
		// TODO(evg): gracefully shutdown the engine
	}

	{
		tradingPlatformsRepository, err := tradingPlatformsRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("can't prepare trading platforms repository: %v", err)
		}

		tradingPlatformsSvc := trading_platforms.NewService(
			tradingPlatformsRepository,
		)
		r.Mount("/trading_platforms", trading_platforms.MakeHTTPHandler(
			trading_platforms.MakeEndpoints(tradingPlatformsSvc, jwtMdw),
			logger,
		))
	}

	{
		puzzleGameRepository, err := puzzleGameRepo.Prepare(ctx, db)
		if err != nil {
			log.Fatalf("can't prepare puzzle game repository: %v", err)
		}

		puzzleGameSvc := puzzle_game.NewService(
			puzzleGameRepository,
			puzzle_game.WithChargeFunction(walletSvcClient.PayForService),
			puzzle_game.WithRewardsFunction(rewardsSvcClient.AddDepositTransaction),
			puzzle_game.WithFileServiceClient(fileSvc),
			puzzle_game.WithUserMultiplierFunction(walletSvcClient.GetMultiplier),
		)

		r.Mount("/puzzle-game", puzzle_game.MakeHTTPHandler(
			puzzle_game.MakeEndpoints(puzzleGameSvc, jwtMdw),
			logger,
		))
	}

	{
		// Init and run http server
		httpServer := &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", a.cfg.AppPort),
		}
		g.Add(func() error {
			log.Printf("[http-server] start listening on :%d...\n", a.cfg.AppPort)
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
			circulatingSupply = a.cfg.TokenCirculatingSupply
			for {
				select {
				case <-tickerDone:
					return nil
				case <-ticker.C:
					circulatingSupply++
				}
			}
		}, func(err error) {
			fmt.Println("going to shutdown ticker")
			tickerDone <- true
			fmt.Println("ticker is shutdown")
		})
	}

	{
		g.Add(func() error {
			<-a.done
			return nil
		}, func(err error) {
			fmt.Println("going to shutdown app")
			a.Shutdown()
			fmt.Println("app is shutdown")
		})
	}

	if err := g.Run(); err != nil {
		log.Println("API terminated with error:", err)
	}
}

func (a *app) Shutdown() {
	if a.shutdown {
		return
	}
	a.shutdown = true

	close(a.done)
}

// returns current build tag
func mkRootHandler(buildTagDO string) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if buildTag == "" {
			buildTag = buildTagDO
		}
		defaultResponse(w, http.StatusOK, map[string]interface{}{"build_tag": buildTag})
	}
}

// returns token circulating supply
func supplyHandler(w http.ResponseWriter, _ *http.Request) {
	defaultResponse(w, http.StatusOK, map[string]interface{}{
		"supply": circulatingSupply,
	})
}

// returns html page to test websocket
// func testWsHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./cmd/api/index.html")
// }

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
