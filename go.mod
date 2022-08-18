module xdao.top/golight

go 1.16

require (
	github.com/elfincafe/mbstring v0.4.2
	github.com/forgoer/openssl v1.2.1
	github.com/gabriel-vasile/mimetype v1.4.1
	github.com/google/uuid v1.3.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tuotoo/qrcode v0.0.0-20190222102259-ac9c44189bf2
	go.etcd.io/etcd/client/pkg/v3 v3.6.0-alpha.0
	go.etcd.io/etcd/client/v3 v3.6.0-alpha.0.0.20220811010006-d4778e78c833
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.22.0 // indirect
	golang.org/x/crypto v0.0.0-20220817201139-bc19a97f63c8
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
	golang.org/x/sys v0.0.0-20220817070843-5a390386f1f2 // indirect
	golang.org/x/text v0.3.8-0.20211105212822-18b340fc7af2
	google.golang.org/genproto v0.0.0-20220817144833-d7fd3f11b9b1 // indirect
	google.golang.org/grpc v1.48.0
)

replace github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6

replace github.com/tuotoo/qrcode v0.0.0-20190222102259-ac9c44189bf2 => github.com/CuteReimu/qrcode v0.0.0-20220122115047-7a37f0abf050
