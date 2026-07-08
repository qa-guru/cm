module github.com/aerokube/cm

go 1.26

toolchain go1.26.5

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/aerokube/selenoid v0.0.0-20240520175821-773c202b01e3
	github.com/docker/go-connections v0.7.0
	github.com/fatih/color v1.17.0
	github.com/fvbommel/sortorder v1.1.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/heroku/docker-registry-client v0.0.0-20211012143308-9463674c8930
	github.com/mattn/go-colorable v0.1.13
	github.com/mitchellh/go-ps v1.0.0
	github.com/moby/moby/api v1.55.0
	github.com/moby/moby/client v0.5.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/text v0.25.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/clipperhouse/uax29/v2 v2.2.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/distribution v0.0.0-20171011171712-7484e51bf6af // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/aerokube/selenoid => github.com/qa-guru/selenoid v0.0.0-20260707221418-7fe0173cf3fd
