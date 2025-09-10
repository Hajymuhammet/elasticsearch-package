package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	es "github.com/elastic/go-elasticsearch/v8"
)

// ClientConfig holds ES connection config
type ClientConfig struct {
	Addresses []string
	Username  string
	Password  string
	Timeout   time.Duration
	TLSConfig *tls.Config
}

// NewClient creates a new ES client with health check
func NewClient(cfg ClientConfig) (*es.Client, error) {
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("elasticsearch: no addresses provided")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	esCfg := es.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			TLSClientConfig: cfg.TLSConfig,
			DialContext: (&net.Dialer{
				Timeout:   cfg.Timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: cfg.Timeout,
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          100,
		},
	}

	client, err := es.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	// Health check with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if _, err := client.Ping(client.Ping.WithContext(ctx)); err != nil {
		return nil, err
	}

	return client, nil
}
