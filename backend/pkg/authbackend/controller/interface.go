package controller

// Repository stores everything for the frontend API
type Repository interface {
	GetOrCreateUser(authOrigin, authSubject string) (ret User, err error)
}
