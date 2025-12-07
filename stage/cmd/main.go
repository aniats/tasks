package main

import (
	"fmt"
	"stage/stage"
	"time"
)

func createStage(name string) stage.Stage {
	return func(in stage.In) stage.Out {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for value := range in {

				time.Sleep(100 * time.Millisecond)
				out <- fmt.Sprintf("%s: %v", name, value)
			}
		}()
		return out
	}
}

func main() {
	input := make(chan interface{})
	done := make(chan interface{})

	stages := []stage.Stage{
		createStage("stage1"),
		createStage("stage2"),
		createStage("stage3"),
		createStage("stage4"),
	}

	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()

	output := stage.ExecutePipeline(input, done, stages...)

	count := 0
	for result := range output {
		count++
		fmt.Printf("[%d] %v\n", count, result)
	}
}
