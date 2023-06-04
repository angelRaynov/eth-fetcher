package authentication

type PasswordGetter interface {
	GetUserPassword(username string) (string, error)
}
