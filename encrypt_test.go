package avtool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	password := "asdf"
	body := "secret"
	var encrypted string
	var err error
	encrypted, err = Encrypt(body, password)
	assert.NoError(t, err)

	var result string
	result, err = Decrypt(encrypted, password)
	assert.NoError(t, err)
	assert.Equal(t, body, result)
}
