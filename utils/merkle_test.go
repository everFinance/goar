package utils

// func TestGenerateChunks(t *testing.T) {
// 	data, _ := Base64Decode("NzcyNg")
// 	assert.Equal(t, "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U", Base64Encode(GenerateChunks(data).DataRoot))
//
// 	data, err := ioutil.ReadFile("./testfile/1mb.bin")
// 	assert.NoError(t, err)
// 	chunks := GenerateChunks(data)
// 	assert.Equal(t, "o1tTTjbC7hIZN6KbUUYjlkQoDl2k8VXNuBDcGIs52Hc", Base64Encode(chunks.DataRoot))
// }

// func TestChunkData(t *testing.T) {
// 	data, _ := Base64Decode("NzcyNg")

// 	assert.Equal(t, "55w6-CA_Um7muHLnJvBlUUKTjpa35cPLv1PCIPQs6M8", Base64Encode(chunkData(data)[0].DataHash))
// }

// func TestGenerateLeaves(t *testing.T) {
// 	data, _ := Base64Decode("NzcyNg")

// 	assert.Equal(t, "z3rQGxyiqdQuOh2dxDst176oOKmW3S9MwQNTEh4DK1U", (Base64Encode(generateLeaves(chunkData(data))[0].ID)))
// }
