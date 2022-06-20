package goar

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// arKeyJs := `{"kty":"RSA","d":"QMU14EBgFCgTaCR6wMzmEdrNAAJ5QhTIWBktJy7Yo3gdC9oAY29awpou9n3Wl2LMacS6-5ihMrVlzeizuyuXQry6Fv1fg5dHCwVZeWWsOErrxu37T2su0BNzI78N2Bxrtx3wGmvnc1vLX4tHAEwkLFxUtzwXteERpRgTwm75OKeE_fiUlCZInBQ78j02-7zjzACruXhYTPXYbKv33eMgeQkeRL0mXkfUyfY0hZ4frAjxDjfxZubMHkm3aXHNHNHr0pWGdct06uPZvg-NVcs3cUV1T2CoyZC1qnf32cOCbxQ92Sm8XfTtkHCxRjZQHc3ACEpQuDiX4dMiewaSGUdW7uGearbZHAUV-AKkfb2Sc0W6Z7aD0dE_uDPxYn-akSiYboNbwX6uY22KZ03FsJJaXbr7W0byQrrTn5eAWso3XlTiNjDhk1GOPmX5YhKMk2Is7XcRCm_H9crhTUkCsTAKIbeJ3YivN_RGO2S8csXA6oc6xjl_jRvGjiGt4Cj3u-WV4KLvkE7x7aEl6eTEvwisbNbOHBMugPPXtYEW_UseRRxWKQIgBcToe18UqfPfH8EhNqOK6YJ5a428JaMjm6yfx79A6hmKxZhtPywW0rrC0B925O4xKFVnq7qzvf0Ik_PD4Snc8fV_p6pzufF4gnJRcnsN9tyktx5eF-oSowgNqSU","n":"sQniDmUoDgJ-ZTy1LeVx9mvjdWiCIs6ftuac7tIkDG-ZtwYe1DCf6wpnw-tBc6gFLHXLqlMNq4fed5YwJsVxGpiU6xl9uxQRk83ztUuZWhfaWgZooJCA7SMATRFjZ1gzM5OOnXv9TTJT-jhWTrjoVa1D-aUrySEKODNCGyt7uwj_MZxpaP7CR91SWyEJZSmJ8cMTH7NQ6a0ZFecLDI5Tr_C0CEdwtb6KacEuqrWK2tGkUbXVjsNQDml8C8HuMpOrGkpNDwAKJZHMv_wFkPzMqPX7CJK-vAlCZA9bVoeI35pSt9-EDyzOStOSBCjac1oO0elaujKhhktRbzYMO-8h_lIKLWyjD1vrHehk6-MTe8w0DxAVV1tku_w18T62gr6uq7FiDyaY2zxboYPW7YQaZNVh-dFwnPeDDVa-5tF6rZbjGPj2BmwW19RgSbVRfEJlw93KUU4xnUmUVhb2CSy_jz3QdUDNniRS7gQ9sSQruJeh-4aDN4DAxcXRlhAJ9mN_YjDnkXt4c_A02TO-CIv4Bnz8uXQ3FcUNTRlIrAwNAfrsKN6nB41l7HlYi_Ttllh2W5B_PFcVTyqrxhkb1kFfPflf9RwNbDw9MMN4sqcO2s4uTusqsOVI01FFjzlewroHRFXKKnwbaftE7aSAAS7fFzOhP4juczZBxkRg577NR0E","e":"AQAB"}`
	// signer, err := NewSigner([]byte(arKeyJs))
	// assert.NoError(t, err)
	//
	// arseedUrl := "https://seed-dev.everpay.io"
	// itemSdk, err := NewItemSdk(signer, arseedUrl)
	// assert.NoError(t, err)
	//
	// data := []byte("123456")
	// item01, err := itemSdk.CreateAndSignItem(data, "", "", nil)
	// assert.NoError(t, err)
	//
	// err = utils.VerifyBundleItem(item01)
	// assert.NoError(t, err)
	// resp, err := itemSdk.SubmitItem(item01, "USDT")
	// assert.NoError(t, err)
	// t.Log(*resp)
}
