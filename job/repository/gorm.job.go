package repository

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/job"
	"github.com/Benyam-S/aseri/tools"
	"github.com/jinzhu/gorm"
)

// JobRepository is a type that defines a job repository type
type JobRepository struct {
	conn *gorm.DB
}

// NewJobRepository is a function that creates a new job repository type
func NewJobRepository(connection *gorm.DB) job.IJobRepository {
	return &JobRepository{conn: connection}
}

// Create is a method that adds a new job to the database
func (repo *JobRepository) Create(newJob *entity.Job) error {
	totalNumOfMembers := tools.CountMembers("jobs", repo.conn)
	newJob.ID = fmt.Sprintf("JB-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)

	for !tools.IsUnique("id", newJob.ID, "jobs", repo.conn) {
		totalNumOfMembers++
		newJob.ID = fmt.Sprintf("JB-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)
	}

	err := repo.conn.Create(newJob).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain job from the database using an identifier,
// also Find() uses only id as a key for selection
func (repo *JobRepository) Find(identifier string) (*entity.Job, error) {

	job := new(entity.Job)
	err := repo.conn.Model(job).Where("id = ?", identifier).First(job).Error

	if err != nil {
		return nil, err
	}
	return job, nil
}

// FindMultiple is a method that find multiple jobs from the database the matches the given identifier
// In FindMultiple() only employer is used as a key
func (repo *JobRepository) FindMultiple(identifier string) []*entity.Job {

	var jobs []*entity.Job
	err := repo.conn.Model(entity.Job{}).Where("employer = ?", identifier).Find(&jobs).Error

	if err != nil {
		return []*entity.Job{}
	}
	return jobs
}

// FindAll is a method that returns set of jobs limited to the page number and status
func (repo *JobRepository) FindAll(status string, pageNum int64) ([]*entity.Job, int64) {
	var jobs []*entity.Job
	var count float64

	switch status {
	case entity.JobStatusPending:
		fallthrough
	case entity.JobStatusOpened:
		fallthrough
	case entity.JobStatusClosed:
		fallthrough
	case entity.JobStatusDecelined:
		repo.conn.Raw("SELECT * FROM jobs WHERE status = ? ORDER BY created_at DESC LIMIT ?, 40", status, pageNum*40).Scan(&jobs)
		repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE status = ?", status).Count(&count)
		break
	case entity.JobStatusAny:
		fallthrough
	default:
		repo.conn.Raw("SELECT * FROM jobs ORDER BY created_at DESC LIMIT ?, 40", pageNum*40).Scan(&jobs)
		repo.conn.Raw("SELECT COUNT(*) FROM jobs").Count(&count)
	}

	var pageCount int64 = int64(math.Ceil(count / 40.0))
	return jobs, pageCount
}

// SearchWRegx is a method that searchs and returns set of jobs limited to the key identifier and page number using regular expersions
func (repo *JobRepository) SearchWRegx(key, status string, pageNum int64, columns ...string) ([]*entity.Job, int64) {
	var jobs []*entity.Job
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	switch status {
	case entity.JobStatusPending:
		fallthrough
	case entity.JobStatusOpened:
		fallthrough
	case entity.JobStatusClosed:
		fallthrough
	case entity.JobStatusDecelined:

		sqlValues = append(sqlValues, status)
		repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE ("+strings.Join(whereStmt, "||")+") && status = ?", sqlValues...).Count(&count)

		sqlValues = append(sqlValues, pageNum*40)
		repo.conn.Raw("SELECT * FROM jobs WHERE "+strings.Join(whereStmt, "||")+" && status = ? ORDER BY created_at DESC LIMIT ?, 40", sqlValues...).Scan(&jobs)

	case entity.JobStatusAny:
		fallthrough
	default:

		repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE ("+strings.Join(whereStmt, "||")+") ", sqlValues...).Count(&count)

		sqlValues = append(sqlValues, pageNum*40)
		repo.conn.Raw("SELECT * FROM jobs WHERE "+strings.Join(whereStmt, "||")+" ORDER BY created_at DESC LIMIT ?, 40", sqlValues...).Scan(&jobs)

	}

	var pageCount int64 = int64(math.Ceil(count / 40.0))
	return jobs, pageCount
}

// Search is a method that searchs and returns set of jobs limited to the key identifier and page number
func (repo *JobRepository) Search(key, status string, pageNum int64, columns ...string) ([]*entity.Job, int64) {
	var jobs []*entity.Job
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	switch status {
	case entity.JobStatusPending:
		fallthrough
	case entity.JobStatusOpened:
		fallthrough
	case entity.JobStatusClosed:
		fallthrough
	case entity.JobStatusDecelined:

		sqlValues = append(sqlValues, status)
		repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE ("+strings.Join(whereStmt, "||")+") && status = ?", sqlValues...).Count(&count)

		sqlValues = append(sqlValues, pageNum*40)
		repo.conn.Raw("SELECT * FROM jobs WHERE "+strings.Join(whereStmt, "||")+" && status = ? ORDER BY created_at DESC LIMIT ?, 40", sqlValues...).Scan(&jobs)

	case entity.JobStatusAny:
		fallthrough
	default:

		repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE ("+strings.Join(whereStmt, "||")+") ", sqlValues...).Count(&count)

		sqlValues = append(sqlValues, pageNum*40)
		repo.conn.Raw("SELECT * FROM jobs WHERE "+strings.Join(whereStmt, "||")+" ORDER BY created_at DESC LIMIT ?, 40", sqlValues...).Scan(&jobs)

	}

	var pageCount int64 = int64(math.Ceil(count / 40.0))
	return jobs, pageCount
}

// All is a method that returns all the jobs found in the database
func (repo *JobRepository) All() []*entity.Job {

	var jobs []*entity.Job

	repo.conn.Model(entity.Job{}).Find(&jobs).Order("created_at ASC")

	return jobs
}

// Total is a method that retruns the total number of jobs for the given status type
func (repo *JobRepository) Total(status string) int64 {

	var count int64
	if status == entity.JobStatusAny {
		repo.conn.Raw("SELECT COUNT(*) FROM jobs").Count(&count)
		return count
	}

	repo.conn.Raw("SELECT COUNT(*) FROM jobs WHERE status = ?", status).Count(&count)
	return count
}

// Update is a method that updates a certain job entries in the database
func (repo *JobRepository) Update(job *entity.Job) error {

	prevJob := new(entity.Job)
	err := repo.conn.Model(prevJob).Where("id = ?", job.ID).First(prevJob).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	job.CreatedAt = prevJob.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(job).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain job single column value in the database
func (repo *JobRepository) UpdateValue(job *entity.Job, columnName string, columnValue interface{}) error {

	prevJob := new(entity.Job)
	err := repo.conn.Model(prevJob).Where("id = ?", job.ID).First(prevJob).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.Job{}).Where("id = ?", job.ID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// CloseDueJobs is a method that updates multiple jobs' status to closed which have reached the given due date
// depending on the provided job status
func (repo *JobRepository) CloseDueJobs(dueDate time.Time, status string) []*entity.Job {

	var closeablejobs []*entity.Job
	err := repo.conn.Model(entity.Job{}).Where("status = ? && due_date <= ?", status, dueDate).Find(&closeablejobs).Error

	if err != nil {
		return []*entity.Job{}
	}

	repo.conn.Model(entity.Job{}).Where("status = ? && due_date <= ?", status, dueDate).
		Update(map[string]interface{}{"status": entity.JobStatusClosed})

	return closeablejobs
}

// Delete is a method that deletes a certain job from the database using an identifier.
// In Delete() id is only used as an key
func (repo *JobRepository) Delete(identifier string) (*entity.Job, error) {
	job := new(entity.Job)
	err := repo.conn.Model(job).Where("id = ?", identifier).First(job).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(job)
	return job, nil
}
