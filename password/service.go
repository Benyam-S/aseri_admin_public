package password

import "github.com/Benyam-S/aseri/entity"

// IService is an interface that defines all the service methods of a password struct
type IService interface {
	FindPassword(identifier string) (*entity.Password, error)
	VerifyPassword(memberPassword *entity.Password, verifyPassword string) error
	UpdatePassword(memberPassword *entity.Password) error
	DeletePassword(identifier string) (*entity.Password, error)
}
