/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//go:generate mockgen -destination oauth_client_mocks_test.go -self_package mocks -package oauth2client_test -source=oauth_client.go -mock_names OAuth2Client=MockOAuth2Client

package oauth2client

import (
	"context"

	"golang.org/x/oauth2"
)

type Factory struct {
}

type Client struct {
	*oauth2.Config
}

type OAuth2Client interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	GeneratePKCE() (verifier string, challenge string, method string, err error)
}

func NewOAuth2ClientFactory() *Factory {
	return &Factory{}
}

func (o *Factory) GetClient(config oauth2.Config) OAuth2Client {
	return &Client{
		Config: &config,
	}
}
