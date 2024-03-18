/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc4ci

import (
	"bytes"
	"context"
	"text/template"

	"github.com/google/uuid"
	"github.com/trustbloc/vc-go/verifiable"
)

type CredentialComposer struct {
}

func NewCredentialComposer() *CredentialComposer {
	return &CredentialComposer{}
}

func (c *CredentialComposer) Compose(
	_ context.Context,
	credential *verifiable.Credential,
	tx *Transaction,
	txCredentialConfiguration *TxCredentialConfiguration,
	_ *PrepareCredentialRequest,
) (*verifiable.Credential, error) {
	if txCredentialConfiguration == nil || txCredentialConfiguration.CredentialComposeConfiguration == nil {
		return credential, nil
	}

	if idTemplate := txCredentialConfiguration.CredentialComposeConfiguration.IDTemplate; idTemplate != "" {
		id, err := c.renderRaw(idTemplate, c.baseParams(tx))
		if err != nil {
			return nil, err
		}

		credential = credential.WithModifiedID(id)
	}

	if txCredentialConfiguration.CredentialComposeConfiguration.OverrideIssuer {
		credential = credential.WithModifiedIssuer(&verifiable.Issuer{ID: tx.DID})
	}

	return credential, nil
}

func (c *CredentialComposer) baseParams(
	tx *Transaction,
) map[string]interface{} {
	result := map[string]interface{}{
		"RandomID":  uuid.NewString(),
		"TxID":      tx.ID,
		"IssuerDID": tx.DID,
	}

	return result
}

func (c *CredentialComposer) renderRaw(
	templateStr string,
	param map[string]interface{},
) (string, error) {
	tpl, err := template.New(uuid.NewString()).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, param); err != nil {
		return "", err
	}

	return buf.String(), nil
}
