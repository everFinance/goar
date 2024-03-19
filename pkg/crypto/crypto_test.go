package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepHash0(t *testing.T) {
	data := []byte{1, 2, 3}
	r := DeepHash(data)
	h, err := hex.DecodeString("41300af79285f856e833164518c7ec4974f5869ec77ca3458113fe6c587680d050f9f6864fd77f9eb62bd4e2faea9ae8")
	assert.NoError(t, err)
	assert.ElementsMatch(t, r[:], h)
}

func TestDeepHash1(t *testing.T) {
	data := []byte{}
	r := DeepHash(data)
	h, err := hex.DecodeString("fbf00cc444f5fea9dc3bedf62a13fba8ae87e7445fc910567a23bec4eb82fadb1143c433069314d8362983dc3c2e4a38")
	assert.NoError(t, err)
	assert.ElementsMatch(t, r[:], h)
}

func TestDeepHash2(t *testing.T) {
	data := [][]byte{{1, 2, 3}, {4, 5, 6}}
	r := DeepHash(data)
	h, err := hex.DecodeString("4dacdcc81acd09f38c77a07a2a7ae81f77c61e6b97ee5cc7b92f3a7f258e8d5ba69d14d7d66070797b083873717c9896")
	assert.NoError(t, err)
	assert.ElementsMatch(t, r[:], h)
}
