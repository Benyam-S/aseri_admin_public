package common

import "github.com/Benyam-S/aseri/entity"

// ICommonRepository is an interface that defines all the common repository methods
type ICommonRepository interface {
	IsUnique(columnName string, columnValue interface{}, tableName string) bool

	CreateJobAttribute(newAttribute *entity.JobAttribute, tableName string) error
	FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)
	AllJobAttributes(tableName string) []*entity.JobAttribute
	UpdateJobAttribute(attribute *entity.JobAttribute, tableName string) error
	DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)

	CreateReferral(newReferral *entity.Referral) error
	FindReferral(code string) (*entity.Referral, error)
	AllReferrals() []*entity.Referral
	UpdateReferral(referral *entity.Referral) error
	DeleteReferral(code string) (*entity.Referral, error)

	CreateReferralVisit(newReferralVist *entity.ReferralVisit) error
	FindReferralVisits(code string) []*entity.ReferralVisit
	AllReferralVisits() []*entity.ReferralVisit
	CountReferralVists(code string) int64
	CountAllReferralVists() int64
}
