package token

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

type TokenService struct {
	Token string
}

func NewTokenService(initialToken string) *TokenService {
	return &TokenService{Token: initialToken}
}

func (s *TokenService) ValidateToken(token string) error {
	if s.Token != token {
		log.Printf("access token %s is incorrect\n", token)
		return errors.New("incorrect access token")
	}
	return nil
}
