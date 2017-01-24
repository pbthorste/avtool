package avtool

import (
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"crypto/hmac"
	"crypto/aes"
	"crypto/cipher"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

func Decrypt(filename, password string) string {
	data, err := ioutil.ReadFile(filename)
	check(err)
	body := splitHeader(data)
	salt, cryptedHmac, ciphertext := decodeData(body)
	key1, key2, iv := genKeyInitctr(password, salt)
	checkDigest(key2, cryptedHmac, ciphertext)
	aes_cipher,err := aes.NewCipher(key1)
	check(err)
	aesBlock := cipher.NewCTR(aes_cipher, iv)
	plaintext := make([]byte, len(ciphertext))
	aesBlock.XORKeyStream(plaintext, ciphertext)
	padding := int(plaintext[len(plaintext)-1])
	return string(plaintext[:len(plaintext)-padding])
}

/*
See _split_header function
https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L288
 */
func splitHeader(data []byte) string {
	contents := string(data)
	lines := strings.Split(contents,"\n")
	header := strings.Split(lines[0], ";")
	cipherName := strings.TrimSpace(header[2])
	if cipherName != "AES256" {
		panic("Error - unsupported cipher: " + cipherName)
	}
	body := strings.Join(lines[1:], "")
	return body
}
/*
See decrypt function (in class VaultAES256)
https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L741
 */
func decodeData(body string) (salt, cryptedHmac, ciphertext []byte) {
	decoded, _ := hex.DecodeString(body)
	elements := strings.SplitN(string(decoded), "\n", 3)
	salt, err1 := hex.DecodeString(elements[0])
	if err1 != nil {
		panic(err1)
	}
	cryptedHmac, err2 := hex.DecodeString(elements[1])
	if err2 != nil {
		panic(err2)
	}
	ciphertext, err3 := hex.DecodeString(elements[2])
	if err3 != nil {
		panic(err3)
	}
	return
}
/*
See function _gen_key_initctr (in class VaultAES256)
https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L685
 */
func genKeyInitctr(password string, salt []byte) (key1, key2, iv []byte) {
	keylength := 32
	ivlength := 16
	key := pbkdf2.Key([]byte(password), salt, 10000, 2*keylength+ivlength, sha256.New )
	key1 = key[:keylength]
	key2 = key[keylength:(keylength * 2)]
	iv = key[(keylength * 2):(keylength * 2) + ivlength]
	return
}
/*
See decrypt function (in class VaultAES256)
https://github.com/ansible/ansible/blob/0b8011436dc7f842b78298848e298f2a57ee8d78/lib/ansible/parsing/vault/__init__.py#L741
 */
func checkDigest(key2, cryptedHmac, ciphertext []byte) {
	hmacDecrypt := hmac.New(sha256.New, key2)
	hmacDecrypt.Write(ciphertext)
	expectedMAC := hmacDecrypt.Sum(nil)
	if !hmac.Equal(cryptedHmac, expectedMAC) {
		fmt.Println("ERROR - digests do not match - exiting")
		panic("error")
	}
}
