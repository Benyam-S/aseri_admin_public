package staff

import (
	"github.com/Benyam-S/aseri/client/http/session"
	"github.com/Benyam-S/aseri/entity"
)

// IStaffRepository is an interface that defines all the repository methods of a staff struct
type IStaffRepository interface {
	Create(newStaffMember *entity.Staff) error
	Find(identifier string) (*entity.Staff, error)
	FindAll(role string, pageNum int64) ([]*entity.Staff, int64)
	Search(key, role string, pageNum int64, columns ...string) ([]*entity.Staff, int64)
	SearchWRegx(key, role string, pageNum int64, columns ...string) ([]*entity.Staff, int64)
	Update(staffMember *entity.Staff) error
	UpdateValue(staffMember *entity.Staff, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.Staff, error)
}

// ISessionRepository is an interface that defines all the repository methods of a user's server side session struct
type ISessionRepository interface {
	Create(newUserSession *session.ServerSession) error
	Find(identifier string) (*session.ServerSession, error)
	Search(identifier string) ([]*session.ServerSession, error)
	Update(userSession *session.ServerSession) error
	Delete(identifier string) (*session.ServerSession, error)
	DeleteMultiple(identifier string) ([]*session.ServerSession, error)
}
