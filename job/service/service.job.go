package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/job"
	"github.com/Benyam-S/aseri/tools"
	"github.com/Benyam-S/aseri/user"
)

// Service is a type that defines a job service
type Service struct {
	jobRepo   job.IJobRepository
	userRepo  user.IUserRepository
	cmService common.IService
}

// NewJobService is a function that returns a new job service
func NewJobService(jobRepository job.IJobRepository, userRepository user.IUserRepository,
	commonService common.IService) job.IService {
	return &Service{jobRepo: jobRepository, userRepo: userRepository, cmService: commonService}
}

// AddJob is a method that adds a new job to the system
func (service *Service) AddJob(newJob *entity.Job) error {

	// Initiating new job
	if newJob.Status == "" {
		newJob.Status = entity.JobStatusPending
	}

	err := service.jobRepo.Create(newJob)
	if err != nil {
		return errors.New("unable to add new job")
	}

	return nil
}

// ValidateJob is a method that validates a job entries.
// It checks if the job has a valid entries or not and return map of errors if any.
func (service *Service) ValidateJob(job *entity.Job) entity.ErrMap {

	errMap := make(map[string]error)
	var isValidJobType bool
	var isValidJobSector bool
	var isValidEducationLevel bool
	var isValidWorkExperience bool
	var isValidContactType bool

	validJobTypes := service.cmService.GetValidJobTypesName()
	validJobSectors := service.cmService.GetValidJobSectorsName()
	validEducationLevels := service.cmService.GetValidEducationLevelsName()
	validWorkExperiences := service.cmService.GetValidWorkExperiences()
	validContactTypes := service.cmService.GetValidContactTypes()

	// Spliting the job sectors and types
	var jobSectors = strings.Split(strings.TrimSpace(job.Sector), ",")
	var jobTypes = strings.Split(strings.TrimSpace(job.Type), ",")

	emptyEmployer, _ := regexp.MatchString(`^\s*$`, job.Employer)
	if emptyEmployer {
		errMap["employer"] = errors.New("employer must be specified")
	}

	emptyTitle, _ := regexp.MatchString(`^\s*$`, job.Title)
	if emptyTitle {
		errMap["title"] = errors.New("job title can not be empty")
	} else if len(job.Title) > 300 {
		errMap["title"] = errors.New("job title can not exceed 300 characters")
	}

	emptyDescription, _ := regexp.MatchString(`^\s*$`, job.Description)
	if emptyDescription {
		errMap["description"] = errors.New("job description can not be empty")
	} else if len(job.Description) > 20000 {
		errMap["description"] = errors.New("job description can not exceed 20,000 characters")
	}

	if job.Link != "" && !tools.IsValidURL(job.Link) {
		errMap["link"] = errors.New("invalid link has bee used")
	} else if len(job.Link) > 1000 {
		errMap["link"] = errors.New("link can not exceed 1000 characters")
	}

	emptyContactInfo, _ := regexp.MatchString(`^\s*$`, job.ContactInfo)
	if emptyContactInfo {
		errMap["contact_info"] = errors.New("contact info must be specified")
	}

	for _, jobType := range jobTypes {
		isValidJobType = false
		for _, validJobType := range validJobTypes {
			if strings.ToLower(strings.TrimSpace(jobType)) == strings.ToLower(validJobType) {
				isValidJobType = true
				break
			}
		}

		if !isValidJobType {
			break
		}
	}

	for _, jobSector := range jobSectors {
		isValidJobSector = false
		for _, validJobSector := range validJobSectors {
			if strings.ToLower(strings.TrimSpace(jobSector)) == strings.ToLower(validJobSector) {
				isValidJobSector = true
				break
			}
		}

		if !isValidJobSector {
			break
		}
	}

	for _, validEducationLevel := range validEducationLevels {
		if strings.ToLower(strings.TrimSpace(job.EducationLevel)) == strings.ToLower(validEducationLevel) {
			isValidEducationLevel = true
			break
		}
	}

	for _, validWorkExperience := range validWorkExperiences {
		if strings.ToLower(strings.TrimSpace(job.Experience)) == strings.ToLower(validWorkExperience) {
			isValidWorkExperience = true
			break
		}
	}

	for _, validContactType := range validContactTypes {
		if strings.ToLower(strings.TrimSpace(job.ContactType)) == strings.ToLower(validContactType) {
			isValidContactType = true
			break
		}
	}

	if !isValidJobType {
		errMap["type"] = errors.New("invalid job type used")
	}

	if !isValidJobSector {
		errMap["sector"] = errors.New("invalid job sector used")
	}

	if !isValidEducationLevel {
		errMap["education_level"] = errors.New("invalid education level used")
	}

	if !isValidWorkExperience {
		errMap["experience"] = errors.New("invalid work experience used")
	}

	if !isValidContactType {
		errMap["contact_type"] = errors.New("invalid contact type selected")
	}

	switch strings.ToLower(job.Gender) {

	case "m", "f", "b", "male", "female", "both":
		if strings.ToLower(job.Gender) == "m" || strings.ToLower(job.Gender) == "male" {
			job.Gender = "M"
		} else if strings.ToLower(job.Gender) == "f" || strings.ToLower(job.Gender) == "female" {
			job.Gender = "F"
		} else if strings.ToLower(job.Gender) == "b" || strings.ToLower(job.Gender) == "both" {
			job.Gender = "B"
		}

	default:
		errMap["gender"] = errors.New("invalid gender selection")
	}

	dateBase := time.Now().Add(time.Hour * 3)
	if job.DueDate != nil && job.DueDate.Unix() < dateBase.Unix() {
		errMap["due_date"] = errors.New("due date must exceed the current time at least by 3 hours")
	}

	// Making the post type 'User' if invalid
	if job.PostType != entity.PostCategoryExternal &&
		job.PostType != entity.PostCategoryUser &&
		job.PostType != entity.PostCategoryInternal {
		job.PostType = entity.PostCategoryUser
	}

	switch job.PostType {
	case entity.PostCategoryUser:
		if errMap["employer"] == nil {
			user, err := service.userRepo.Find(job.Employer)
			if err != nil {
				errMap["employer"] = errors.New("no user found for the provided employer id")
			} else if user.Category == entity.UserCategoryJobSeeker {
				errMap["employer"] = errors.New("can not perform operation for job seeker")
			}
		}

		delete(errMap, "contact_info")
		delete(errMap, "link")

		// Cleaning unused data for security purpose
		job.ContactInfo = ""
		job.Link = ""

	case entity.PostCategoryInternal:

		delete(errMap, "contact_type")
		delete(errMap, "link")

		// Cleaning unused data for security purpose
		job.ContactType = ""
		job.Link = ""

	case entity.PostCategoryExternal:
		delete(errMap, "contact_type")
		delete(errMap, "contact_info")
		delete(errMap, "gender")

		emptyJobType, _ := regexp.MatchString(`^\s*$`, job.Type)
		if emptyJobType {
			delete(errMap, "type")
		}

		emptyEducationLevel, _ := regexp.MatchString(`^\s*$`, job.EducationLevel)
		if emptyEducationLevel {
			delete(errMap, "education_level")
		}

		// Cleaning unused data for security purpose
		job.ContactType = ""
		job.ContactInfo = ""
		job.Gender = ""
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindJob is a method that find and return a job that matchs the identifier value
func (service *Service) FindJob(identifier string) (*entity.Job, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no job found")
	}

	job, err := service.jobRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no job found")
	}
	return job, nil
}

