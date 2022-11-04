package auth

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

type AuthService struct {
	Token string
}

func NewAuthService(initialToken string) *AuthService {
	return &AuthService{Token: initialToken}
}

func (s *AuthService) ValidateToken(token string) error {
	if s.Token != token {
		log.Printf("access token %s is incorrect\n", token)
		return errors.New("incorrect access token")
	}
	return nil
}
