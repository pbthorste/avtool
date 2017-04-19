package avtool

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"io/ioutil"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EncryptFile(filename, password string) (result string, err error) {
	data, err := ioutil.ReadFile(filename)
	check(err)
	result, err = Encrypt(string(data), password)
	return
}

// see https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L710
func Encrypt(body, password string) (result string, err error) {
	salt, err := GenerateRandomBytes(32)
	check(err)
	//salt_64 := "2262970e2309d5da757af6c473b0ed3034209cc0d48a3cc3d648c0b174c22fde"
	//salt,_ = hex.DecodeString(salt_64)
	key1, key2, iv := genKeyInitctr(password, salt)
	ciphertext := createCipherText(body, key1, iv)
	combined := combineParts(ciphertext,key2,salt)
	vaultText := hex.EncodeToString([]byte(combined))
	result = formatOutput(vaultText)
	return
}

func createCipherText(body string, key1,iv []byte) []byte {
	bs := aes.BlockSize
	padding := (bs - len(body) % bs)
	if padding == 0 {
		padding = bs
	}
	padChar := rune(padding)
	padArray := make([]byte, padding)
	for i := range padArray {
		padArray[i] = byte(padChar)
	}
	plaintext := []byte(body)
	plaintext = append(plaintext, padArray...)

	aes_cipher, err := aes.NewCipher(key1)
	check(err)
	ciphertext := make([]byte, len(plaintext))

	aesBlock := cipher.NewCTR(aes_cipher, iv)
	aesBlock.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func combineParts(ciphertext, key2, salt []byte) string {
	hmacEncrypt := hmac.New(sha256.New, key2)
	hmacEncrypt.Write(ciphertext)
	hexSalt := hex.EncodeToString(salt)
	hexHmac := hmacEncrypt.Sum(nil)
	hexCipher := hex.EncodeToString(ciphertext)
	combined := string(hexSalt) + "\n" + hex.EncodeToString([]byte(hexHmac)) + "\n" + string(hexCipher)
	return combined
}

// https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L268
func formatOutput(vaultText string) string {
	heading := "$ANSIBLE_VAULT"
	version := "1.1"
	cipherName := "AES256"

	headerElements := make([]string, 3)
	headerElements[0] = heading
	headerElements[1] = version
	headerElements[2] = cipherName
	header := strings.Join(headerElements, ";")

	elements := make([]string, 1)
	elements[0] = header
	for i := 0; i < len(vaultText); i += 80 {
		end := i + 80
		if end > len(vaultText) {
			end = len(vaultText)
		}
		elements = append(elements, vaultText[i:end])
	}
	elements = append(elements, "")

	whole := strings.Join(elements, "\n")
	return whole
}