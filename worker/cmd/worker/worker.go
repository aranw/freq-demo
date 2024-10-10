package main

import (
	"context"
	"log"

	"github.com/aranw/freq-demo/workflow"

	"github.com/sourcegraph/conc/pool"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func NewWorker() *Worker {
	return &Worker{}
}

type Worker struct{}

func (w *Worker) Run(ctx context.Context, args []string) error {
	// TODO: Load config here

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	worker := worker.New(c, workflow.FrequencyBatchTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	worker.RegisterWorkflow(workflow.ProcessFrequencyBatch)
	worker.RegisterActivity(workflow.Min)
	worker.RegisterActivity(workflow.Max)
	worker.RegisterActivity(workflow.Avg)
	worker.RegisterActivity(workflow.StdDev)

	p := pool.New().
		WithContext(ctx).
		WithCancelOnError()

	p.Go(func(ctx context.Context) error {
		return worker.Run(nil)
	})

	p.Go(func(ctx context.Context) error {
		<-ctx.Done()

		worker.Stop()

		return nil
	})

	// Start listening to the Task Queue.
	return p.Wait()
}
