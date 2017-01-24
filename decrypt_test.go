package avtool

import (
	"testing"
)

func Test1(t *testing.T) {
	password := "asdf"
	filename := "./testdata/test1/secrets.yaml"
	Decrypt(filename, password)
}
