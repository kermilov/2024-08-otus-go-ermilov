package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		doneProxy := func(in In) Out {
			out := stage(in)
			outProxy := make(Bi)
			go func() {
				defer close(outProxy)
				for v := range out {
					select {
					case <-done:
						return
					case outProxy <- v:
					}
				}
			}()
			return outProxy
		}
		in = doneProxy(in)
	}
	return in
}
