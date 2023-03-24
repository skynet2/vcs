/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package stress

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/trustbloc/logutil-go/pkg/log"

	"github.com/trustbloc/vcs/component/wallet-cli/pkg/walletrunner"
	"github.com/trustbloc/vcs/component/wallet-cli/pkg/walletrunner/vcprovider"
)

type TestCase struct {
	walletRunner         *walletrunner.Service
	httpClient           *http.Client
	vcsAPIURL            string
	issuerProfileID      string
	verifierProfileID    string
	credentialTemplateID string
	credentialType       string
	credentialFormat     string
	token                string
	claimData            map[string]interface{}
}

type TestCaseOptions struct {
	vcProviderOptions    []vcprovider.ConfigOption
	httpClient           *http.Client
	vcsAPIURL            string
	issuerProfileID      string
	verifierProfileID    string
	credentialTemplateID string
	credentialType       string
	credentialFormat     string
	token                string
	claimData            map[string]interface{}
}

type TestCaseOption func(opts *TestCaseOptions)

func NewTestCase(options ...TestCaseOption) (*TestCase, error) {
	opts := &TestCaseOptions{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		credentialFormat: "jwt_vc_json-ld",
	}

	for _, opt := range options {
		opt(opts)
	}

	if opts.vcsAPIURL == "" {
		return nil, fmt.Errorf("vcs api url is empty")
	}

	if opts.issuerProfileID == "" {
		return nil, fmt.Errorf("issuer profile id is empty")
	}

	if opts.verifierProfileID == "" {
		return nil, fmt.Errorf("verifier profile id is empty")
	}

	if opts.credentialType == "" {
		return nil, fmt.Errorf("credential type is empty")
	}

	runner, err := walletrunner.New(vcprovider.ProviderVCS, opts.vcProviderOptions...)
	if err != nil {
		return nil, fmt.Errorf("create wallet runner: %w", err)
	}

	return &TestCase{
		walletRunner:         runner,
		httpClient:           opts.httpClient,
		vcsAPIURL:            opts.vcsAPIURL,
		issuerProfileID:      opts.issuerProfileID,
		verifierProfileID:    opts.verifierProfileID,
		credentialTemplateID: opts.credentialTemplateID,
		credentialType:       opts.credentialType,
		credentialFormat:     opts.credentialFormat,
		token:                opts.token,
		claimData:            opts.claimData,
	}, nil
}

func WithVCProviderOption(opt vcprovider.ConfigOption) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.vcProviderOptions = append(opts.vcProviderOptions, opt)
	}
}

func WithHTTPClient(client *http.Client) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.httpClient = client
	}
}

func WithVCSAPIURL(apiURL string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.vcsAPIURL = apiURL
	}
}

func WithIssuerProfileID(issuerProfileID string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.issuerProfileID = issuerProfileID
	}
}

func WithVerifierProfileID(verifierProfileID string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.verifierProfileID = verifierProfileID
	}
}

func WithCredentialTemplateID(credentialTemplateID string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.credentialTemplateID = credentialTemplateID
	}
}

func WithCredentialType(credentialType string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.credentialType = credentialType
	}
}

func WithClaimData(data map[string]interface{}) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.claimData = data
	}
}

func WithToken(token string) TestCaseOption {
	return func(opts *TestCaseOptions) {
		opts.token = token
	}
}

type stressTestPerfInfo map[string]time.Duration

func (c *TestCase) Invoke() (interface{}, error) {
	credentialOfferURL, pin, err := c.fetchCredentialOfferURL()
	if err != nil {
		return nil, fmt.Errorf("fetch credential offer url: %w", err)
	}

	// run pre-auth flow and save credential in the wallet
	if err = c.walletRunner.RunOIDC4CIPreAuth(&walletrunner.OIDC4CIConfig{
		InitiateIssuanceURL: credentialOfferURL,
		CredentialType:      c.credentialType,
		CredentialFormat:    c.credentialFormat,
		Pin:                 pin,
	}); err != nil {
		return nil, fmt.Errorf("run pre-auth issuance: %w", err)
	}

	providerConf := c.walletRunner.GetConfig()
	providerConf.WalletUserId = providerConf.WalletParams.UserID
	providerConf.WalletPassPhrase = providerConf.WalletParams.Passphrase
	providerConf.WalletDidID = providerConf.WalletParams.DidID
	providerConf.WalletDidKeyID = providerConf.WalletParams.DidKeyID
	providerConf.SkipSchemaValidation = true

	authorizationRequest, err := c.fetchAuthorizationRequest()
	if err != nil {
		return nil, fmt.Errorf("fetch authorization request: %w", err)
	}

	err = c.walletRunner.RunOIDC4VPFlow(authorizationRequest)
	if err != nil {
		return nil, fmt.Errorf("run vp: %w", err)
	}

	b, err := json.Marshal(c.walletRunner.GetPerfInfo())
	if err != nil {
		return nil, fmt.Errorf("marshal perf info: %w", err)
	}

	var perfInfo stressTestPerfInfo

	if err = json.Unmarshal(b, &perfInfo); err != nil {
		return nil, fmt.Errorf("unmarshal perf info into stressTestPerfInfo: %w", err)
	}

	return perfInfo, nil
}

func (c *TestCase) fetchCredentialOfferURL() (string, string, error) {
	b, err := json.Marshal(&initiateOIDC4CIRequest{
		ClaimData:            &c.claimData,
		CredentialTemplateId: c.credentialTemplateID,
		UserPinRequired:      true,
	})
	if err != nil {
		return "", "", fmt.Errorf("marshal initiate oidc4ci request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(
			"%v/issuer/profiles/%v/interactions/initiate-oidc",
			c.vcsAPIURL,
			c.issuerProfileID,
		),
		bytes.NewBuffer(b))
	if err != nil {
		return "", "", fmt.Errorf("create initiate oidc4ci request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("send initiate oidc4ci request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("initiate oidc4ci request failed: %v", resp.Status)
	}

	if resp.Body != nil {
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				logger.Error("failed to close response body", log.WithError(err))
			}
		}()
	}

	var parsedResp initiateOIDC4CIResponse

	if err = json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return "", "", fmt.Errorf("decode initiate oidc4ci response: %w", err)
	}

	pin := ""
	if parsedResp.UserPin != nil {
		pin = *parsedResp.UserPin
	}

	return parsedResp.OfferCredentialURL, pin, nil
}

func (c *TestCase) fetchAuthorizationRequest() (string, error) {
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(
			"%s/verifier/profiles/%s/interactions/initiate-oidc",
			c.vcsAPIURL,
			c.verifierProfileID,
		),
		http.NoBody)
	if err != nil {
		return "", fmt.Errorf("create initiate oidc4vp request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send initiate oidc4vp request: %w", err)
	}

	if resp.Body != nil {
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				logger.Error("failed to close response body", log.WithError(err))
			}
		}()
	}

	var parsedResp initiateOIDC4VPResponse

	if err = json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return "", fmt.Errorf("decode initiate oidc4vp response: %w", err)
	}

	return parsedResp.AuthorizationRequest, nil
}
