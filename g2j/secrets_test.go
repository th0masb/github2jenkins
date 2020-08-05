package g2j

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessfulLoadJsonSecrets(t *testing.T) {
	// assemble
	secretsLocation := "/path/to/some/secrets.json"
	data := `
	{
		"first": "second",
		"third": "fourth"
	}
	`
	expectedSecrets := Secrets{
		"first": "second",
		"third": "fourth",
	}

	// act
	secrets, err := interpretSecrets(secretsLocation, []byte(data))

	// assert
	assert.Nil(t, err, "Expected success")
	assert.Equal(t, expectedSecrets, secrets)
}

func TestFailedLoadJsonSecrets(t *testing.T) {
	// assemble
	secretsLocation := "/path/to/some/secrets.json"
	data := `
	{
		"first": "second,
		"third": "fourth"
	}
	`
	expectedSecrets := Secrets{}

	// act
	secrets, err := interpretSecrets(secretsLocation, []byte(data))

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, expectedSecrets, secrets)
}

func TestUnrecognisedExtension(t *testing.T) {
	// assemble
	secretsLocation := "/path/to/some/secrets.jso"
	data := `
	{
		"first": "second",
		"third": "fourth"
	}
	`
	expectedSecrets := Secrets{}

	// act
	secrets, err := interpretSecrets(secretsLocation, []byte(data))

	// assert
	assert.NotNil(t, err, "Expected success")
	assert.Equal(t, expectedSecrets, secrets)
}
