package utils

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	rightKey, err := GenerateRsaKey(4096)
	assert.NoError(t, err)
	wrongKey, err := GenerateRsaKey(4096)
	assert.NoError(t, err)
	msg := []byte("123")

	sig, err := Sign(msg[:], rightKey)
	assert.NoError(t, err)
	assert.NoError(t, Verify(msg, &rightKey.PublicKey, sig))
	assert.Error(t, Verify(msg, &wrongKey.PublicKey, sig))
}

func TestVerify(t *testing.T) {
	pub, err := OwnerToPubKey("5oZAXKfBfZ32JIKSIUksXO6K27T9gnbJCnndpGTvKf1HFUh5uh1y38Zcxwn8D5hhOfgWoonVmWcuxvNBw3LaW4q7NUEu72ukg0KpaipGOqVvzChsP4MuD79mXgQd_N-J18MT51mO3uGwiMLJUohkWT_XkhSDSsPyJoS9xhoRZN9OseDB-4cIkevfIyewrBpmX6wSZ7yBf6lml5btAi57Aha-DvhUSEKV0mtsx7C81jf1Aw8B639zPomT7eYYRWepkGwfQ7_JZMJ7ddAC0hEUpTDcmBe20k5i9XzlmqB1bFIu236BrswD-xVTGvWidrEVTUafL8jOpSML875iOKGgC_y6aCgwMgTxemneWNgsUjkSOvi24fPxpGoFFzcRPWnQ1Ru5Wuw6RqO5RlVpvEvztXEE5IprXWdP286lgLTzcduagT4dHouhDrMaIP68lq17u2p3RaU2Awgn4kOQA-p8_iVESj9hPlLwhbQ48I46Vh2Eq6SgJffvr-mYqOY6jJB3AttsftMk8zxhmMl20nDT0RbKNGkKicezjbxXi9Pp2j4s_H_ynvrNvFWj_JS3wHITuX_vKcYFxgmExcpTa4FQsuwBeZxI2Ls9g4kkIM-CcXDwfp-aQ4JIpfWLWt5MKKx5ouhOXBQa5GvNs3YR3imhUTv1FOpm9n_jscWxwZJMTls")
	assert.NoError(t, err)
	data := "0xeaa110e6abf9423481fa6d5a72a3d66118e05ea4759fd4a6d6344de83e18dd52"
	dataBy, err := hexutil.Decode(data)
	assert.NoError(t, err)

	sig, err := Base64Decode("5fS8LHBda15AMKBH8QQUuDCk23RJdYhWFq_rNoyPusqEdFGeiawSGwypcIhTeP9cOXKVZjf9LT9rPcEpJm49EdDb5FcqV7xCBSVV0CpLwJKCoYUMk1tDWFBN2wUbRDUe2PI7ycUKrurXrQeT1WfjEc-Nh7nGCR9kKv9MJfBoo76GhbC_CKtrmGFUeok3t9uoVjEOGC-iXtvMjh6G3hIV2JFIiSUdBS3qMIRa7FZiJ2ZkJaS-4vTwh70esP6KGkLYW5nhJWZ9Y0jIndLAxnFz-vcJJDmaq4fPUcWDO49sH0u8ZZe6mAs3Kdj2B7rGmzozYi4-j_nivnbWUGM7yCw6w1bops8Cxwr0g9l16lBPkAHBgxE1CGS0_KbsJ_GkbNlEU-6znEuFJpDubhC12QNmEeH8TIIwJqz8aw8G6zOyTP2wXyP6TIoFAemMIbyzAoV99QQ-2RWjnoLG0srgfacfJlhHcSSAaXFWB1Mr-bCCMjetLLcRrK5wwa-h7NAT6feZIUni0bq-SGc9oIVNDFvOsNJnHVKx644wpmdmNkZSINUccuJDwzNYKrbqqDIIpe_vzaAgLbS940LS5UmYtbkZkMvLQBFzbEW9C1A4rpJoDy-nibD0cqP1gg0AIdueElKrA1Szuj5Kh3_uOz-Fh81SYNLye4pX1bYWuV_wPt7qd3E")
	assert.NoError(t, err)
	hashed := sha256.Sum256(dataBy)
	err = Verify(hashed[:], pub, sig)
	assert.NoError(t, err)
}
