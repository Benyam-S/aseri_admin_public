package jobapplication

import "github.com/Benyam-S/aseri/entity"

// IService is an interface that defines all the service methods of a job application struct
type IService interface {
	AddJobApplication(newJobApplication *entity.JobApplication) error
	FindJobApplications(identifier string) []*entity.JobApplication
	JobApplicationExists(jobID, jobSeekerID string) bool
	DeleteJobApplication(jobID, jobSeekerID string) (*entity.JobApplication, error)
	DeleteMultipleJobApplications(identifier string) []*entity.JobApplication
}
