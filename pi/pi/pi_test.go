package pi

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestNewCalculator(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid number of workers",
			numWorkers: 4,
			wantErr:    false,
		},
		{
			name:       "zero workers",
			numWorkers: 0,
			wantErr:    true,
			errMsg:     "number of workers must be positive, got 0",
		},
		{
			name:       "negative workers",
			numWorkers: -1,
			wantErr:    true,
			errMsg:     "number of workers must be positive, got -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := NewCalculator(tt.numWorkers)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, calc)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, calc)
				assert.Equal(t, tt.numWorkers, calc.numWorkers)
				assert.NotNil(t, calc.stopChan)
				assert.NotNil(t, calc.results)
				assert.NotNil(t, calc.wg)
			}
		})
	}
}

func TestCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name        string
		numWorkers  int
		computeTime time.Duration
		minResult   float64
		maxResult   float64
		maxDiff     float64
	}{
		{
			name:        "single worker",
			numWorkers:  1,
			computeTime: 50 * time.Millisecond,
			minResult:   2.0,
			maxResult:   4.0,
			maxDiff:     2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := NewCalculator(tt.numWorkers)
			assert.NoError(t, err)

			done := make(chan float64)
			go func() {
				result := calc.Calculate()
				done <- result
			}()

			time.Sleep(tt.computeTime)
			close(calc.stopChan)

			select {
			case result := <-done:
				assert.True(t, result > tt.minResult)
				assert.True(t, result < tt.maxResult)
				diff := math.Abs(result - math.Pi)
				assert.True(t, diff < tt.maxDiff)
			case <-time.After(5 * time.Second):
				t.Fatal("Test timed out")
			}
		})
	}
}

func TestResult(t *testing.T) {
	tests := []struct {
		name     string
		sum      float64
		expected float64
	}{
		{
			name:     "positive sum",
			sum:      3.14159,
			expected: 3.14159,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Result{Sum: tt.sum}
			assert.Equal(t, tt.expected, result.Sum)
		})
	}
}
