package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/messgeplusio/httpq/pkg/queue"
	"github.com/messgeplusio/httpq/pkg/queue/pg"

	"github.com/messgeplusio/httpq/pkg/types"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer initializes the OpenTelemetry tracer
func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

// processJob processes an HTTP job
func processJob(ctx context.Context, job queue.Job[*types.HTTPRequest]) error {
	if job.Job == nil {
		return fmt.Errorf("job is nil")
	}
	req, err := http.NewRequestWithContext(ctx, job.Job.Method, job.Job.URL, strings.NewReader(job.Job.Body))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	for key, value := range job.Job.Headers {
		req.Header.Add(key, value)
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		if resp != nil {
			return fmt.Errorf("error making HTTP request: %v, status code: %d", err, resp.StatusCode)
		}
		return fmt.Errorf("error making HTTP request: %v", err)
	}

	return nil
}

func main() {
	connStr := "passw
	queue, err := pg.NewPostgreSQLQueue(connStr)
	if err != nil {
		log.Fatal(err)
	}

	go PollAndProcess(context.Background(), queue)
	log.Println("Service is running")
	select {}
}

// PollAndProcess polls the queue and processes jobs
func PollAndProcess(ctx context.Context, queue *pg.PostgreSQLQueue) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := queue.Dequeue(ctx)
			if err != nil {
				log.Printf("Error dequeuing job: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}
			err = processJob(ctx, job)
			if err != nil {
				log.Printf("Error processing job %s: %v\n", job.ID, err)
				queue.Requeue(ctx, job, 5*time.Minute)
				continue
			}

			err = queue.Acknowledge(ctx, job.ID)
			if err != nil {
				log.Printf("Error acknowledging job %s: %v\n", job.ID, err)
				queue.Requeue(ctx, job, 5*time.Minute)
				continue
			}

			log.Printf("Job %s completed successfully\n", job.ID)
		}
	}
}
