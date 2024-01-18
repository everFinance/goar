package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	// rpId
	LocalhostRpId = "localhost"
	EverpayRpId   = "everpay.io"

	// rp origins
	LocalhostOrg      = "http://localhost:8080"
	EverpayOrg        = "https://app.everpay.io"
	EverpayDevOrg     = "https://app-dev.everpay.io"
	BetaDevEverpayOrg = "https://beta-dev.everpay.io"
	BetaEverpayOrg    = "https://beta.everpay.io"
)

type Authn struct {
	Id                string `json:"id"`
	RawId             string `json:"rawId"`
	ClientDataJSON    string `json:"clientDataJSON"`
	AuthenticatorData string `json:"authenticatorData"`
	Signature         string `json:"signature"`
	UserHandle        string `json:"userHandle"`
}

type User struct {
	Id     string
	Name   string
	Public webauthn.Credential // key: publicType，val: credential
}

func NewUser(id, name string, public webauthn.Credential) *User {
	return &User{
		Id:     id,
		Name:   name,
		Public: public,
	}
}

func (u *User) WebAuthnID() []byte {
	id, _ := Base64Decode(u.Id)
	return id
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.Name
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{u.Public}
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func GetWebAuthn(rpIdHash []byte) (*webauthn.WebAuthn, error) {
	localhostRpIdHash := sha256.Sum256([]byte(LocalhostRpId))
	everpayioRpIdHash := sha256.Sum256([]byte(EverpayRpId))
	switch string(rpIdHash) {
	case string(localhostRpIdHash[:]):
		return webAuthn(LocalhostRpId)
	case string(everpayioRpIdHash[:]):
		return webAuthn(EverpayRpId)
	default:
		return nil, errors.New("err_rp_id_not_exist")
	}
}

func webAuthn(rpId string) (*webauthn.WebAuthn, error) {
	return webauthn.New(&webauthn.Config{
		RPID:          rpId,
		RPDisplayName: "everpay",
		RPOrigins: []string{EverpayOrg, EverpayDevOrg, BetaDevEverpayOrg, BetaEverpayOrg,
			LocalhostOrg},
		AttestationPreference: "",
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationRequired,
		},
		Debug:                false,
		EncodeUserIDAsString: false,
		Timeouts:             webauthn.TimeoutsConfig{},
	})
}
func VerifyFidoAuthnSig(sig, hexHash string, accName, userId string, public webauthn.Credential) (*webauthn.Credential, error) {
	sigBy, err := decodeBase64(sig)
	if err != nil {
		return nil, err
	}
	authn := Authn{}
	if err = json.Unmarshal(sigBy, &authn); err != nil {
		return nil, err
	}
	// parse car
	rawId, err := decodeBase64(authn.RawId)
	if err != nil {
		return nil, err
	}
	ClientDataJSON, err := decodeBase64(authn.ClientDataJSON)
	if err != nil {
		return nil, err
	}
	AuthenticatorData, err := decodeBase64(authn.AuthenticatorData)
	if err != nil {
		return nil, err
	}
	Signature, err := decodeBase64(authn.Signature)
	if err != nil {
		return nil, err
	}
	UserHandle, err := decodeBase64(authn.UserHandle)
	if err != nil {
		return nil, err
	}

	car := protocol.CredentialAssertionResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				ID:   authn.Id,
				Type: "public-key",
			},
			RawID:                   rawId,
			ClientExtensionResults:  nil,
			AuthenticatorAttachment: "platform",
		},
		AssertionResponse: protocol.AuthenticatorAssertionResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: ClientDataJSON,
			},
			AuthenticatorData: AuthenticatorData,
			Signature:         Signature,
			UserHandle:        UserHandle,
		},
	}
	pca, err := car.Parse()
	if err != nil {
		return nil, err
	}

	// new user
	user := NewUser(userId, accName, public)
	session := webauthn.SessionData{
		Challenge:            Base64Encode([]byte(hexHash)),
		UserID:               user.WebAuthnID(),
		UserDisplayName:      user.WebAuthnDisplayName(),
		AllowedCredentialIDs: nil,
		Expires:              time.Time{},
		UserVerification:     protocol.VerificationRequired,
		Extensions:           nil,
	}

	webAuthn, err := GetWebAuthn(pca.Response.AuthenticatorData.RPIDHash)
	if err != nil {
		return nil, err
	}
	credential, err := webAuthn.ValidateLogin(user, session, pca)
	if err != nil {
		return nil, err
	}
	return credential, nil
}

func decodeBase64(s string) (protocol.URLEncodedBase64, error) {
	// StdEncoding: the standard base64 encoded character set defined by RFC 4648, with the result padded with = so that the number of bytes is a multiple of 4
	// URLEncoding: another base64 encoded character set defined by RFC 4648, replacing '+' and '/' with '-' and '_'.
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")

	bs := &protocol.URLEncodedBase64{}
	err := bs.UnmarshalJSON([]byte(s))
	if err != nil {
		return nil, err
	}
	return *bs, nil
}

func GenUserId(eid string, chainID int) string {
	data := []byte(strings.ToLower(eid) + strconv.Itoa(chainID))
	hash := sha256.Sum256(data)
	userId := Base64Encode(hash[:10])
	return userId
}
