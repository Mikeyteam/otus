package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline run with pipeline.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	if len(stages) > 1 {
		out = ExecutePipeline(in, done, stages[:len(stages)-1]...)
	}
	return worker(done, stages[len(stages)-1](out))
}

// worker write in channel.
func worker(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()
	return out
}
