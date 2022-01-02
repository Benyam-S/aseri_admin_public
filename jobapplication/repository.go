package jobapplication

import "github.com/Benyam-S/aseri/entity"

// IJobApplicationRepository is an interface that defines all the repository methods of a job application struct
type IJobApplicationRepository interface {
	Create(newJobApplication *entity.JobApplication) error
	Find(identifier string) []*entity.JobApplication
	HasApplied(jobID, jobSeekerID string) bool
	Delete(jobID, jobSeekerID string) (*entity.JobApplication, error)
	DeleteMultiple(identifier string) []*entity.JobApplication
}
