module xdao.top/golight

go 1.16

require (
	github.com/coreos/bbolt v1.3.6 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elfincafe/mbstring v0.4.2
	github.com/forgoer/openssl v1.2.0
	github.com/gabriel-vasile/mimetype v1.4.0
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd v3.3.27+incompatible
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/text v0.3.6
	google.golang.org/grpc v1.33.1
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
