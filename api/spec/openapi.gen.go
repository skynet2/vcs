// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	externalRef0 "github.com/trustbloc/vcs/pkg/restapi/v1/common"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
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

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./common.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
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
