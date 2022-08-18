package srv

import (
	"crypto/tls"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type EtcdConfig struct {
	Host       string
	Username   string
	Password   string
	UseTls     bool
	CertFile   string
	KeyFile    string
	CaCertFile string
	Timeout    int
}

func ParseEtcdCfg(host string) (*EtcdConfig, error) {
	var cfg EtcdConfig
	// cfg.Host = host
	cfg.UseTls = false
	cfg.Timeout = 5
	host = strings.Replace(host, "etcd://", "", 1)
	param := strings.Split(host, "?")
	pos := strings.Index(param[0], "@")
	if pos == -1 {
		return nil, fmt.Errorf("invalid etcd host no @: %s", host)
	}
	userInfo := strings.Split(param[0][:pos], ":")
	if len(userInfo) != 2 {
		return nil, fmt.Errorf("invalid etcd host no userinfo: %s", host)
	}
	cfg.Username = userInfo[0]
	cfg.Password = userInfo[1]
	cfg.Host = strings.Trim(param[0][pos+1:], "()")
	if len(param) == 2 {
		// param[1]
		v, err := url.ParseQuery(param[1])
		if err != nil {
			return nil, fmt.Errorf("invalid etcd host param err: %s", err.Error())
		}
		if v.Get("useTls") == "true" {
			cfg.UseTls = true
		}
		if v.Get("certFile") != "" {
			cfg.CertFile = v.Get("certFile")
		}
		if v.Get("keyFile") != "" {
			cfg.KeyFile = v.Get("keyFile")
		}
		if v.Get("caCertFile") != "" {
			cfg.CaCertFile = v.Get("caCertFile")
		}
		if v.Get("timeout") != "" {
			cfg.Timeout, _ = strconv.Atoi(v.Get("timeout"))
		}
	}

	return &cfg, nil
}

// NewEtcdClient creates a new etcd client.
// etcdDsn: connect string like  etcd://user:password@host:port?useTls=true&certFile=xxx&keyFile=xxx&caCertFile=xxx&timeout=xxx
func NewEtcdClient(etcdDsn string) (*clientv3.Client, error) {
	var err error
	config, err := ParseEtcdCfg(etcdDsn)
	if err != nil {
		return nil, err
	}
	endpoints := []string{config.Host}

	// use tls if usetls is true
	var tlsConfig *tls.Config
	if config.UseTls {
		tlsInfo := transport.TLSInfo{
			CertFile:           config.CertFile,
			KeyFile:            config.KeyFile,
			TrustedCAFile:      config.CaCertFile,
			InsecureSkipVerify: true,
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			log.Println(err.Error())
		}
	}

	conf := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * time.Duration(config.Timeout),
		TLS:         tlsConfig,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	if config.Username != "" {
		conf.Username = config.Username
		conf.Password = config.Password
	}

	var c *clientv3.Client
	c, err = clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	return c, nil
}
