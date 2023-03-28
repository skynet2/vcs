/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//go:generate mockgen -destination service_mocks_test.go -self_package mocks -package issuecredential_test -source=issuecredential_service.go -mock_names profileService=MockProfileService,kmsRegistry=MockKMSRegistry,vcStatusManager=MockVCStatusManager

package issuecredential

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/trustbloc/logutil-go/pkg/log"

	"github.com/trustbloc/vcs/pkg/doc/vc"
	"github.com/trustbloc/vcs/pkg/doc/vc/crypto"
	"github.com/trustbloc/vcs/pkg/doc/vc/vcutil"
	vcskms "github.com/trustbloc/vcs/pkg/kms"
	profileapi "github.com/trustbloc/vcs/pkg/profile"
	"github.com/trustbloc/vcs/pkg/service/credentialstatus"
)

const (
	defaultCredentialPrefix = "urn:uuid:" //nolint:gosec
)

type vcCrypto interface {
	SignCredential(signerData *vc.Signer, vc *verifiable.Credential,
		opts ...crypto.SigningOpts) (*verifiable.Credential, error)
}

type kmsRegistry interface {
	GetKeyManager(config *vcskms.Config) (vcskms.VCSKeyManager, error)
}

type vcStatusManager interface {
	CreateStatusListEntry(ctx context.Context, profileID, credentialID string) (*credentialstatus.StatusListEntry, error)
}

type Config struct {
	VCStatusManager vcStatusManager
	Crypto          vcCrypto
	KMSRegistry     kmsRegistry
}

type Service struct {
	vcStatusManager vcStatusManager
	crypto          vcCrypto
	kmsRegistry     kmsRegistry
}

func New(config *Config) *Service {
	return &Service{
		vcStatusManager: config.VCStatusManager,
		crypto:          config.Crypto,
		kmsRegistry:     config.KMSRegistry,
	}
}

var logger = log.New("issuer-credential-svc")

func (s *Service) IssueCredential(
	ctx context.Context,
	credential *verifiable.Credential,
	issuerSigningOpts []crypto.SigningOpts,
	profile *profileapi.Issuer) (*verifiable.Credential, error) {
	kms, err := s.kmsRegistry.GetKeyManager(profile.KMSConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get kms: %w", err)
	}

	signer := &vc.Signer{
		DID:                     profile.SigningDID.DID,
		Creator:                 profile.SigningDID.Creator,
		KMSKeyID:                profile.SigningDID.KMSKeyID,
		SignatureType:           profile.VCConfig.SigningAlgorithm,
		KeyType:                 profile.VCConfig.KeyType,
		KMS:                     kms,
		Format:                  profile.VCConfig.Format,
		SignatureRepresentation: profile.VCConfig.SignatureRepresentation,
		VCStatusListType:        profile.VCConfig.Status.Type,
		SDJWT:                   profile.VCConfig.SDJWT,
	}

	var statusListEntry *credentialstatus.StatusListEntry

	// update credential prefix.
	vcutil.PrependCredentialPrefix(credential, defaultCredentialPrefix)

	if !profile.VCConfig.Status.Disable {
		st := time.Now()
		statusListEntry, err = s.vcStatusManager.CreateStatusListEntry(ctx, profile.ID, credential.ID)
		logger.Info(fmt.Sprintf("s.vcStatusManager.CreateStatusListEntry took %v", time.Since(st)))
		if err != nil {
			return nil, fmt.Errorf("failed to add credential status: %w", err)
		}

		credential.Context = append(credential.Context, statusListEntry.Context)
		credential.Status = statusListEntry.TypedID
	}

	// update context
	vcutil.UpdateSignatureTypeContext(credential, profile.VCConfig.SigningAlgorithm)

	// update credential issuer
	vcutil.UpdateIssuer(credential, profile.SigningDID.DID, profile.Name, true)

	// sign the credential
	st := time.Now()
	signedVC, err := s.crypto.SignCredential(signer, credential, issuerSigningOpts...)
	logger.Info(fmt.Sprintf("s.crypto.SignCredential took %v", time.Since(st)))
	if err != nil {
		return nil, fmt.Errorf("failed to sign credential: %w", err)
	}

	return signedVC, nil
}
