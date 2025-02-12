package cipher

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	aesccm "github.com/pschlump/AesCCM"
)

const (
	iter    = 10000
	dkLen   = 16
	tagSize = 8
	ext     = ".enc"
)

var calcLen = aesccm.CalculateNonceLengthFromMessageLength

var errBlankKey = errors.New("blank key")

// EncryptText encrypts text by key.
func EncryptText(key, plaintext string) string {
	if key == "" {
		return strings.ReplaceAll(
			base64.StdEncoding.EncodeToString([]byte(plaintext)), "=", "",
		)
	}

	return strings.ReplaceAll(
		base64.StdEncoding.EncodeToString(Encrypt([]byte(key), []byte(plaintext))), "=", "",
	)
}

// DecryptText decrypts text by key.
func DecryptText(key, ciphertext string) (string, error) {
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

	plaintext, err := Decrypt([]byte(key), data)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptFile encrypts file by key.
func EncryptFile(key, file string) error {
	if key == "" {
		return errBlankKey
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return os.WriteFile(file+ext, Encrypt([]byte(key), data), 0666)
}

// DecryptFile decrypts file by key.
func DecryptFile(key, file string) (string, error) {
	if key == "" {
		return "", errBlankKey
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	data, err = Decrypt([]byte(key), data)
	if err != nil {
		return "", err
	}

	if filepath.Ext(file) == ext {
		file = file[:len(file)-len(ext)]
	} else {
		file = file + ".dec"
	}

	err = os.WriteFile(file, data, 0666)
	if err != nil {
		return "", err
	}

	return file, nil
}

// Encrypt encrypts bytes by key.
func Encrypt(key, data []byte) []byte {
	if len(key) == 0 {
		return data
	}

	salt := make([]byte, 8)
	rand.Read(salt)
	dk, err := pbkdf2.Key(sha256.New, string(key), salt, iter, dkLen)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err)
	}
	AesCCM, err := aesccm.NewCCM(block, tagSize, calcLen(len(data)))
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, 16)
	rand.Read(nonce)
	data, compression := compress(data)
	encrypted := AesCCM.Seal(nil, nonce, data, nil)

	return concat(salt, nonce, encrypted, []byte{compression})
}

// Decrypt decrypts bytes by key.
func Decrypt(key, data []byte) ([]byte, error) {
	if len(key) == 0 {
		return data, nil
	}
	if len(data) < 25 {
		return nil, errors.New("data below minimum length")
	}

	salt := data[:8]
	dk, err := pbkdf2.Key(sha256.New, string(key), salt, iter, dkLen)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}
	AesCCM, err := aesccm.NewCCM(block, tagSize, calcLen(len(data)-len(salt)-16-1-tagSize))
	if err != nil {
		return nil, err
	}
	decrypted, err := AesCCM.Open(nil, data[8:24], data[24:len(data)-1], nil)
	if err != nil {
		return nil, err
	}

	if data[len(data)-1] == '0' {
		return decrypted, nil
	}
	return decompress(decrypted)
}

func concat(b ...[]byte) (c []byte) {
	for _, i := range b {
		c = append(c, i...)
	}
	return
}

func compress(data []byte) ([]byte, byte) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()

	if b.Len() < len(data) {
		return b.Bytes(), '1'
	}
	return data, '0'
}

func decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)

	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
