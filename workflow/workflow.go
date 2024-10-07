package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessFrequencyBatch(ctx workflow.Context, input FrequencyBatch) (FrequencyBatchResult, error) {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    500, // 0 is unlimited retries
		// NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	var minFreq float64
	if err := workflow.ExecuteActivity(ctx, Min, input).Get(ctx, &minFreq); err != nil {
		return FrequencyBatchResult{}, fmt.Errorf("Min: failed to determine min frequency: %w", err)
	}

	var maxFreq float64
	if err := workflow.ExecuteActivity(ctx, Max, input).Get(ctx, &maxFreq); err != nil {
		return FrequencyBatchResult{}, fmt.Errorf("Max: failed to determine max frequency: %w", err)
	}

	var avgFreq float64
	if err := workflow.ExecuteActivity(ctx, Avg, input).Get(ctx, &avgFreq); err != nil {
		return FrequencyBatchResult{}, fmt.Errorf("Avg: failed to determine average frequency: %w", err)
	}

	var stdDev float64
	if err := workflow.ExecuteActivity(ctx, StdDev, input).Get(ctx, &stdDev); err != nil {
		return FrequencyBatchResult{}, fmt.Errorf("StdDev: failed to determine standard deviation for frequency batch: %w", err)
	}

	return FrequencyBatchResult{
		BatchSize:                 len(input.Readings),
		FirstReadingFrequency:     input.Readings[0].Frequency,
		FirstReadingFrequencyTime: input.Readings[0].Time,
		LastReadingFrequency:      input.Readings[len(input.Readings)-1].Frequency,
		LastReadingFrequencyTime:  input.Readings[len(input.Readings)-1].Time,
		AverageFrequency:          avgFreq,
		MinimumFrequency:          minFreq,
		MaximumFrequency:          maxFreq,
		StandardDeviation:         stdDev,
	}, nil
}
