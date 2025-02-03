package auth

import "testing"

func TestHash(t *testing.T) {
	hash, _ := hashPassword("password")

	result := verifyPassword("password", hash)

	if !result {
		t.Errorf("failed to verify password")
	}

}
