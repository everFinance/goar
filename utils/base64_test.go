package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	assert.Equal(t, "QXBwLU5hbWU", Base64Encode([]byte("App-Name")))

	res, err := Base64Decode("QXBwLU5hbWU")
	assert.NoError(t, err)
	assert.Equal(t, "App-Name", string(res))
}
