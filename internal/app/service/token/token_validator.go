package token

type TokenValidator interface {
	ValidateToken(token string) error
}
