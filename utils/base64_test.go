package utils

import (
	"encoding/json"
	"github.com/everFinance/goar/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	assert.Equal(t, "QXBwLU5hbWU", Base64Encode([]byte("App-Name")))

	res, err := Base64Decode("QXBwLU5hbWU")
	assert.NoError(t, err)
	assert.Equal(t, "App-Name", string(res))
}

func TestBase64Decode(t *testing.T) {
	tagsBy := `[{"name": "Data-Protocol", "value": "ao"}, {"name": "Variant", "value": "ao.TN.1"}, {"name": "Type", "value": "Message"}, {"name": "SDK", "value": "goar"}, {"name": "Action", "value": "Mint-Token"}, {"name": "TokenId", "value": "Wb1N5IIQzyNv0tTHmGZL54D_6mbSi4tRKhYFrp5b0xQ"}, {"name": "TxHash", "value": "7xWzbEahU3NMa14T1rLgn-Jl5xPKsnJ2sCalgRernBQ"}, {"name": "Recipient", "value": "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE"}, {"name": "Quantity", "value": "5000000000"}]`
	tags := make([]types.Tag, 0)
	err := json.Unmarshal([]byte(tagsBy), &tags)
	assert.NoError(t, err)
	item := types.BundleItem{
		SignatureType: 1,
		Signature:     "oXpKZdJ2sevzM_43GOnRh-5El7dLl95LvzPFXwFnaGHBv49IF9YPav2tSSA4FRL8rsEE1wKRxs4nUrNSVOUC99gAow5OdpD3UzL6D0S-tBSVDGk5h_v_NyKUbhD7CLQ33iVKqJhGRrIaofD9TqA98xeiTm2iiJnCqu7NbP_8XsSwdWLlgC18Js7hRUzPDdmK3QysR-uVCIh_icSEHqTTNSj5XmWcjXi3xr7GOkzOdtOLKf5onWKynkMjGvjwhCoNRQjlLkvZ0MrfHImVxeS_UdMOfjxzINCa659GD3cBw0E0lYvkgAOuT3X7ZUuC-eDLqYNCdZg9Lg1tJki2Kg0L90OhwXul7Aa7tiptwVbEWKr4a4GTlxtsgCA2Pw-0FcP9fbtAnVUxsEvjKrb9BsnTQxknAGJjEIOgjU3mpoOsYaUmxnTkEPiERI5T7BGoACnLIVKAHDbp9nKdY5Hu7JxqAKsHPIUZfWAa1QwJ8_ubFZR_IM35TXx-Bal5Q8dLm6-2Ec6ylK4NnA0oyONs09hAbSfBo2XMS2NoIileG9TxtnjoW6_kPsZirCMIipiGtbANGaZTlm29YE3SQuEBLSvx4HYF81XL2tJfa6BuHljvNr0gXSodS47swWHlBMYIxM98deggKeDgNGUJekrFZj9A1smQ0NR-34mP1DBehXHfy60",
		Owner:         "r1LsVhfzPE0wVgavmVdxpzH0lESO7pEJS1pmhlBAqNaiDykPm-IDXXAfN51WnR1pX5PDEz6ivB3o212D0-_JUoKzcSzjHRLkPKLsPwr0ab4_TC8bdzensbCSDZ9NbhprOAkhlOA8Ii6WtBI1CwV7Nk5lS9aKOW7keKRlBaLQELXCr5HluIwmBDT_oMKISv5ixorUMuNrD22uUlbqDkbwcsOD1TOyM5n9PhV7aoYDql5cXETay6WsJ9q9oQAGrh6HUbXkJK3TFOQx4Hb8JLT45Yi-W_9LQz4cvloTlhOoPl42op6hdaflhNmNhSzTgCRhJyU67yq6SRJidIBio_z8-dQOuMzu8rqxDp31DvZjPRFcMV77gBgJS8Sc7mYKoptHZMTGCAWu6eQB2OwgRN_V3DAZAGUubWdX9A1J6jv4cB9lX7Bcn-nImxVQnVxgwVuh74K9YnbieO_Sr3pk5svhtq6iAEkgTh01aRNwDimFR1naKgQigiyMz4EodgAACsQqxFrnprIQ5TRrgDtYErtchz2YjWaCcfikx5Ex98pId2SkexZCsTtfibEHJpnM2X59ZSwOM_BOgRhRKk6c49bPCyuvnIxraA2Kbf-rrMlNl6yTx4xcUJ4Ll2g1UIYmjh_gzoRL4JnJLAPkiHIfg-5r-b_gu_UsGwOorc3q24sN_sM",
		Target:        "lp-ihMlyEo5NSoW5Rkd7kU1PoTqCIYAGfDrn2G4ihiI",
		Anchor:        "",
		Tags:          tags,
		Data:          "MTExMQ",
		Id:            "xFj8tieuPT9XF-tY8FUr5ZF-9yoj-4vp2nVxA-Q6jZ8",
	}
	// err = utils.VerifyBundleItem(item)
	// assert.NoError(t, err)
	itemBinary, err := GenerateItemBinary(&item)
	assert.NoError(t, err)
	i, err := DecodeBundleItem(itemBinary)
	assert.NoError(t, err)
	t.Log(i.Id)
}
