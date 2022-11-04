package oauth2client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPKCE(t *testing.T) {
	p := NewPKCEService()
	verifier, challenge, err := p.Generate()
	assert.NoError(t, err)
	assert.NotEmpty(t, verifier)
	assert.NotEmpty(t, challenge)
}
