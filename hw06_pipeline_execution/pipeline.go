package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func worker(done In, out Out, stage Stage) Out {
	ch := make(Bi)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case val, ok := <-out:
				if !ok {
					return
				}
				ch <- val
			}
		}
	}()
	return stage(ch)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = worker(done, out, stage)
	}
	return out
}
