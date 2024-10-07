package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/aranw/freq-demo/generator/internal/freqgen"
	"github.com/aranw/freq-demo/workflow"

	"github.com/sourcegraph/conc/pool"
	"go.temporal.io/sdk/client"
)

func NewGenerator() (*Generator, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("creating Temporal client: %w", err)
	}
	return &Generator{temporal: c}, nil
}

type Generator struct {
	temporal client.Client
}

func (g *Generator) Run(ctx context.Context, args []string) error {
	p := pool.New().
		WithContext(ctx).
		WithCancelOnError()

	bpool := pool.New().WithContext(ctx)

	// Batch variables
	batchSize := 1000
	batch := make([]workflow.FrequencyReading, 0, batchSize)

	p.Go(func(ctx context.Context) error {
		// Create a ticker that fires every 50 milliseconds
		ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				reading := freqgen.Generator()

				// Add the reading to the batch
				batch = append(batch, reading)

				// Check if batch is full
				if len(batch) == batchSize {
					// Process the batch
					bpool.Go(g.processBatch(batch))

					// Reset the batch
					batch = batch[:0]
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	p.Go(func(ctx context.Context) error {
		<-ctx.Done()

		g.temporal.Close()

		return bpool.Wait()
	})

	return p.Wait()
}

func (g *Generator) processBatch(batch []workflow.FrequencyReading) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// TODO: Ideally this shouldn't be here and should be in a separate package
		input := workflow.FrequencyBatch{
			Readings: batch,
		}

		options := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("frequency-batch-%d", rand.Int32()),
			TaskQueue: workflow.FrequencyBatchTaskQueueName,
		}

		slog.Info("processing batch", "id", options.ID)

		we, err := g.temporal.ExecuteWorkflow(context.Background(), options, workflow.ProcessFrequencyBatch, input)
		if err != nil {
			return fmt.Errorf("starting workflow: %w", err)
		}

		slog.Info("executed workflow", "workflow_id", we.GetID(), "run_id", we.GetRunID())

		var result workflow.FrequencyBatchResult

		if err := we.Get(context.Background(), &result); err != nil {
			return fmt.Errorf("getting workflow result: %w", err)
		}

		// Example: Print batch size and first/last readings
		fmt.Printf("Processing batch of %d readings\n", result.BatchSize)
		fmt.Printf("First reading: Time %s, Frequency = %.5f Hz\n",
			result.FirstReadingFrequencyTime.Format("15:04:05.000"), result.FirstReadingFrequency)
		fmt.Printf("Last reading:  Time %s, Frequency = %.5f Hz\n",
			result.LastReadingFrequencyTime.Format("15:04:05.000"), result.LastReadingFrequency)
		fmt.Printf("Average Frequency:     %.5f Hz\n", result.AverageFrequency)
		fmt.Printf("Minimum Frequency:     %.5f Hz\n", result.MinimumFrequency)
		fmt.Printf("Maximum Frequency:     %.5f Hz\n", result.MaximumFrequency)
		fmt.Printf("Standard Deviation:    %.5f Hz\n", result.StandardDeviation)
		return nil
	}
}
