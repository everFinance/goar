package goar

import (
	"testing"

	"github.com/daqiancode/goar/types"
	"github.com/daqiancode/goar/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewSigner(t *testing.T) {
	signer, err := NewSignerFromPath("./example/testKey.json")
	assert.NoError(t, err)
	tags := []types.Tag{
		{Name: "GOAR", Value: "sendAR"},
	}
	tx := &types.Transaction{
		Format:    2,
		ID:        "EyPNVxI-zv1WGjiMmeb20VimIGjfvnkFQtgVnlSMX1o",
		LastTx:    "0cJJGyXdeIJ-azX3I-jJ3XJEIYti-_qjHZhaRsXQxi1D3wXQVv6Px5WQhj_j1W8O",
		Owner:     "zu2HnkLgrxsqqtdOzsSx6qezj8ujqYzJCaz9r6XosF5oHRdyqyxZFmuFeQaWS9_8YP73XFDUn2FJRNMzd6lclbIlgfiorufLfa6ez_4Y3yMVfYC56anQs0f62TwfcMuQCfTnPyPDds3evODLLRyuwVqtfe6jLZZ4LGbcgXySi0Hn184VxOTkDzvYz3zWz0tqpXX2T8qNeELNHH-j6pXMH7LNzU2akZhZUeuxsstxyedaZOP_EUwx-BtooWXWcVNcjM_elMitgZelU60J012AdgUXuj9kULxT78_BSFhiD4vGn5ZN8V1lpB410wGqQNuiGZGjUrorY9zmm3IDIhywxJrzcaTshxPBb8-gw-38mCOTP0yp4URs4lae75fxrD-m4kXzKl27XZ10K21U47W0hdJI1s934MPpuLudyvlXbxkvm19LzZj9wkmqf08RqghMS5X_LVgAT0bs6UMnLoJcQFckBWpOVDVY-spkPqxBhEomhDD3C9xvXz5KqR8Q9kgjX6nSXyn5svAsGGviHBzU27FtR9L01q6QCeCUEehCG_n0eVM4IsHIySh7ZZk1sFjBvgEEXiEpxA-av5ob2nfcR93cWILlhmPs4wdfzi8Q52jM2hShZXTYGvAOEdXdZmBatoVgm9Sca6U-pffmFqYl2BM-Bq_OjBpDnGvvQBV_990",
		Tags:      utils.TagsEncode(tags),
		Target:    "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE",
		Quantity:  "1000000000",
		Data:      "",
		DataSize:  "0",
		DataRoot:  "",
		Reward:    "747288",
		Signature: "",
	}
	err = signer.SignTx(tx)
	assert.NoError(t, err)

	err = utils.VerifyTransaction(*tx)
	assert.NoError(t, err)
}
