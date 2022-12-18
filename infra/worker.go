package infra

type Worker[T any] struct {
	workerNum      int
	producerStream chan T
	consumerStream chan T
	doneStreams    []chan bool
}

func NewWorker[T any](workerNum int) *Worker[T] {
	doneStreams := make([]chan bool, workerNum)
	for i := 0; i < workerNum; i++ {
		doneStreams[i] = make(chan bool)
	}

	return &Worker[T]{
		workerNum:      workerNum,
		producerStream: make(chan T),
		consumerStream: make(chan T),
		doneStreams:    doneStreams,
	}
}

func (w *Worker[T]) Produce(producer func(producerStream chan<- T, stop func())) {
	go func() {
		producer(w.producerStream, func() {
			close(w.producerStream)

			for _, doneStream := range w.doneStreams {
				<-doneStream

				close(doneStream)
			}

			close(w.consumerStream)
		})
	}()
}

func (w *Worker[T]) Consume(consumer func(val T)) {
	for i := 0; i < w.workerNum; i++ {
		go func(idx int) {
			doneStream := w.doneStreams[idx]

			for val := range w.producerStream {
				consumer(val)
			}

			doneStream <- true
		}(i)
	}
}

func (w *Worker[T]) Wait(waitFn func(data T)) {
	for val := range w.consumerStream {
		if waitFn != nil {
			waitFn(val)
		}
	}
}
