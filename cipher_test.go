package cipher

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestDecryptText(t *testing.T) {
	plaintext, err := DecryptText("key", "B7y+ifBCt1AUqly/5Lku2EottQJ64wxUHIAt6KpF0m1wMOxhhP9cH//+pqx/MA")
	if err != nil {
		t.Fatal(err)
	}
	if expect := "Hello, 世界"; plaintext != expect {
		t.Errorf("expected %q; got %q", expect, plaintext)
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	keyLen := []int{0, 5, 20, 50}
	dataLen := []int{10, 2000, 30000, 400000, 5000000}
	for _, kl := range keyLen {
		for i := 0; i < 5; i++ {
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
			plaintext := "Hello, 世界"
			result, err := DecryptText(string(key), EncryptText(string(key), plaintext))
			if err != nil {
				t.Fatal(err)
			}
			if result != "" && result != plaintext {
				t.Errorf("expected %q; got %q", plaintext, result)
			}
		}
	}
}

func random(len int) []byte {
	buff := make([]byte, len)
	rand.Read(buff)
	return buff
}
