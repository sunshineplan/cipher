package cipher

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"strings"

	aesccm "github.com/pschlump/AesCCM"
	"golang.org/x/crypto/pbkdf2"
)

// Encrypt encrypts text by key.
func Encrypt(key, plaintext string) string {
	if key == "" {
		return strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(plaintext)), "=", "")
	}
	salt := make([]byte, 8)
	rand.Read(salt)
	dk := pbkdf2.Key([]byte(key), salt, 10000, 16, sha256.New)
	Aes, _ := aes.NewCipher(dk)
	AesCCM, _ := aesccm.NewCCM(Aes, 8, 13)
	nonce := make([]byte, 16)
	rand.Read(nonce)
	data, compression := compress(plaintext)
	ciphertext := AesCCM.Seal(nil, nonce, data, nil)
	return strings.ReplaceAll(base64.StdEncoding.EncodeToString(concat(salt, nonce, ciphertext, compression)), "=", "")
}

// Decrypt decrypts text by key.
func Decrypt(key, ciphertext string) (string, error) {
	if r := len(ciphertext) % 4; r > 0 {
		ciphertext += strings.Repeat("=", 4-r)
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if key == "" {
		return string(data), nil
	}
	salt := data[:8]
	dk := pbkdf2.Key([]byte(key), salt, 10000, 16, sha256.New)
	Aes, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}
	AesCCM, _ := aesccm.NewCCM(Aes, 8, 13)
	plaintext, err := AesCCM.Open(nil, data[8:24], data[24:len(data)-1], nil)
	if err != nil {
		return "", err
	}
	if data[len(data)-1] == []byte("0")[0] {
		return string(plaintext), nil
	}
	return decompress(plaintext)
}

func concat(b ...[]byte) (c []byte) {
	for _, i := range b {
		c = append(c, i...)
	}
	return
}

func compress(data string) ([]byte, []byte) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(data))
	w.Close()
	if b.Len() < len([]byte(data)) {
		return b.Bytes(), []byte("1")
	}
	return []byte(data), []byte("0")
}

func decompress(data []byte) (string, error) {
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	defer r.Close()
	decompressed, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(decompressed), nil
}
