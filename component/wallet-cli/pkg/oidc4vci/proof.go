package oidc4vci

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/piprate/json-gold/ld"
	"github.com/samber/lo"
	vdrapi "github.com/trustbloc/did-go/vdr/api"
	"github.com/trustbloc/kms-go/doc/jose"
	"github.com/trustbloc/kms-go/doc/jose/jwk"
	"github.com/trustbloc/kms-go/spi/kms"
	"github.com/trustbloc/vc-go/cwt"
	"github.com/trustbloc/vc-go/dataintegrity"
	"github.com/trustbloc/vc-go/dataintegrity/suite/ecdsa2019"
	"github.com/trustbloc/vc-go/jwt"
	"github.com/trustbloc/vc-go/proof"
	"github.com/trustbloc/vc-go/proof/creator"
	"github.com/trustbloc/vc-go/proof/jwtproofs/eddsa"
	"github.com/trustbloc/vc-go/proof/jwtproofs/es256"
	"github.com/trustbloc/vc-go/proof/jwtproofs/es256k"
	"github.com/trustbloc/vc-go/proof/jwtproofs/es384"
	"github.com/trustbloc/vc-go/proof/jwtproofs/es521"
	"github.com/trustbloc/vc-go/proof/jwtproofs/ps256"
	"github.com/trustbloc/vc-go/proof/jwtproofs/rs256"
	"github.com/trustbloc/vc-go/proof/ldproofs/bbsblssignature2020"
	"github.com/trustbloc/vc-go/proof/ldproofs/ecdsasecp256k1signature2019"
	"github.com/trustbloc/vc-go/proof/ldproofs/ed25519signature2018"
	"github.com/trustbloc/vc-go/proof/ldproofs/ed25519signature2020"
	"github.com/trustbloc/vc-go/proof/ldproofs/jsonwebsignature2020"
	"github.com/trustbloc/vc-go/verifiable"
	cwt2 "github.com/trustbloc/vc-go/verifiable/cwt"
	"github.com/veraison/go-cose"
)

type ProofBuilder interface {
	Build(
		ctx context.Context,
		req *CreateProofRequest,
	) (*Proof, error)
}

type CWTProofBuilder struct {
}

func NewCWTProofBuilder() *CWTProofBuilder {
	return &CWTProofBuilder{}
}

func (b *CWTProofBuilder) newProofCreator(signer jose.Signer) (*creator.ProofCreator, []proof.JWTProofDescriptor) {
	desc := []proof.JWTProofDescriptor{
		eddsa.New(),
		es256.New(),
		es256k.New(),
		es384.New(),
		es521.New(),
		rs256.New(),
		ps256.New(),
	}

	opt := []creator.Opt{
		creator.WithLDProofType(bbsblssignature2020.New(), signer),
		creator.WithLDProofType(ecdsasecp256k1signature2019.New(), signer),
		creator.WithLDProofType(ed25519signature2018.New(), signer),
		creator.WithLDProofType(ed25519signature2020.New(), signer),
		creator.WithLDProofType(jsonwebsignature2020.New(), signer),
	}

	for _, d := range desc {
		opt = append(opt, creator.WithJWTAlg(d, signer))
	}

	return creator.New(
		opt...,
	), desc
}

func (b *CWTProofBuilder) Build(
	ctx context.Context,
	req *CreateProofRequest,
) (*Proof, error) {
	req.Claims.Issuer = req.WalletDID // per spec

	encoded, err := cbor.Marshal(req.Claims)
	if err != nil {
		return nil, fmt.Errorf("marshal proof claims: %w", err)
	}

	proofCreator, descriptors := b.newProofCreator(req.Signer)

	algo, _ := req.Signer.Headers().Algorithm()
	var targetAlgo cose.Algorithm
	for _, d := range descriptors {
		if d.JWTAlgorithm() == algo && d.CWTAlgorithm() != 0 {
			targetAlgo = d.CWTAlgorithm()
		}
	}

	if targetAlgo == 0 {
		return nil, fmt.Errorf("unsupported cosg algorithm: %s", algo)
	}

	coseKey, err := cose.NewKeyFromPublic(req.PubKeyJWK.Key)
	if err != nil {
		return nil, fmt.Errorf("get public key bytes: %w", err)
	}

	keyID, _ := req.Signer.Headers().KeyID()

	keyBytes, err := coseKey.MarshalCBOR()
	if err != nil {
		return nil, fmt.Errorf("COSE_KEY marshal cbor: %w", err)
	}
	pubKeyStr := hex.EncodeToString(keyBytes)

	msg := &cose.Sign1Message{
		Headers: cose.Headers{
			Protected: cose.ProtectedHeader{
				cose.HeaderLabelAlgorithm:   targetAlgo,
				cose.HeaderLabelContentType: "application/openid4vci-proof+cwt",
				proof.COSEKeyHeader:         pubKeyStr,
			},
		},
		Payload: encoded,
	}

	signData, err := cwt2.GetProofValue(msg)
	if err != nil {
		return nil, fmt.Errorf("get proof value: %w", err)
	}

	signed, err := proofCreator.SignCWT(cwt.SignParameters{
		KeyID:  keyID,
		CWTAlg: targetAlgo,
	}, signData)
	if err != nil {
		return nil, fmt.Errorf("sign cwt: %w", err)
	}

	msg.Signature = signed

	finalMsg, err := msg.MarshalCBOR()
	if err != nil {
		return nil, fmt.Errorf("marshal cbor: %w", err)
	}

	return &Proof{
		CWT:       hex.EncodeToString(finalMsg),
		ProofType: "cwt",
	}, nil
}

