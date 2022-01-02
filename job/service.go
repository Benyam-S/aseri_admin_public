package job

import (
	"time"

	"github.com/Benyam-S/aseri/entity"
)

// IService is an interface that defines all the service methods of a job struct
type IService interface {
	AddJob(newJob *entity.Job) error
	ValidateJob(job *entity.Job) entity.ErrMap
	FindJob(identifier string) (*entity.Job, error)
	FindMultipleJobs(identifier string) []*entity.Job
	AllJobs() []*entity.Job
	AllJobsWithPagination(status string, pageNum int64) ([]*entity.Job, int64)
	SearchJobs(key, status string, pageNum int64, extra ...string) ([]*entity.Job, int64)
	TotalJobs(status string) int64
	UpdateJob(job *entity.Job) error
	UpdateJobSingleValue(jobID, columnName string, columnValue interface{}) error
	ChangeJobStatus(jobID, status string) (*entity.Job, error)
	CloseDueJobs(time.Time, string) []*entity.Job
	DeleteJob(jobID string) (*entity.Job, error)
}
