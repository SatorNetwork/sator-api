module github.com/SatorNetwork/sator-api

go 1.17

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.2+incompatible

require (
	filippo.io/edwards25519 v1.0.0-rc.1
	github.com/anytypeio/go-slip10 v0.0.0-20200330114100-25f30c832993
	github.com/aws/aws-sdk-go v1.31.4
	github.com/cenkalti/backoff v2.2.1+incompatible
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
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/golang/mock v1.5.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/keighl/postmark v0.0.0-20190821160221-28358b1a94e3
	github.com/lib/pq v1.10.1
	github.com/mcnijman/go-emailaddress v1.1.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/mr-tron/base58 v1.2.0
	github.com/nats-io/nats.go v1.9.1
	github.com/near/borsh-go v0.3.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/oklog/run v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/portto/solana-go-sdk v0.1.0
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/thedevsaddam/govalidator v1.9.10
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/zeebo/errs v1.2.2
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	syreclabs.com/go/faker v1.2.3
)

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea // indirect
	github.com/denisenkom/go-mssqldb v0.9.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/go-gorp/gorp/v3 v3.0.2 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/godror/godror v0.24.2 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-oci8 v0.1.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/mitchellh/cli v1.1.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nkeys v0.1.3 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/posener/complete v1.1.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/teserakt-io/golang-ed25519 v0.0.0-20210104091850-3888c087a4c8 // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	goji.io v2.0.2+incompatible // indirect
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
