package elasticsearch

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	es "github.com/elastic/go-elasticsearch/v8"
)

type ClientConfig struct {
	Addresses []string
	Username  string
	Password  string
	Timeout   time.Duration
	TLSConfig *tls.Config
}

func NewClient(cfg ClientConfig) (*es.Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	esCfg := es.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			TLSClientConfig: cfg.TLSConfig,
		},
	}

	client, err := es.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	_, err = client.Info()
	if err != nil {
		log.Printf("warning: es client Info() error: %v", err)
	}

	return client, nil
}
