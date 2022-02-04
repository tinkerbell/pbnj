module github.com/tinkerbell/pbnj

go 1.16

require (
	github.com/bmc-toolbox/bmclib v0.4.16-0.20220202200536-079d247f718c
	github.com/cristalhq/jwt/v3 v3.1.0
	github.com/equinix-labs/otel-init-go v0.0.5
	github.com/fatih/color v1.13.0
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zerologr v1.2.1
	github.com/go-test/deep v1.0.8
	github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/manifoldco/promptui v0.9.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/onsi/gomega v1.18.1
	github.com/packethost/pkg/grpc/authz v0.0.0-20211110202003-387414657e83
	github.com/philippgille/gokv v0.6.0
	github.com/philippgille/gokv/freecache v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rs/xid v1.3.0
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	goa.design/goa v2.2.5+incompatible
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)
