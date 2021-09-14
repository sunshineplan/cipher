package cipher

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	keyLen := []int{0, 5, 20, 50}
	dataLen := []int{10, 50, 200}
	for _, kl := range keyLen {
		key := random(kl)
		for _, dl := range dataLen {
			data := random(dl)
			result, err := Decrypt(key, Encrypt(key, data))
			if err != nil {
				t.Fatal(kl, dl, err)
			}
			if len(result) != 0 && !bytes.Equal(result, data) {
				t.Errorf("expected %v; got %v", data, result)
			}
		}
		plaintext := "测试"
		result, err := DecryptText(string(key), EncryptText(string(key), plaintext))
		if err != nil {
			t.Fatal(err)
		}
		if result != "" && result != plaintext {
			t.Errorf("expected %q; got %q", plaintext, result)
		}
	}
}

func random(len int) []byte {
	buff := make([]byte, len)
	rand.Read(buff)
	return buff
}
