package pow_test

import (
	"encoding/binary"
	"testing"

	"golang.org/x/crypto/sha3"
)

func BenchmarkX16(b *testing.B) {
	nonce := []byte{}
	for i := uint16(0); i < 16; i++ {
		n := make([]byte, 2)
		binary.LittleEndian.PutUint16(n[0:], i)
		nonce = append(nonce, n...)
	}

	for i := 0; i < b.N; i++ {
		sha3.Sum256(nonce)
	}
}

func BenchmarkX1024(b *testing.B) {
	nonce := []byte{}
	for i := uint16(0); i < 1024; i++ {
		n := make([]byte, 2)
		binary.LittleEndian.PutUint16(n[0:], i)
		nonce = append(nonce, n...)
	}

	for i := 0; i < b.N; i++ {
		sha3.Sum256(nonce)
	}
}
