// Package common provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package common

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Defines values for DIDMethod.
const (
	DIDMethodKey DIDMethod = "key"
	DIDMethodOrb DIDMethod = "orb"
	DIDMethodWeb DIDMethod = "web"
)

// Defines values for KMSConfigType.
const (
	KMSConfigTypeAws   KMSConfigType = "aws"
	KMSConfigTypeLocal KMSConfigType = "local"
	KMSConfigTypeWeb   KMSConfigType = "web"
)

// Defines values for VCFormat.
const (
	JwtVcJsonLd VCFormat = "jwt_vc_json-ld"
	LdpVc       VCFormat = "ldp_vc"
)

// Defines values for VPFormat.
const (
	JwtVp VPFormat = "jwt_vp"
	LdpVp VPFormat = "ldp_vp"
)

// Model to convey the details about the Credentials the Client wants to obtain.
type AuthorizationDetails struct {
	// String representing a format in which the Credential is requested to be issued. Valid values defined by OIDC4VC are jwt_vc_json-ld and ldp_vc. Issuer can refuse the authorization request if the given credential type and format combo is not supported.
	Format *string `json:"format,omitempty"`

	// An array of strings that allows a client to specify the location of the resource server(s) allowing the Authorization Server to mint audience restricted access tokens.
	Locations *[]string `json:"locations,omitempty"`

	// String that determines the authorization details type. MUST be set to "openid_credential" for OIDC4VC.
	Type string `json:"type"`

	// String array denoting the types of the requested Credential.
	Types []string `json:"types"`
}

// DID method of the DID to be used for signing.
type DIDMethod string

// Model for KMS configuration.
type KMSConfig struct {
	// Prefix of database used by local kms.
	DbPrefix *string `json:"dbPrefix,omitempty"`

	// Type of database used by local kms.
	DbType *string `json:"dbType,omitempty"`

	// URL to database used by local kms.
	DbURL *string `json:"dbURL,omitempty"`

	// KMS endpoint.
	Endpoint *string `json:"endpoint,omitempty"`

	// Path to secret lock used by local kms.
	SecretLockKeyPath *string `json:"secretLockKeyPath,omitempty"`

	// Type of kms used to create and store DID keys.
	Type KMSConfigType `json:"type"`
}

// Type of kms used to create and store DID keys.
type KMSConfigType string

// Supported VC formats.
type VCFormat string

// Supported VP formats.
type VPFormat string

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5xW0W/buA/+Vwg9bYCbFr/9nvLWc3ZA0PYWLGv2cBsKWaJjNbLkk+i4uaH/+4Fy3DSN",
	"c9jupagp8SP58ROZH0L5uvEOHUUx/SGiqrCW6d/rliofzN+SjHczJGlssmuMKpiGrWIq7rxGC+RBebfF",
	"HVCFoPvLIAvfUrLkATU6MtLG/tsadASddBTZ2RckjZuITDTBNxjIYIpV+lBLOo26pGDcGgI2ASMDuzVI",
	"6G+DcdBVRlVvIoOJEPCvFiOh5qAFgomxRT2BlbRGw1baFiNoLI1DDcUOPs1n+f9XOciA8NjRw1Y9PEbv",
	"LqwG6TRY3Txs1QTmDBNASQcByzZiCi1fEziEBlOmw7XZogN1yI52DSbQfRXK14XnnJ0niG3T+EComSK+",
	"KaYiJg7EcyasVynGSHuuHcgQ5A58Cb0DN0ASSGt9F0GC6ltBHmKDypR9CwdI9uPvgNG3QSFEDFsM7+L7",
	"HoGJ5/MjscAyXWLM2jgC2WqDTiUUCkYx/1IpjNz7DbrIVRnCOhVwUt7ekOo4fJ+RRCpOI2GojcM40ohB",
	"nQwzgbv75RdWQsTEwTfhG3RGPxw6801wSwYpjDaADfFsSn0DNDpPA13J4cDtIMqDWH+FkedMMIQJqMX0",
	"z/5wyOl7JsiQ5duj7/kFyxePqIjBZ/PZHVLl9WlBs/kM6nQ25M6W/im1EZN2IZq1M27NFaBra07Jh0Jk",
	"okP+u8FdyuptTTd3y9y70qzPzRjGvrlb8qApzboNqY7TkaGLRcDSPJ3C9HbOXEuShYz7pItdkruFTR1H",
	"26uLL6OSY+t/grv/fHuKdv/5lqn8RTB0uvHGjcxI5mo4HXWNqALSrVebG9wtJFUjlEmq0mhIVzmVzU/m",
	"Rf/K2KaOPQ4vjoCS+tEXyYdeUxvcxdcKSsFeNCS7OKKhsXfwSv8Hgb0VfSaeLkiuI3uljRDE9+dMrPLf",
	"z62fYRzDKt/P66Nsj1eFyES/KF5n8wI+wt1q8ROBF2cDN0PA5ijg4lxAJs640o+0K7SRfrNewSpfDkvo",
	"9dJitiRPdu7eFoMpzX5vtJFn3dcPOazyi+vFHKT1bg2doQo+NejmM96rTfDkle8fd0/9ZYLBAMYRBqkS",
	"WnL7Kq3FJGVrFLqY9OVknYZbI1WFF/+bXIlMtMGKqaiImji9vOy6biLT8cSH9eXeN17ezvOPfyw/ss+E",
	"ntLoG7jKfV17t5/SnNoqlSYLe/SDgtecUQjvVvnyvcjEFkPsibuacCbPWdomsjFiKj5MrlJyjaQqiqlr",
	"rX3+JwAA///XafVnfgkAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
