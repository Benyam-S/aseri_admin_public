package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
)

// Service is a type that defines a common service
type Service struct {
	commonRepo common.ICommonRepository
}

// NewCommonService is a function that returns a new common service
func NewCommonService(commonRepository common.ICommonRepository) common.IService {
	return &Service{commonRepo: commonRepository}
}

// AddJobAttribute is a method that adds a new job attribute to the system
func (service *Service) AddJobAttribute(newJobAttribute *entity.JobAttribute, tableName string) error {

	err := service.commonRepo.CreateJobAttribute(newJobAttribute, tableName)
	if err != nil {
		return errors.New("unable to add new job attribute")
	}

	return nil
}

// ValidateJobAttribute is a method that checks if the provided table and attribtue name is valid or not
func (service *Service) ValidateJobAttribute(tableName string, jobAttribute *entity.JobAttribute) error {

	if err := service.ValidateJobAttributeTable(tableName); err != nil {
		return err
	}

	empty, _ := regexp.MatchString(`^\s*$`, jobAttribute.Name)
	if empty {
		return errors.New("job attribute name can not be empty")
	}

	jobAttribute.Name = strings.TrimSpace(jobAttribute.Name)

	prevJobAttribute, _ := service.commonRepo.FindJobAttribute(jobAttribute.Name, tableName)
	if prevJobAttribute != nil {
		return errors.New("job attribute already exist")
	}

	return nil
}

// FindJobAttribute is a method that find and return job attribute that matches the given identifier
func (service *Service) FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("job attribute not found")
	}

	if err := service.ValidateJobAttributeTable(tableName); err != nil {
		return nil, err
	}

	jobAttribute, err := service.commonRepo.FindJobAttribute(identifier, tableName)
	if err != nil {
		return nil, errors.New("job attribute not found")
	}

	return jobAttribute, nil
}

// AllJobAttributes is a method that returns all the job attributes a given table
func (service *Service) AllJobAttributes(tableName string) []*entity.JobAttribute {

	if err := service.ValidateJobAttributeTable(tableName); err != nil {
		return []*entity.JobAttribute{}
	}

	return service.commonRepo.AllJobAttributes(tableName)
}

// UpdateJobAttribute is a method that updates a job attribute in the system
func (service *Service) UpdateJobAttribute(jobAttribute *entity.JobAttribute, tableName string) error {

	if err := service.ValidateJobAttributeTable(tableName); err != nil {
		return err
	}

	err := service.commonRepo.UpdateJobAttribute(jobAttribute, tableName)
	if err != nil {
		return errors.New("unable to update job attribute")
	}

	return nil
}

// DeleteJobAttribute is a method that deletes a job attribute from the system
func (service *Service) DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {

	if err := service.ValidateJobAttributeTable(tableName); err != nil {
		return nil, err
	}

	jobAttribute, err := service.commonRepo.DeleteJobAttribute(identifier, tableName)
	if err != nil {
		return nil, errors.New("unable to delete job attribute")
	}

	return jobAttribute, nil
}

// ValidateJobAttributeTable is a method that checks if the provided table name is valid or not
func (service *Service) ValidateJobAttributeTable(tableName string) error {

	switch tableName {
	case "job_types", "job_sectors", "education_levels":
		return nil
	}

	return errors.New("table does not exist")
}

// GetValidJobTypesName is a method that gets the valid job types name allowed by the system
func (service *Service) GetValidJobTypesName() []string {

	jobTypes := service.AllJobAttributes("job_types")
	validJobTypes := make([]string, 0)

	for _, jobType := range jobTypes {
		validJobTypes = append(validJobTypes, jobType.Name)
	}

	if len(validJobTypes) > 0 {
		validJobTypes = append(validJobTypes, "Other")
	}

	return validJobTypes
}

// GetValidJobTypes is a method that gets the valid job types allowed by the system
func (service *Service) GetValidJobTypes() []*entity.JobAttribute {

	jobTypes := service.AllJobAttributes("job_types")

	if len(jobTypes) > 0 {
		jobType := new(entity.JobAttribute)
		jobType.Name = "Other"
		jobTypes = append(jobTypes, jobType)
	}

	return jobTypes
}

// GetValidJobTypesForSubscription is a method that gets the valid job types allowed by the system to be used for subscription
func (service *Service) GetValidJobTypesForSubscription() []*entity.JobAttribute {
	jobTypes := service.AllJobAttributes("job_types")

	if len(jobTypes) > 0 {
		jobType := new(entity.JobAttribute)
		jobType.ID = "any"
		jobType.Name = "Any"
		jobTypes = append(jobTypes, jobType)
	}

	return jobTypes
}

