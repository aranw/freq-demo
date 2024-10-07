package freqgen

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/aranw/freq-demo/workflow"
)

// Ornstein-Uhlenbeck parameters
const (
	mu    = 50.0 // Mean frequency (Hz)
	theta = 0.1  // Speed of mean reversion
	sigma = 0.05 // Volatility (Hz)
	dt    = 0.05 // Time step in seconds (50 ms)
)

func Generator() workflow.FrequencyReading {
	// Generate a standard normal random variable
	eps := rand.NormFloat64()

	// Initial frequency set to mean
	X := mu

	// Update frequency using the Ornstein-Uhlenbeck formula
	dX := theta*(mu-X)*dt + sigma*math.Sqrt(dt)*eps
	X += dX

	// Get the current timestamp
	timestamp := time.Now()

	// Create a FrequencyReading
	return workflow.FrequencyReading{
		Time:      timestamp,
		Frequency: X,
	}
}
