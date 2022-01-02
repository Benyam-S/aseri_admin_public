package repository

import (
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/jobapplication"
	"github.com/jinzhu/gorm"
)

// JobApplicationRepository is a type that defines a job application repository type
type JobApplicationRepository struct {
	conn *gorm.DB
}

// NewJobApplicationRepository is a function that creates a new job application repository type
func NewJobApplicationRepository(connection *gorm.DB) jobapplication.IJobApplicationRepository {
	return &JobApplicationRepository{conn: connection}
}

// Create is a method that adds a new job application to the database
func (repo *JobApplicationRepository) Create(newJobApplication *entity.JobApplication) error {
	err := repo.conn.Create(newJobApplication).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that searches and returns a set of job applications that are limited to the provided identifier.
// In Find() job_id or job_seeker_id can either be used as a key.
func (repo *JobApplicationRepository) Find(identifier string) []*entity.JobApplication {
	var jobApplications []*entity.JobApplication
	err := repo.conn.Model(entity.JobApplication{}).
		Where("job_id = ? || job_seeker_id = ?", identifier, identifier).
		Find(&jobApplications).Error

	if err != nil {
		return []*entity.JobApplication{}
	}
	return jobApplications
}

// HasApplied is a method that checks whether a job seeker has applied to a certain job by checking
// the job application relation table
func (repo *JobApplicationRepository) HasApplied(jobID, jobSeekerID string) bool {

	jobApplication := new(entity.JobApplication)
	err := repo.conn.Model(jobApplication).
		Where("job_id = ? && job_seeker_id = ?", jobID, jobSeekerID).
		First(jobApplication).Error

	if err != nil || jobApplication.JobID == "" {
		return false
	}

	return true
}

// Delete is a method that deletes a certain job application from the database using job_id and job_seeker_id.
func (repo *JobApplicationRepository) Delete(jobID, jobSeekerID string) (*entity.JobApplication, error) {
	jobApplication := new(entity.JobApplication)
	err := repo.conn.Model(jobApplication).Where("job_id = ? && job_seeker_id = ?",
		jobID, jobSeekerID).First(jobApplication).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(jobApplication)
	return jobApplication, nil
}

// DeleteMultiple is a method that deletes a set of job applications from the database using an identifier.
// In DeleteMultiple() job_id or job_seeker_id can either be used as a key.
func (repo *JobApplicationRepository) DeleteMultiple(identifier string) []*entity.JobApplication {
	var jobApplications []*entity.JobApplication
	repo.conn.Model(jobApplications).Where("job_id = ? || job_seeker_id = ?", identifier, identifier).
		Find(&jobApplications)

	for _, jobApplication := range jobApplications {
		repo.conn.Delete(jobApplication)
	}

	return jobApplications
}
