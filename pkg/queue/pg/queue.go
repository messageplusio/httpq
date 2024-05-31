package pg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/messgeplusio/httpq/pkg/queue"
)

// PostgreSQLQueue is the implementation of the job queue using PostgreSQL
type PostgreSQLQueue[T queue.ScanValuer] struct {
	db *sqlx.DB
}

// NewPostgreSQLQueue initializes a new PostgreSQLQueue
func NewPostgreSQLQueue[T queue.ScanValuer](connStr string) (*PostgreSQLQueue[T], error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgreSQLQueue[T]{db: db}, nil
}

// Enqueue adds a new job to the queue
func (q *PostgreSQLQueue[T]) Enqueue(ctx context.Context, job queue.Job[T]) (queue.Job[T], error) {
	job.Status = "pending"
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}
	stmt := `INSERT INTO "jobs" ("id", "job_type", "job", "retry_count", "status", "created_at", "updated_at") VALUES (:id, :job_type, :job, :retry_count, :status, :created_at, :updated_at)`
	_, err := q.db.NamedExecContext(ctx, stmt, job)
	return job, err
}

// Dequeue retrieves and locks the next pending job from the queue
func (q *PostgreSQLQueue[T]) Dequeue(ctx context.Context) (queue.Job[T], error) {
	var job queue.Job[T]
	// update the status of the job to 'processing' and return the job
	// for one job that is in pending status and has the oldest created_at timestamp
	err := q.db.GetContext(ctx, &job, `UPDATE jobs SET status = $1 WHERE id = (SELECT id FROM jobs WHERE status = $2 ORDER BY delay_until ASC LIMIT 1)  RETURNING *`,
		queue.Processing,
		queue.Pending)
	return job, err
}

// Dequeue retrieves and locks the next pending job from the queue
func (q *PostgreSQLQueue[T]) CurrentJobs(ctx context.Context) ([]queue.Job[T], error) {
	var jobs []queue.Job[T]
	err := q.db.SelectContext(ctx, &jobs, `SELECT * FROM jobs WHERE status = $1`, queue.Processing)
	return jobs, err
}

// Acknowledge marks a job as completed
func (q *PostgreSQLQueue[T]) Acknowledge(ctx context.Context, jobID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, `UPDATE jobs SET status = $1 WHERE id = $2 and status= $3`, jobID)
	return err
}

// Requeue requeues a job with a delay
func (q *PostgreSQLQueue[T]) Requeue(ctx context.Context, job queue.Job[T], delay time.Duration) error {
	stmt := `UPDATE jobs SET status = $1, delay_until = $2 WHERE id = $3`
	_, err := q.db.ExecContext(ctx, stmt, queue.Pending, time.Now().Add(delay), job.ID)
	return err
}
