package service

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/jobapplication"
)

// Service is a type that defines a job application service
type Service struct {
	jobApplicationRepo jobapplication.IJobApplicationRepository
	commonRepo         common.ICommonRepository
}

// NewJobApplicationService is a function that returns a new job application service
func NewJobApplicationService(jobApplicationRepository jobapplication.IJobApplicationRepository,
	commonRepository common.ICommonRepository) jobapplication.IService {
	return &Service{jobApplicationRepo: jobApplicationRepository, commonRepo: commonRepository}
}

// AddJobApplication is a method that adds a new job application to the system
func (service *Service) AddJobApplication(newJobApplication *entity.JobApplication) error {
	err := service.jobApplicationRepo.Create(newJobApplication)
	if err != nil {
		return errors.New("unable to add new job application")
	}

	return nil
}

// FindJobApplications is a method that find and return multiple job applications that match the identifier value
func (service *Service) FindJobApplications(identifier string) []*entity.JobApplication {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.JobApplication{}
	}

	return service.jobApplicationRepo.Find(identifier)
}

// JobApplicationExists is a method that checks whether the job application exists or not
func (service *Service) JobApplicationExists(jobID, jobSeekerID string) bool {
	return service.jobApplicationRepo.HasApplied(jobID, jobSeekerID)
}

// DeleteJobApplication is a method that deletes a certain job application from the system
func (service *Service) DeleteJobApplication(jobID, jobSeekerID string) (*entity.JobApplication, error) {

	jobApplication, err := service.jobApplicationRepo.Delete(jobID, jobSeekerID)
	if err != nil {
		return nil, errors.New("unable to delete job application")
	}

	return jobApplication, nil
}

// DeleteMultipleJobApplications is a method that deletes multiple jobs from the system that match the given identifier
func (service *Service) DeleteMultipleJobApplications(identifier string) []*entity.JobApplication {
	return service.jobApplicationRepo.DeleteMultiple(identifier)
}
