// Package client provides intergation with api calls.
package client

import (
	"errors"
	"fmt"
)

// constants for the client config.
const (
	USRegion        = "us"
	EURegion        = "eu"
	BaseURLTemplate = "https://api.%s.onelogin.com"
	DefaultTimeout  = 5
)

var (
	errRegion            = errors.New("region is missing or unsupported")
	errClientIDEmpty     = errors.New("client_id is missing")
	errClientSecretEmpty = errors.New("client_secret is missing")
)

// APIClientConfig is the configuration for the APIClient.
type APIClientConfig struct {
	Timeout      int
	ClientID     string
	ClientSecret string
	Region       string
	Url          string
}

func (cfg *APIClientConfig) Initialize() (*APIClientConfig, error) {

	// Initialize clientID
	if len(cfg.ClientID) == 0 {
		return cfg, errClientIDEmpty
	}

	// Initialize clientSecret
	if len(cfg.ClientSecret) == 0 {
		return cfg, errClientSecretEmpty
	}
	if len(cfg.Url) == 0 {
		// Initialize the region if no url given
		if !isSupportedRegion(cfg.Region) {
			return cfg, errRegion
		}
		cfg.Url = fmt.Sprintf(BaseURLTemplate, cfg.Region)
	}

	// Initialize the timeout
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	return cfg, nil
}

func isSupportedRegion(region string) bool {
	return region == EURegion || region == USRegion
}
