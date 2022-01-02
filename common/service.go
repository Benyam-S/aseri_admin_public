package common

import "github.com/Benyam-S/aseri/entity"

// IService is an interface that defines all the common service methods
type IService interface {
	AddJobAttribute(newJobAttribute *entity.JobAttribute, tableName string) error
	ValidateJobAttribute(tableName string, jobAttribute *entity.JobAttribute) error
	FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)
	AllJobAttributes(tableName string) []*entity.JobAttribute
	UpdateJobAttribute(jobAttribute *entity.JobAttribute, tableName string) error
	DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)
	ValidateJobAttributeTable(tableName string) error

	GetValidJobTypesName() []string
	GetValidJobTypes() []*entity.JobAttribute
	GetValidJobTypesForSubscription() []*entity.JobAttribute
	GetValidJobSectorsName() []string
	GetValidJobSectors() []*entity.JobAttribute
	GetValidJobSectorsForSubscription() []*entity.JobAttribute
	GetValidEducationLevelsName() []string
	GetValidEducationLevels() []*entity.JobAttribute
	GetValidEducationLevelsForSubscription() []*entity.JobAttribute
	GetValidWorkExperiences() []string
	GetValidWorkExperiencesForSubscription() []string
	GetValidContactTypes() []string

	FindReferral(code string) (*entity.Referral, error)
	AllReferrals() []*entity.Referral

	AddReferralVisit(newReferralVist *entity.ReferralVisit) error
	CountReferralVists(code string) int64
	CountAllReferralVists() int64
}

// IPushQueue is an interface all the method required from push queue struct
type IPushQueue interface {
	GetQueue() []*entity.ChannelRequest
	AddToQueue(request *entity.ChannelRequest)
	RemoveFromQueueWithIndex(index int)
}
