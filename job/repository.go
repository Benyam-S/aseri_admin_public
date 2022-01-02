package job

import (
	"time"

	"github.com/Benyam-S/aseri/entity"
)

// IJobRepository is an interface that defines all the repository methods of a job struct
type IJobRepository interface {
	Create(newJob *entity.Job) error
	Find(identifier string) (*entity.Job, error)
	FindMultiple(identifier string) []*entity.Job
	FindAll(status string, pageNum int64) ([]*entity.Job, int64)
	SearchWRegx(key, status string, pageNum int64, columns ...string) ([]*entity.Job, int64)
	Search(key, status string, pageNum int64, columns ...string) ([]*entity.Job, int64)
	All() []*entity.Job
	Total(status string) int64
	Update(job *entity.Job) error
	UpdateValue(job *entity.Job, columnName string, columnValue interface{}) error
	CloseDueJobs(time.Time, string) []*entity.Job
	Delete(identifier string) (*entity.Job, error)
}
