package app_config

import (
	"time"

	sator_app "github.com/SatorNetwork/sator-api/cmd/api/app"
)

const serverRSAPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAyzOP1bef/mRQELQSU99VMXx5JSS/uoEjDVUYbv54qZ0ZNr/x
mr/LVDmklAH85YXoVI9Ze116ms4o27/iSyDRDQzZQJMZQeZkhsrCjl047KCEtuSB
KlrWV514C0PZAI8E5dIOmcThqDL6Oe3ko3cxc5On27O8i1L02ZR63lzAHsiYUKZJ
LCuumsgRjN3Qmnxgsj7blCNlo2Inc04hIU7FxNePnDQeIq9InTQwe8kpdrR3+z7O
M2dOwsgzwQXfRmhRVmMv9WyBOA7ppbhED2xDV0HtwlE7zvVXVzz0XlIH0Bcg2hJV
nOUspMNOVN98hovTNW1zYapq8K4WzdzTLWu6ewIDAQABAoIBAEk+mb8YhFRHRlDv
B5gx7Vx8GJLZ7z5D5CWfIcKIeWfyF6+TenwkXB9M89Th7o2xOfCZB1Effd0lKLcl
MGWSO6hvlDPhkf4eFOg9V+nHEojAx3XBYgmpWo+UVDwDRcocw1Av6lFlFU3yqh2N
cZe2lB2sAJqB7BlwWo+/JJbYaevuD6i6SQUy6kXrbDAFbMY0phcuL4qkW3KsywpY
2bijso8HeLVurk74MaAezbgwVzrhReA89RQruJtEIT4EYECknxammam5PuCSeflV
Mpsh31KP8PSu6Le4ejMYqKnQhLBV1kpA2tZ00NwP1T8PgRcxmBW+YhxDIK/PIwJv
lz0Z8FECgYEA1YVgVX3NZ6HMTmokM/u+rIjAg08meSVokUjnosvdef8RTEYT5DlZ
3RoCmKzth4T9GJE56Cvo+0C6Xvcba2QiHBCuvwa0NCOqNv37qtXsAOwKwBxVx8nn
9LSsrpmrXl407C64mM4pPXH5FMkHrK7AwGikxiceCWZFl4ULcZZDB/MCgYEA86Cb
5uu9Ah3i6X5Sm/wV4WAejLYKIgw046BRXF2ZJB+Rk66l9NG28hIgszQr4tgPVqQ8
m4u/xctII5fSLNN+fKxz7Pqb5uPDrWbBvyk4BAlz+oCV6DpQQrmX79OSSz3GLLxa
IG2CXYe41s02wJZ+WF6Ap/+SRt8qXEOkZpm27VkCgYApC2iHRpWTlECn2jN3Yq82
j1siYraslwpQ00jjvHiomOWEWfw85OFnZTaWjrdiU6grbs9I2BgDJGAvHSVCMY91
Aaf1xJ4jX6+6vnwATPr++mDeqRO8Qg26tnhzX8rXaxiVRi3qAcdfrmcJHdPB2B3p
XrQ9+wsFF4nNJKAch5v/DQKBgQCQtB6lg/OZpEK4yQ0sFQix+rNqhF10Z6eqY/iv
UfC54f5Hp35u8XkmQtolVqGSdR53KcnN4a2gP+OzMGPnuB7y0kNwyFF9TF9XSSde
8Y6R50N50JI5gxlU6IN0MUg9ZI2m2KD3jdPW1dxVyUHyFfEpb8gfAM/TRI4Wix7E
yhw60QKBgQCmmaR1Eg8/sJV/5SiAgHzXkOjPRmQjmwESazhe10mja0s/dyHk7Si3
x7YDoJeQm7c+GO6lb46/ccFU/MGXNmtEDeaLT5bVaZ81zkEWED2AjdHK6S4yttx/
o5uW70B4smBr2njkjibrTjaY1Mb7z4zFVPpK7ohEAulTnqKdHdFeqQ==
-----END RSA PRIVATE KEY-----`

var (
	AppConfigForTests = sator_app.Config{
		AppPort:                     8080,
		AppBaseURL:                  "XXXXXXXXXX",
		HttpRequestTimeout:          30 * time.Second,
		DBConnString:                "postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable",
		DBMaxOpenConns:              20,
		DBMaxIdleConns:              2,
		JwtSigningKey:               "secret",
		JwtTTL:                      24 * time.Hour,
		OtpLength:                   5,
		MasterOTPHash:               "$2a$04$JEj1CnjccUr237U8lOWMVOUPcm4xG/a3SHcJM00uNQKAx.ujaP5Pa",
		QuizWsConnURL:               "https://aec45cb3e117.ngrok.io/quiz",
		QuizBotsTimeout:             5 * time.Second,
		QuizLobbyLatency:            5 * time.Second,
		TokenCirculatingSupply:      11839844,
		SolanaEnv:                   "devnet",
		SolanaApiBaseUrl:            "http://localhost:8899/",
		SolanaAssetAddr:             "3yKB53R6DCuq2VL7aBfJY4VT9jv3w67NixyWoWoZZe5v",
		SolanaFeePayerAddr:          "67CXqkdKLhZxeDaHos2dxNGpqaiJvvva77TDnEipxXPx",
		SolanaFeePayerPrivateKey:    "tg3BEHU1lH24lo9JccmqLL13DLOzLMptxh0aa3NXJUtL4PVdkvwOmbpCqMTFG7a8CJles911d0uu7SYeuck8Uw==",
		SolanaTokenHolderAddr:       "uFhu3UDp2ymFYKRwPf1jrvfhDj1R7eiWjDnkdVQJhGQ",
		SolanaTokenHolderPrivateKey: "I52q0J0qsUY2NLTSScSKre1lH6XZRu69FGS0pa3xypsNYtRHIr9ICfw0SXUd1Vcr0sf3tqQuG3whne/UvJfBNQ==",
		SolanaStakePoolAddr:         "4pm3G48wWGrbUVF3JHLDrgVniQi7eSkRyx5bwXnawC2z",
		SolanaSystemProgram:         "11111111111111111111111111111111",
		SolanaSysvarRent:            "SysvarRent111111111111111111111111111111111",
		SolanaSysvarClock:           "SysvarC1ock11111111111111111111111111111111",
		SolanaSplToken:              "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		SolanaStakeProgramID:        "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
		PostmarkServerToken:         "b8de3efc-df98-4302-9d58-7bac75d7fdef",
		PostmarkAccountToken:        "local",
		NotificationFromName:        "Sator.io",
		NotificationFromEmail:       "notifications@sator.io",
		ProductName:                 "Sator.io",
		ProductURL:                  "https://sator.io",
		SupportURL:                  "https://sator.io",
		SupportEmail:                "support@sator.io",
		CompanyName:                 "Sator",
		CompanyAddress:              "New York",
		HoldRewardsPeriod:           0,
		InvitationReward:            0,
		InvitationURL:               "https://sator.io",
		FileStorageKey:              "XXXXXXXXXX",
		FileStorageSecret:           "XXXXXXXXXX",
		FileStorageEndpoint:         "XXXXXXXXXX",
		FileStorageRegion:           "XXXXXXXXXX",
		FileStorageBucket:           "XXXXXXXXXX",
		FileStorageUrl:              "XXXXXXXXXX",
		FileStorageDisableSsl:       false,
		FileStorageForcePathStyle:   true,
		BaseFirebaseURL:             "https://satorio.page.link",
		FBWebAPIKey:                 "XXXXXXXXXXXXXXXXXXXXXX",
		MainSiteLink:                "https://sator.io/",
		AndroidPackageName:          "com.satorio.app",
		IOSBundleId:                 "io.sator",
		SuffixOption:                "UNGUESSABLE",
		MinAmountToTransfer:         0,
		MinAmountToClaim:            0,
		KycAppToken:                 "XXXXXXXXXX",
		KycAppSecret:                "XXXXXXXXXX",
		KycAppBaseURL:               "XXXXXXXXXX",
		KycAppTTL:                   1200,
		KycSkip:                     false,
		NatsURL:                     "nats://127.0.0.1:4222",
		NatsWSURL:                   "ws://127.0.0.1:8080",
		QuizV2ShuffleQuestions:      false,
		ServerRSAPrivateKey:         serverRSAPrivateKey,
		SatorAPIKey:                 "582e89d8-69ca-4206-8e7f-1fc822b41307",
	}
)

func RunAndWait() func() {
	app := sator_app.WithConfig(&AppConfigForTests)
	go app.Run()
	time.Sleep(3 * time.Second)

	return func() {
		app.Shutdown()
		time.Sleep(5 * time.Second)
	}
}
