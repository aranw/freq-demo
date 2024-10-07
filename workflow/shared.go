package workflow

import "time"

const FrequencyBatchTaskQueueName = "FREQUENCY_BATCH_TASK_QUEUE"

type FrequencyBatch struct {
	Readings []FrequencyReading
}

// FrequencyReading holds a timestamp and a frequency value
type FrequencyReading struct {
	Time      time.Time
	Frequency float64
}

type FrequencyBatchResult struct {
	BatchSize                 int
	FirstReadingFrequency     float64
	FirstReadingFrequencyTime time.Time
	LastReadingFrequency      float64
	LastReadingFrequencyTime  time.Time
	AverageFrequency          float64
	MinimumFrequency          float64
	MaximumFrequency          float64
	StandardDeviation         float64
}

func (f *FrequencyBatchResult) String() string {
	return ""
}
