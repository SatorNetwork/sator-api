module github.com/SatorNetwork/sator-api

go 1.16

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.2+incompatible

require (
	github.com/aws/aws-sdk-go v1.31.4
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/dmitrymomot/distlock v0.1.1
	github.com/dmitrymomot/go-env v0.1.1
	github.com/dmitrymomot/go-signature v0.0.0-20190805202920-725320ef06d3
	github.com/dmitrymomot/random v0.0.0-20190806074213-235e86f90ac3
	github.com/dustin/go-broadcast v0.0.0-20171205050544-f664265f5a66
	github.com/ethereum/go-ethereum v1.10.12
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/keighl/postmark v0.0.0-20190821160221-28358b1a94e3
	github.com/lib/pq v1.10.1
	github.com/mr-tron/base58 v1.2.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/oklog/run v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/portto/solana-go-sdk v0.1.0
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea
	github.com/stretchr/testify v1.7.0
	github.com/thedevsaddam/govalidator v1.9.10
	github.com/zeebo/errs v1.2.2
	github.com/ziutek/mymysql v1.5.4 // indirect
	goji.io v2.0.2+incompatible // indirect
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	syreclabs.com/go/faker v1.2.3
)
