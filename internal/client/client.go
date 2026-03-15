package client

import (
	"fmt"

	"github.com/kavinkumar999/kite-cli/internal/config"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

type Client struct {
	kc  *kiteconnect.Client
	cfg *config.Config
}

func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	kc := kiteconnect.New(cfg.APIKey)

	if cfg.AccessToken != "" {
		kc.SetAccessToken(cfg.AccessToken)
	}

	return &Client{
		kc:  kc,
		cfg: cfg,
	}, nil
}

func (c *Client) Kite() *kiteconnect.Client {
	return c.kc
}

func (c *Client) GenerateSession(requestToken string) error {
	data, err := c.kc.GenerateSession(requestToken, c.cfg.APISecret)
	if err != nil {
		return fmt.Errorf("failed to generate session: %w", err)
	}

	c.cfg.AccessToken = data.AccessToken
	c.kc.SetAccessToken(data.AccessToken)

	return config.Save(c.cfg)
}
