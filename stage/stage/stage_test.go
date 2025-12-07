package stage

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecutePipeline(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		stages   []Stage
		expected []interface{}
	}{
		{
			name:     "empty_stages",
			input:    []interface{}{1, 2, 3},
			stages:   []Stage{},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:  "single_stage",
			input: []interface{}{1, 2, 3},
			stages: []Stage{
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- value.(int) * 2
						}
					}()
					return out
				},
			},
			expected: []interface{}{2, 4, 6},
		},
		{
			name:  "multiple_stages",
			input: []interface{}{1, 2, 3},
			stages: []Stage{
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- value.(int) * 2
						}
					}()
					return out
				},
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- value.(int) + 1
						}
					}()
					return out
				},
			},
			expected: []interface{}{3, 5, 7},
		},
		{
			name:  "concurrent_processing",
			input: []interface{}{1, 2, 3, 4, 5},
			stages: []Stage{
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- fmt.Sprintf("stage1: %v", value)
						}
					}()
					return out
				},
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- fmt.Sprintf("stage2: %v", value)
						}
					}()
					return out
				},
			},
			expected: []interface{}{
				"stage2: stage1: 1",
				"stage2: stage1: 2",
				"stage2: stage1: 3",
				"stage2: stage1: 4",
				"stage2: stage1: 5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				in := make(chan interface{})
				done := make(chan interface{})

				go func() {
					defer close(in)
					for _, value := range tt.input {
						in <- value
					}
				}()

				out := ExecutePipeline(in, done, tt.stages...)

				var results []interface{}
				for value := range out {
					results = append(results, value)
				}

				require.Len(t, results, len(tt.expected), "unexpected number of results")
				require.Equal(t, tt.expected, results, "results don't match expected")
			})
		})
	}
}

func TestExecutePipelineWithDone(t *testing.T) {
	tests := []struct {
		name           string
		input          []interface{}
		stages         []Stage
		closeDoneEarly bool
		expectedResult interface{}
		expectedTime   time.Duration
	}{
		{
			name:  "timing_precision",
			input: []interface{}{42},
			stages: []Stage{
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							time.Sleep(10 * time.Millisecond)
							out <- value
						}
					}()
					return out
				},
			},
			closeDoneEarly: false,
			expectedResult: 42,
			expectedTime:   10 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()

				in := make(chan interface{})
				done := make(chan interface{})

				if tt.closeDoneEarly {
					close(done)
				}

				go func() {
					defer close(in)
					for _, value := range tt.input {
						in <- value
					}
				}()

				out := ExecutePipeline(in, done, tt.stages...)

				if tt.closeDoneEarly {
					_, ok := <-out
					require.False(t, ok, "pipeline should be closed when done is signaled")
				} else {
					result := <-out
					elapsed := time.Since(start)

					require.Equal(t, tt.expectedResult, result, "unexpected result value")
					require.Equal(t, tt.expectedTime, elapsed, "timing should be precise in synctest")
				}
			})
		})
	}
}

func TestExecutePipelineWithPanic(t *testing.T) {
	tests := []struct {
		name         string
		stages       []Stage
		expectClosed bool
	}{
		{
			name: "panic_in_first_stage",
			stages: []Stage{
				func(in In) Out {
					panic("first stage panic")
				},
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- value.(int) + 1
						}
					}()
					return out
				},
			},
			expectClosed: true,
		},
		{
			name: "panic_in_second_stage",
			stages: []Stage{
				func(in In) Out {
					out := make(chan interface{})
					go func() {
						defer close(out)
						for value := range in {
							out <- value.(int) * 2
						}
					}()
					return out
				},
				func(in In) Out {
					panic("second stage panic")
				},
			},
			expectClosed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				in := make(chan interface{})
				done := make(chan interface{})

				close(in)

				out := ExecutePipeline(in, done, tt.stages...)

				_, ok := <-out
				require.Equal(t, tt.expectClosed, !ok, "unexpected channel state after panic")
			})
		})
	}
}