type JWTProofBuilder struct {
	proofFn JWTProofFn
}

type JWTProofFn func(
	ctx context.Context,
	req *CreateProofRequest,
) (string, error)

func NewJWTProofBuilder() *JWTProofBuilder {
	return &JWTProofBuilder{
		proofFn: func(
			ctx context.Context,
			req *CreateProofRequest,
		) (string, error) {
			signedJWT, jwtErr := jwt.NewJoseSigned(req.Claims, req.CustomHeaders, req.Signer)
			if jwtErr != nil {
				return "", fmt.Errorf("create signed jwt: %w", jwtErr)
			}

			jws, jwtErr := signedJWT.Serialize(false)
			if jwtErr != nil {
				return "", fmt.Errorf("serialize signed jwt: %w", jwtErr)
			}

			return jws, nil
		},
	}
}

func (b *JWTProofBuilder) WithCustomProofFn(
	proofFn JWTProofFn,
) *JWTProofBuilder {
	b.proofFn = proofFn

	return b
}

func (b *JWTProofBuilder) Build(
	ctx context.Context,
	req *CreateProofRequest,
) (*Proof, error) {
	headers := map[string]interface{}{
		jose.HeaderType: jwtProofTypeHeader,
	}
	for k, v := range headers {
		req.CustomHeaders[k] = v
	}

	jws, err := b.proofFn(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("build proof: %w", err)
	}

	return &Proof{
		JWT:       jws,
		ProofType: "jwt",
	}, nil
}

type LDPProofBuilder struct {
	proofFn JWTProofFn
}

func NewLDPProofBuilder() *LDPProofBuilder {
	return &LDPProofBuilder{}
}

func (b *LDPProofBuilder) Build(
	_ context.Context,
	req *CreateProofRequest,
) (*Proof, error) {
	pres, err := verifiable.NewPresentation()
	if err != nil {
		return nil, fmt.Errorf("new presentation: %w", err)
	}

	signerSuite := ecdsa2019.NewSignerInitializer(&ecdsa2019.SignerInitializerOptions{
		SignerGetter:     ecdsa2019.WithStaticSigner(req.Signer),
		LDDocumentLoader: ld.NewDefaultDocumentLoader(http.DefaultClient),
	})

	signer, err := dataintegrity.NewSigner(&dataintegrity.Options{
		DIDResolver: req.VDR,
	}, signerSuite)
	if err != nil {
		return nil, fmt.Errorf("new signer: %w", err)
	}

	pres.Holder = req.WalletDID
	if err = pres.AddDataIntegrityProof(&verifiable.DataIntegrityProofContext{
		SigningKeyID: fmt.Sprintf("%s#%s", req.WalletDID, req.WalletKeyID),
		CryptoSuite:  ecdsa2019.SuiteType,
		Created:      lo.ToPtr(time.Now().UTC()),
		Domain:       req.CredentialIssuer,
		Challenge:    req.Claims.Nonce,
	}, signer); err != nil {
		return nil, fmt.Errorf("add data integrity proof: %w", err)
	}

	return &Proof{
		LdpVp:     pres,
		ProofType: "ldp_vp",
	}, nil
}

type CreateProofRequest struct {
	Signer           jose.Signer
	CustomHeaders    map[string]interface{}
	WalletKeyID      string
	WalletDID        string
	WalletKeyType    kms.KeyType
	Claims           *ProofClaims
	VDR              vdrapi.Registry
	CredentialIssuer string
	PubKeyJWK        *jwk.JWK
}