// GetValidJobSectorsName is a method that gets the valid job sectors name allowed by the system
func (service *Service) GetValidJobSectorsName() []string {
	jobSectors := service.AllJobAttributes("job_sectors")
	validJobSectors := make([]string, 0)

	for _, jobSector := range jobSectors {
		validJobSectors = append(validJobSectors, jobSector.Name)
	}

	if len(validJobSectors) > 0 {
		validJobSectors = append(validJobSectors, "Other")
	}

	return validJobSectors
}

// GetValidJobSectors is a method that gets the valid job sectors allowed by the system
func (service *Service) GetValidJobSectors() []*entity.JobAttribute {
	jobSectors := service.AllJobAttributes("job_sectors")

	if len(jobSectors) > 0 {
		jobSector := new(entity.JobAttribute)
		jobSector.Name = "Other"
		jobSectors = append(jobSectors, jobSector)
	}

	return jobSectors
}

// GetValidJobSectorsForSubscription is a method that gets the valid job sectors allowed by the system to used for subscription
func (service *Service) GetValidJobSectorsForSubscription() []*entity.JobAttribute {
	return service.AllJobAttributes("job_sectors")
}

// GetValidEducationLevelsName is a method that gets the valid education levels name allowed by the system
func (service *Service) GetValidEducationLevelsName() []string {
	educationLevels := service.AllJobAttributes("education_levels")
	validEducationLevels := make([]string, 0)

	for _, educationLevel := range educationLevels {
		validEducationLevels = append(validEducationLevels, educationLevel.Name)
	}

	if len(validEducationLevels) > 0 {
		validEducationLevels = append(validEducationLevels, "Other")
	}

	return validEducationLevels
}

// GetValidEducationLevels is a method that gets the valid education levels allowed by the system
func (service *Service) GetValidEducationLevels() []*entity.JobAttribute {
	educationLevels := service.AllJobAttributes("education_levels")

	if len(educationLevels) > 0 {
		educationLevel := new(entity.JobAttribute)
		educationLevel.Name = "Other"
		educationLevels = append(educationLevels, educationLevel)
	}

	return educationLevels
}

// GetValidEducationLevelsForSubscription is a method that gets the valid education levels allowed by the system to used for subscription
func (service *Service) GetValidEducationLevelsForSubscription() []*entity.JobAttribute {
	educationLevels := service.AllJobAttributes("education_levels")

	if len(educationLevels) > 0 {
		educationLevel := new(entity.JobAttribute)
		educationLevel.ID = "any"
		educationLevel.Name = "Any"
		educationLevels = append(educationLevels, educationLevel)
	}

	return educationLevels
}

// GetValidWorkExperiences is a method that gets the valid work experiences allowed by the system
func (service *Service) GetValidWorkExperiences() []string {
	return entity.ValidWorkExperiences
}

// GetValidWorkExperiencesForSubscription is a method that gets the valid work experiences allowed by the system to used for subscription
func (service *Service) GetValidWorkExperiencesForSubscription() []string {
	experiences := entity.ValidWorkExperiences

	if len(experiences) > 0 {
		experience := new(entity.JobAttribute)
		experience.ID = "any"
		experience.Name = "Any"
		experiences = append(experiences, experience.Name)
	}

	return experiences
}

// GetValidContactTypes is a method that gets the valid contact type that can be used to contact job owner
func (service *Service) GetValidContactTypes() []string {
	return entity.ValidContactTypes
}

// FindReferral is a method that returns a certain referral that match the given code
func (service *Service) FindReferral(code string) (*entity.Referral, error) {

	empty, _ := regexp.MatchString(`^\s*$`, code)
	if empty {
		return nil, errors.New("no referral found")
	}

	referral, err := service.commonRepo.FindReferral(code)
	if err != nil {
		return nil, errors.New("no referral found")
	}
	return referral, nil
}

// AllReferrals is a method that returns all the referrals
func (service *Service) AllReferrals() []*entity.Referral {
	return service.commonRepo.AllReferrals()
}

// AddReferralVisit is a method that adds a new referral visit to the system
func (service *Service) AddReferralVisit(newReferralVist *entity.ReferralVisit) error {
	err := service.commonRepo.CreateReferralVisit(newReferralVist)
	if err != nil {
		return errors.New("unable to add new referral visit")
	}

	return nil
}

// CountReferralVists is a method that counts the number of visits for a certain referral
func (service *Service) CountReferralVists(code string) int64 {
	return service.commonRepo.CountReferralVists(code)
}

// CountAllReferralVists is a method that returns all the referrals
func (service *Service) CountAllReferralVists() int64 {
	return service.commonRepo.CountAllReferralVists()
}
