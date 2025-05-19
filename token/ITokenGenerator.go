package token

type ITokenGenerator interface {
	Generate(email string) (string, error)
}
