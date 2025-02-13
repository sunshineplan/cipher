package cipher

import (
	"bytes"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"testing"
	"time"
)

func benchmarkPBKDF2(iter int, t *testing.T) {
	start := time.Now()
	pbkdf2.Key(sha256.New, string("testpassword"), []byte("testsalt"), iter, keyLength)
	duration := time.Since(start)
	t.Logf("Iterations: %d, Duration: %v", iter, duration)
}

func TestPBKDF2Performance(t *testing.T) {
	iterationsList := []int{10000, 50000, 200000, 500000, 1000000}
	for _, iter := range iterationsList {
		benchmarkPBKDF2(iter, t)
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pbkdf2.Key(sha256.New, string("testpassword"), []byte("testsalt"), iter, keyLength)
	}
	b.ReportAllocs()
}

func TestDecryptText(t *testing.T) {
	plaintext, err := DecryptText("测试Key", "ZcDV7Bew3jHc0rnfI7u6FYSdc4kTp7R8Cs8QTt1+oPqoL/eoEI6eXqQ+5WsY1DplxgA")
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
