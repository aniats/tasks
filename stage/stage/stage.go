package stage

type (
	In  = <-chan interface{}
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	out := in
	for _, s := range stages {
		out = runStageWithDone(done, out, s)
	}
	return out
}

func runStageWithDone(done In, in In, stage Stage) Out {
	var stageOut Out

	// to make defer recover from panic right after calling the stage func
	// could be replaced with other func, but I decided to keep it like that
	func() {
		defer func() {
			if r := recover(); r != nil {
				stageOut = nil
			}
		}()
		stageOut = stage(in)
	}()

	out := make(chan interface{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
			close(out)
		}()

		if stageOut == nil {
			return
		}

		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}
