package crypto

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

func CompareLegacyHashAndPassword(ctx context.Context, hash, password string) error {
	if password == "" || hash == "" {
		return errors.New("password or hash can't empty")
	}
	splitedHash := strings.Split(hash, ":")
	if len(splitedHash) < 3 {
		return errors.New("failed legacy hash format")
	}

	legacyHash := splitedHash[0]
	salt := splitedHash[1]

	// use ARGON2ID13
	hashedArgon := getArgonHash([]byte(password), salt)

	match := compareHashes([]byte(legacyHash), []byte(hashedArgon))
	if match {
		return nil
	} else {
		return errors.New("password is invalid")
	}
}

func getArgonHash(data []byte, salt string) string {
	if salt == "" {
		saltBytes := make([]byte, 16)
		rand.Read(saltBytes)
		salt = hex.EncodeToString(saltBytes)
	} else {
		salt = salt[:16]
		if len(salt) < 16 {
			salt = padString(salt, 16)
		}
	}

	hash := argon2.IDKey(data, []byte(salt), 2, 64*1024, 1, 32)

	return hex.EncodeToString(hash)
}

func padString(s string, length int) string {
	for len(s) < length {
		s += s
	}
	return s[:length]
}

func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}

	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			return false
		}
	}

	return true
}
