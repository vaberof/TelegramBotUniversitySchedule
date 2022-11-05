package auth

type TokenService interface {
	ValidateToken(token string) error
}
