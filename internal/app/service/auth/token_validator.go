package auth

type TokenValidator interface {
	ValidateToken(token string) error
}
