package apple

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func GetTestEnvironment(audienceKey string) (string, *httptest.Server, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return "", nil, err
	}

	key, err := jwk.FromRaw(raw)

	if err != nil {
		return "", nil, err
	}

	err = key.Set(jwk.KeyIDKey, "kidTest")
	if err != nil {
		return "", nil, err
	}

	token := jwt.New()
	err = token.Set(jwt.SubjectKey, "abc")
	if err != nil {
		return "", nil, err
	}
	err = token.Set(jwt.AudienceKey, audienceKey)
	if err != nil {
		return "", nil, err
	}
	err = token.Set("email", "test@gmail.com")
	if err != nil {
		return "", nil, err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return "", nil, err
	}

	publicKey, err := key.PublicKey()
	if err != nil {
		return "", nil, err
	}

	err = publicKey.Set(jwk.KeyIDKey, "kidTest")
	if err != nil {
		return "", nil, err
	}
	buf, err := json.MarshalIndent(map[string][]jwk.Key{"keys": {publicKey}}, "", "  ")
	if err != nil {
		return "", nil, err
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(buf)
	}))

	return string(signed), testServer, nil
}