// FindMultipleJobs is a method that find and return multiple jobs that matchs the identifier value
func (service *Service) FindMultipleJobs(identifier string) []*entity.Job {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.Job{}
	}

	return service.jobRepo.FindMultiple(identifier)
}

// AllJobs is a method that returns all the jobs in the system
func (service *Service) AllJobs() []*entity.Job {
	return service.jobRepo.All()
}

// CloseDueJobs is a method that closes all the jobs that have reached their due date for given job status
func (service *Service) CloseDueJobs(dueDate time.Time, status string) []*entity.Job {
	return service.jobRepo.CloseDueJobs(dueDate, status)
}

// AllJobsWithPagination is a method that returns all the jobs with pagination
func (service *Service) AllJobsWithPagination(status string, pageNum int64) ([]*entity.Job, int64) {

	if status != entity.JobStatusPending && status != entity.JobStatusOpened &&
		status != entity.JobStatusClosed && status != entity.JobStatusDecelined {
		status = entity.JobStatusAny
	}

	return service.jobRepo.FindAll(status, pageNum)
}

// TotalJobs is a method that returns the total number of jobs for a given job status
func (service *Service) TotalJobs(status string) int64 {

	if status != entity.JobStatusClosed && status != entity.JobStatusOpened &&
		status != entity.JobStatusPending && status != entity.JobStatusDecelined {
		status = entity.JobStatusAny
	}

	return service.jobRepo.Total(status)
}

