package queue

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	Pending    JobStatus = "pending"
	Processing JobStatus = "processing"
	Completed  JobStatus = "completed"
	Failed     JobStatus = "failed"
)

type ScanValuer interface {
	sql.Scanner
	driver.Valuer
}

type Queue[T ScanValuer] interface {
	Enqueue(job Job[T]) error
	Dequeue() (Job[T], error)
	Acknowledge(jobID uuid.UUID) error
	Requeue(job Job[T], delay time.Duration) error
}

// HTTPJob represents the structure of a job to be processed
type Job[T ScanValuer] struct {
	ID         uuid.UUID `json:"id"`
	JobType    string    `json:"job_type"`
	Job        T         `json:"job"`
	Status     JobStatus `json:"status"`
	RetryCount int       `json:"retry_count"`
	DelayUntil time.Time `json:"delay_until"` // used for requeueing with delay
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
