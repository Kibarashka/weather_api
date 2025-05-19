package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type TokenService interface {
	GenerateToken(byteLength int) (string, error)
}

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", fmt.Errorf("byteLength must be positive")
	}
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("tokenService.GenerateToken: failed to read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
