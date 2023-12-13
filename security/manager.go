package security

type AuthManager interface {
	GetUser(any) (IUser, error)
}
