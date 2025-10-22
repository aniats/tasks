package pi

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"testing/synctest"
)

func TestNew(t *testing.T) {
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
			calc, err := New(tt.numWorkers)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, calc)
				require.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, calc)
				require.Equal(t, tt.numWorkers, calc.numWorkers)
				require.NotNil(t, calc.stopChan)
				require.NotNil(t, calc.results)
				require.NotNil(t, calc.wg)
			}
		})
	}
}

func TestCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		minResult  float64
		maxResult  float64
		maxDiff    float64
	}{
		{
			name:       "single worker",
			numWorkers: 1,
			minResult:  0.0,
			maxResult:  4.0,
			maxDiff:    4.0,
		},
		{
			name:       "multiple workers",
			numWorkers: 4,
			minResult:  0.0,
			maxResult:  4.0,
			maxDiff:    4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				calc, err := New(tt.numWorkers)
				require.NoError(t, err)

				done := make(chan float64, 1)
				go func() {
					result := calc.Calculate()
					done <- result
				}()

				calc.Stop()

				result := <-done
				require.True(t, result >= tt.minResult, "Result %f should be >= %f", result, tt.minResult)
				require.True(t, result <= tt.maxResult, "Result %f should be <= %f", result, tt.maxResult)
				diff := math.Abs(result - math.Pi)
				require.True(t, diff < tt.maxDiff, "Difference %f should be < %f", diff, tt.maxDiff)
			})
		})
	}
}

func TestCalculateLeibnizTerm(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected float64
	}{
		{
			name:     "first term (n=0)",
			n:        0,
			expected: 1.0,
		},
		{
			name:     "second term (n=1)",
			n:        1,
			expected: -1.0 / 3.0,
		},
		{
			name:     "third term (n=2)",
			n:        2,
			expected: 1.0 / 5.0,
		},
		{
			name:     "fourth term (n=3)",
			n:        3,
			expected: -1.0 / 7.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateLeibnizTerm(tt.n)
			require.InDelta(t, tt.expected, result, 0.0001, "Term calculation incorrect")
		})
	}
}

func TestCalculator_Stop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calc, err := New(2)
		require.NoError(t, err)

		done := make(chan float64, 1)
		go func() {
			result := calc.Calculate()
			done <- result
		}()

		calc.Stop()

		result := <-done
		require.True(t, result >= 0, "Result should be non-negative")
	})
}