// SearchJobs is a method that searchs and returns a set of jobs related to the key identifier
func (service *Service) SearchJobs(key, status string, pageNum int64, extra ...string) ([]*entity.Job, int64) {

	if status != entity.JobStatusPending && status != entity.JobStatusOpened &&
		status != entity.JobStatusClosed && status != entity.JobStatusDecelined {
		status = entity.JobStatusAny
	}

	defaultSearchColumnsRegx := []string{"title"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "employer", "type"}

	result1 := make([]*entity.Job, 0)
	result2 := make([]*entity.Job, 0)
	results := make([]*entity.Job, 0)
	resultsMap := make(map[string]*entity.Job)
	var pageCount1 int64 = 0
	var pageCount2 int64 = 0
	var pageCount int64 = 0

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results, 0
	}

	result1, pageCount1 = service.jobRepo.Search(key, status, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2, pageCount2 = service.jobRepo.SearchWRegx(key, status, pageNum, defaultSearchColumnsRegx...)
	}

	for _, job := range result1 {
		resultsMap[job.ID] = job
	}

	for _, job := range result2 {
		resultsMap[job.ID] = job
	}

	for _, uniqueJob := range resultsMap {
		results = append(results, uniqueJob)
	}

	pageCount = pageCount1
	if pageCount < pageCount2 {
		pageCount = pageCount2
	}

	return results, pageCount
}

// UpdateJob is a method that updates a job in the system
func (service *Service) UpdateJob(job *entity.Job) error {
	err := service.jobRepo.Update(job)
	if err != nil {
		return errors.New("unable to update job")
	}

	return nil
}

// UpdateJobSingleValue is a method that updates a single column entry of a job
func (service *Service) UpdateJobSingleValue(jobID, columnName string, columnValue interface{}) error {
	job := entity.Job{ID: jobID}
	err := service.jobRepo.UpdateValue(&job, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update job")
	}

	return nil
}

// ChangeJobStatus is a method that changes the given job status
func (service *Service) ChangeJobStatus(jobID, status string) (*entity.Job, error) {

	job, err := service.jobRepo.Find(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if (job.Status != entity.JobStatusPending && job.Status != entity.JobStatusOpened) ||
		(status != entity.JobStatusDecelined && status != entity.JobStatusOpened &&
			status != entity.JobStatusClosed) {
		return nil, errors.New("unable to perform operation")
	}

	if job.Status == entity.JobStatusPending && (status != entity.JobStatusDecelined &&
		status != entity.JobStatusOpened) {
		return nil, errors.New("unable to perform operation")
	} else if job.Status == entity.JobStatusOpened && status != entity.JobStatusClosed {
		return nil, errors.New("unable to perform operation")
	}

	job.Status = status
	err = service.jobRepo.Update(job)
	if err != nil {
		return nil, errors.New("unable to update job")
	}

	return job, nil
}

// DeleteJob is a method that deletes a job from the system
func (service *Service) DeleteJob(jobID string) (*entity.Job, error) {

	job, err := service.jobRepo.Delete(jobID)
	if err != nil {
		return nil, errors.New("unable to delete job")
	}

	return job, nil
}
