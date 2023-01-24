package handler

type TokenService interface {
	ValidateToken(token string) error
}
