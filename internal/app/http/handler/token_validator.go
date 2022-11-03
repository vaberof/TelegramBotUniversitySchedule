package handler

type TokenValidator interface {
	ValidateToken(token string) error
}
