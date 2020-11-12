package cipher

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	keyLen := []int{0, 5, 20, 50}
	plaintextLen := []int{10, 50, 200}
	for _, kl := range keyLen {
		key := randomString(kl)
		for _, pl := range plaintextLen {
			plaintext := randomString(pl)
			pt, err := Decrypt(key, Encrypt(key, plaintext))
			if err != nil {
				t.Error(kl, pl, err)
			}
			if pt != "" && pt != plaintext {
				t.Errorf("expected %q; got %q", plaintext, pt)
			}
		}
		plaintext := "测试"
		pt, err := Decrypt(key, Encrypt(key, plaintext))
		if err != nil {
			t.Error(err)
		}
		if pt != "" && pt != plaintext {
			t.Errorf("expected %q; got %q", plaintext, pt)
		}
	}
}

func randomString(l int) string {
	buff := make([]byte, int(math.Round(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}
