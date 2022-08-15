module github.com/SatorNetwork/sator-api

go 1.17

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.2+incompatible

require (
	filippo.io/edwards25519 v1.0.0-rc.1
	firebase.google.com/go v3.13.0+incompatible
	github.com/SatorNetwork/gopuzzlegame v0.0.0-20220429113459-e0ae698bfb8a
	github.com/anytypeio/go-slip10 v0.0.0-20200330114100-25f30c832993
	github.com/awa/go-iap v1.3.16
	github.com/aws/aws-sdk-go v1.31.4
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/dmitrymomot/go-env v0.1.1
	github.com/dmitrymomot/go-signature v0.0.0-20190805202920-725320ef06d3
	github.com/dmitrymomot/random v1.0.6
	github.com/dustin/go-broadcast v0.0.0-20171205050544-f664265f5a66
	github.com/ethereum/go-ethereum v1.10.17
	github.com/go-chi/httprate v0.6.0
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.6.1
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/keighl/postmark v0.0.0-20190821160221-28358b1a94e3
	github.com/lib/pq v1.10.1
	github.com/mcnijman/go-emailaddress v1.1.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/mr-tron/base58 v1.2.0
	github.com/nats-io/nats.go v1.15.0
	github.com/near/borsh-go v0.3.1-0.20210831082424-4377deff6791
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/oklog/run v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/portto/solana-go-sdk v1.16.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea
	github.com/segmentio/ksuid v1.0.4
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.7.0
	github.com/superoo7/go-gecko v1.0.0
	github.com/thedevsaddam/govalidator v1.9.10
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/zeebo/errs v1.2.2
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	golang.org/x/text v0.3.7
	google.golang.org/api v0.59.0
	syreclabs.com/go/faker v1.2.3
)

require github.com/cespare/xxhash/v2 v2.1.2 // indirect

require (
	cloud.google.com/go v0.97.0 // indirect
	cloud.google.com/go/firestore v1.6.1 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/nats-io/nats-server/v2 v2.8.4 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.opencensus.io v0.23.0 // indirect
	goji.io v2.0.2+incompatible // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211028162531-8db9c33dc351 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
