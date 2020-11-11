package auth

type IAuth interface {
	IsValid() bool
	GetUsername() string
}
