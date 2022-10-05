package vc_devapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/hyperledger/aries-framework-go/pkg/common/log"

	bddcontext "github.com/trustbloc/vcs/test/bdd/pkg/context"
)

var logger = log.New("bdd-test")

type Steps struct {
	bddContext     *bddcontext.BDDContext
	responseStatus int
	responseBody   []byte
}

// NewSteps returns new Steps context.
func NewSteps(ctx *bddcontext.BDDContext) *Steps {
	return &Steps{bddContext: ctx}
}

// RegisterSteps registers VC scenario steps.
func (s *Steps) RegisterSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I request did config for "([^"]*)" with ID "([^"]*)"$`, s.httpGet)
	//sc.Step(`^I receive response with status code "([^"]*)"$`, s.checkResponseStatus)
	//sc.Step(`^response contains "([^"]*)" with value "([^"]*)"$`, s.checkResponseValue)
}

func (s *Steps) httpGet(profileType string, id string) error {
	logger.Infof("I request did config : %v and id %v", string(profileType), id)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("http://localhost:8075/%s/profiles/%s/well-known/did-config",
			profileType, id), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Errorf("Failed to close response body: %s\n", closeErr.Error())
		}
	}()

	s.responseStatus = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	s.responseBody = body

	logger.Infof("body : %v", string(body))

	return nil
}
