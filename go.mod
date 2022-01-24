module github.com/tinkerbell/pbnj

go 1.16

require (
	bou.ke/monkey v1.0.2
	github.com/bmc-toolbox/bmclib v0.4.16-0.20211230160158-5afdbf3b6a65
	github.com/cristalhq/jwt/v3 v3.0.9
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/gebn/bmc v0.0.0-20200904230046-a5643220ab2a
	github.com/gin-gonic/gin v1.6.3
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2
	github.com/go-test/deep v1.0.7
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jacobweinstock/registrar v0.4.5
	github.com/manifoldco/promptui v0.8.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/onsi/gomega v1.10.4
	github.com/packethost/pkg v0.0.0-20211110202003-387414657e83
	github.com/packethost/pkg/grpc/authz v0.0.0-20211110202003-387414657e83
	github.com/packethost/pkg/log/logr v0.0.0-20211110202003-387414657e83
	github.com/philippgille/gokv v0.6.0
	github.com/philippgille/gokv/freecache v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/xid v1.2.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stmcginnis/gofish v0.12.0
	github.com/stretchr/testify v1.7.0
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/zap v1.19.1
	goa.design/goa v2.2.5+incompatible
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/packethost/pkg/log/logr => ../../joelrebel/pkg/log/logr/v1
