package apple

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var audienceKey string
var token string
var verifier *AppleVerifier

func TestMain(m *testing.M) {
	audienceKey = "clientId"
	signed, serverMock, err := GetTestEnvironment(audienceKey)

	token = signed

	if err != nil {
		panic(err)
	}

	verifier = &AppleVerifier{
		JwksAppleURI: serverMock.URL,
	}
	os.Exit(m.Run())
}

func TestVerifyAppleIdToken(t *testing.T) {
	response, err := verifier.VerifyAppleIdToken(token, audienceKey)

	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	a.Equal(response.Id, "abc")
	a.Equal(response.Email, "test@gmail.com")

	response, err = verifier.VerifyAppleIdToken("test_fail", audienceKey)

	a.Nil(response)
	a.NotNil(err)
	a.True(strings.Contains(err.Error(), "invalid compact serialization format: invalid number of segments"))

	response, err = verifier.VerifyAppleIdToken(token, "audience_fail")

	a.Nil(response)
	a.NotNil(err)
	a.True(strings.Contains(err.Error(), "The audience parameter does not include this client"))
}
