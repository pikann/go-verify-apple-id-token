package apple

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AppleVerifier struct {
	JwksAppleURI string
}

func NewAppleVerifier() *AppleVerifier {
	return &AppleVerifier{
		JwksAppleURI: "https://appleid.apple.com/auth/keys",
	}
}

func (v *AppleVerifier) GetApplePublicKey(kid string) (*jwk.Key, error) {
	resp, err := http.Get(v.JwksAppleURI)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	set, err := jwk.ParseString(string(bodyText))

	if err != nil {
		return nil, err
	}

	publicKey, _ := set.LookupKeyID(kid)

	return &publicKey, nil
}

func (v *AppleVerifier) VerifyAppleIdToken(token string, clientId string) (*VerifyTokenResponse, error) {
	headers, err := jws.Parse([]byte(token))

	if err != nil {
		return nil, err
	}

	kid := headers.Signatures()[0].ProtectedHeaders().KeyID()

	publicKey, err := v.GetApplePublicKey(kid)

	if err != nil {
		return nil, err
	}

	if publicKey == nil {
		return nil, errors.New("Key ID invalid")
	}

	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, *publicKey))

	if err != nil {
		return nil, err
	}

	tokenClaims := verifiedToken.PrivateClaims()

	clientIdVerified := false

	for _, aud := range verifiedToken.Audience() {
		if aud == clientId {
			clientIdVerified = true
		}
	}

	if clientIdVerified {
		return &VerifyTokenResponse{
			Id:    verifiedToken.Subject(),
			Email: tokenClaims["email"].(string),
		}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("The audience parameter does not include this client - is: %v | expected: %v", clientId, verifiedToken.Audience()))
	}
}
