package model

type Interface interface {
	GetResourceOwnerId(token string) (userId uint)
}
