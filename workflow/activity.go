package workflow

import (
	"context"
	"math"
)

func Min(ctx context.Context, data FrequencyBatch) (float64, error) {
	var sum float64
	minFreq := data.Readings[0].Frequency

	frequencies := make([]float64, len(data.Readings))

	for i, reading := range data.Readings {
		freq := reading.Frequency
		sum += freq
		frequencies[i] = freq

		if freq < minFreq {
			minFreq = freq
		}
	}

	return minFreq, nil
}

func Max(ctx context.Context, data FrequencyBatch) (float64, error) {
	var sum float64
	maxFreq := data.Readings[0].Frequency

	frequencies := make([]float64, len(data.Readings))

	for i, reading := range data.Readings {
		freq := reading.Frequency
		sum += freq
		frequencies[i] = freq

		if freq > maxFreq {
			maxFreq = freq
		}
	}

	return maxFreq, nil
}

func Avg(ctx context.Context, data FrequencyBatch) (float64, error) {
	var sum float64
	for _, reading := range data.Readings {
		freq := reading.Frequency
		sum += freq
	}

	return sum / float64(len(data.Readings)), nil
}

func StdDev(ctx context.Context, data FrequencyBatch) (float64, error) {
	frequencies := make([]float64, len(data.Readings))

	var sum float64
	for i, reading := range data.Readings {
		freq := reading.Frequency
		sum += freq
		frequencies[i] = freq
	}

	avgFreq := sum / float64(len(data.Readings))

	var sumSqDiff float64
	for _, freq := range frequencies {
		diff := freq - avgFreq
		sumSqDiff += diff * diff
	}
	variance := sumSqDiff / float64(len(frequencies))

	return math.Sqrt(variance), nil
}
