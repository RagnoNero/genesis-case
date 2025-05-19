package unit

import (
	"testing"
	"weather-subscription/token"
)

func TestTokenGenerator_GeneratesUnique(t *testing.T) {
	gen := token.NewTokenGenerator(10)
	t1, _ := gen.Generate("email1@example.com")
	t2, _ := gen.Generate("email1@example.com")

	if t1 == t2 {
		t.Error("expected unique tokens for different inputs")
	}
}
