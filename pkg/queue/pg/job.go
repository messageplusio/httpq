package pg

import (
	"time"

	"github.com/google/uuid"
	"github.com/messgeplusio/httpq/pkg/queue"
)

// Job represents the structure of a job to be processed
type Job[T queue.ScanValuer] struct {
	ID         uuid.UUID       `db:"id"`
	JobType    string          `db:"job_type"`
	Job        T               `db:"job"`
	RetryCount int             `db:"retry_count"`
	Status     queue.JobStatus `db:"status"`
	DelayUntil time.Time       `db:"delay_until"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

func (j Job[T]) AsEntity() queue.Job[T] {
	return queue.Job[T]{
		ID:         j.ID,
		Job:        j.Job,
		JobType:    j.JobType,
		RetryCount: j.RetryCount,
		Status:     j.Status,
		CreatedAt:  j.CreatedAt,
		UpdatedAt:  j.UpdatedAt,
	}
}

func (j *Job[T]) FromEntity(e queue.Job[T]) {
	if j == nil {
		return
	}
	j.ID = e.ID
	j.Job = e.Job
	j.JobType = e.JobType
	j.RetryCount = e.RetryCount
	j.Status = e.Status
	j.CreatedAt = e.CreatedAt
	j.UpdatedAt = e.UpdatedAt
}
