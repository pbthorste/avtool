package avtool

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt1(t *testing.T) {
	password := "asdf"
	body := "secret"
	encrypted,_ := Encrypt(body, password)

	result, err := Decrypt(encrypted, password)
	if(err != nil) {
		fmt.Println(err.Error())
	}
	assert.Equal(t, body, result)
}

