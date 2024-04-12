package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGenerateChunks(t *testing.T) {
// 	data, _ := Base64Decode("NzcyNg")
// 	assert.Equal(t, "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U", Base64Encode(GenerateChunks(data).DataRoot))
//
// 	data, err := os.ReadFile("./testfile/1mb.bin")
// 	assert.NoError(t, err)
// 	chunks := GenerateChunks(data)
// 	assert.Equal(t, "o1tTTjbC7hIZN6KbUUYjlkQoDl2k8VXNuBDcGIs52Hc", Base64Encode(chunks.DataRoot))
// }

func TestChunkData(t *testing.T) {
	data, _ := Base64Decode("NzcyNg")

	assert.Equal(t, "55w6-CA_Um7muHLnJvBlUUKTjpa35cPLv1PCIPQs6M8", Base64Encode(chunkData(data)[0].DataHash))
}

func TestGenerateLeaves(t *testing.T) {
	data, _ := Base64Decode("NzcyNg")

	assert.Equal(t, "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U", (Base64Encode(generateLeaves(chunkData(data))[0].ID)))
}

func TestChunkStream(t *testing.T) {
	data, err := os.ReadFile("img.jpeg")
	assert.NoError(t, err)
	dataReader, err := os.Open("img.jpeg")
	assert.NoError(t, err)
	defer dataReader.Close()
	chunks01 := chunkData(data)
	chunks02, err := chunkStreamData(dataReader)
	assert.NoError(t, err)
	assert.Equal(t, len(chunks01), len(chunks02))
	for i := 0; i < len(chunks01); i++ {
		assert.Equal(t, chunks01[i].MaxByteRange, chunks02[i].MaxByteRange)
	}
}

func TestDecodeEmptyString(t *testing.T) {
	data := ""
	by, err := Base64Decode(data)
	assert.NoError(t, err)
	t.Log(by)
}
