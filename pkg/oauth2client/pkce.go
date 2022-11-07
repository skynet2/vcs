package oauth2client

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
)

const (
	defaultLength = 56
	defaultMethod = "S256"
)

func (c *Client) GeneratePKCE() (verifier string, challenge string, method string, err error) {
	b := make([]byte, defaultLength)
	_, err = rand.Read(b)

	if err != nil {
		return "", "", "", err
	}

	return c.GeneratePKCEFromBytes(b)
}

func (c *Client) GeneratePKCEFromBytes(b []byte) (verifier string, challenge string, method string, err error) {
	verifier = c.encode(b)

	h := sha256.New()
	h.Write([]byte(verifier))
	challenge = c.encode(h.Sum(nil))

	return verifier, challenge, defaultMethod, nil
}

func (c *Client) encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.Replace(base64.StdEncoding.EncodeToString(data), "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)

	return encoded
}
