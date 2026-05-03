package util

import (
	"context"
	"errors"
	"runtime"
	"sync"
)

var (
	ErrNilPool    = errors.New("util: pool is nil")
	ErrNilTask    = errors.New("util: task is nil")
	ErrPoolClosed = errors.New("util: pool is closed")
)

type (
	poolConfig struct {
		maxWorkers int
	}
	runConfig struct {
		workers int
	}

	// PoolOption configures a pool at creation time.
	PoolOption func(*poolConfig)
	// RunOption configures a single batch execution.
	RunOption func(*runConfig)

	// Pool is a shared goroutine control center.
	Pool struct {
		maxWorkers int
		jobs       chan func()
		done       chan struct{}
		closeOnce  sync.Once
		workersWg  sync.WaitGroup
	}

	// Result is the async result of a single task.
	Result[T any] struct {
		Value T
		Err   error
	}

	// BatchResult is the result of one item in a batch task.
	BatchResult[I any, O any] struct {
		Index int
		Input I
		Value O
		Err   error
	}

	// Future wraps the result of a single async task.
	Future[T any] struct {
		ch chan Result[T]
	}

	// BatchFuture wraps the result stream of an async batch task.
	BatchFuture[I any, O any] struct {
		ch      chan BatchResult[I, O]
		mu      sync.Mutex
		results []BatchResult[I, O]
	}
)

// WithMaxWorkers limits the total worker goroutines owned by the pool.
func WithMaxWorkers(n int) PoolOption {
	return func(cfg *poolConfig) {
		if n > 0 {
			cfg.maxWorkers = n
		}
	}
}

// WithWorkers limits worker usage for a single batch execution.
func WithWorkers(n int) RunOption {
	return func(cfg *runConfig) {
		if n > 0 {
			cfg.workers = n
		}
	}
}

// NewPool creates a configurable goroutine control center.
func NewPool(opts ...PoolOption) *Pool {
	cfg := poolConfig{
		maxWorkers: defaultMaxWorkers(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.maxWorkers <= 0 {
		cfg.maxWorkers = 1
	}

	pool := &Pool{
		maxWorkers: cfg.maxWorkers,
		jobs:       make(chan func()),
		done:       make(chan struct{}),
	}

	pool.workersWg.Add(pool.maxWorkers)
	for i := 0; i < pool.maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

// MaxWorkers returns the configured upper bound of worker goroutines.
func (p *Pool) MaxWorkers() int {
	if p == nil {
		return 0
	}
	return p.maxWorkers
}

// Close stops the pool from accepting new tasks and waits for workers to exit.
func (p *Pool) Close() {
	if p == nil {
		return
	}
	p.closeOnce.Do(func() {
		close(p.done)
		p.workersWg.Wait()
	})
}

// OK reports whether the single-task result succeeded.
func (r Result[T]) OK() bool {
	return r.Err == nil
}

// OK reports whether the batch item succeeded.
func (r BatchResult[I, O]) OK() bool {
	return r.Err == nil
}

// Await blocks until the async task finishes.
func (f *Future[T]) Await() Result[T] {
	if f == nil || f.ch == nil {
		return Result[T]{Err: ErrNilTask}
	}
	result, ok := <-f.ch
	if !ok {
		return Result[T]{Err: ErrNilTask}
	}
	return result
}

// Next returns the next completed batch result in completion order.
func (f *BatchFuture[I, O]) Next() (BatchResult[I, O], bool) {
	if f == nil || f.ch == nil {
		return BatchResult[I, O]{}, false
	}

	result, ok := <-f.ch
	if !ok {
		return BatchResult[I, O]{}, false
	}

	f.record(result)
	return result, true
}

// Await waits for all batch items to finish and returns results in input order.
func (f *BatchFuture[I, O]) Await() []BatchResult[I, O] {
	if f == nil {
		return nil
	}
	for {
		_, ok := f.Next()
		if !ok {
			return f.snapshot()
		}
	}
}

// Submit runs a single async task through the pool.
func Submit[T any](pool *Pool, ctx context.Context, task func(context.Context) (T, error)) *Future[T] {
	future := &Future[T]{
		ch: make(chan Result[T], 1),
	}

	if pool == nil {
		future.resolve(Result[T]{Err: ErrNilPool})
		return future
	}
	if task == nil {
		future.resolve(Result[T]{Err: ErrNilTask})
		return future
	}

	ctx = normalizeContext(ctx)
	if err := pool.submit(func() {
		value, runErr := task(ctx)
		future.resolve(Result[T]{Value: value, Err: runErr})
	}); err != nil {
		future.resolve(Result[T]{Err: err})
	}

	return future
}

// Call runs a single task synchronously through the pool.
func Call[T any](pool *Pool, ctx context.Context, task func(context.Context) (T, error)) (T, error) {
	result := Submit(pool, ctx, task).Await()
	return result.Value, result.Err
}

// Map runs a batch task synchronously and returns ordered results.
func Map[I any, O any](pool *Pool, ctx context.Context, items []I, task func(context.Context, I) (O, error), opts ...RunOption) []BatchResult[I, O] {
	return MapAsync(pool, ctx, items, task, opts...).Await()
}

// MapAsync runs a batch task asynchronously.
func MapAsync[I any, O any](pool *Pool, ctx context.Context, items []I, task func(context.Context, I) (O, error), opts ...RunOption) *BatchFuture[I, O] {
	future := newBatchFuture[I, O](len(items))

	if len(items) == 0 {
		future.close()
		return future
	}
	if pool == nil {
		completeBatchWithError(future, items, 0, ErrNilPool)
		future.close()
		return future
	}
	if task == nil {
		completeBatchWithError(future, items, 0, ErrNilTask)
		future.close()
		return future
	}

	ctx = normalizeContext(ctx)
	workers := resolveWorkers(pool.maxWorkers, len(items), opts...)

	go func() {
		slots := make(chan struct{}, workers)
		var wg sync.WaitGroup

		defer func() {
			wg.Wait()
			future.close()
		}()

		for index, item := range items {
			select {
			case <-ctx.Done():
				completeBatchWithError(future, items, index, ctx.Err())
				return
			case slots <- struct{}{}:
			}

			wg.Add(1)
			idx := index
			input := item

			err := pool.submit(func() {
				defer func() {
					<-slots
					wg.Done()
				}()

				value, runErr := task(ctx, input)
				future.push(BatchResult[I, O]{
					Index: idx,
					Input: input,
					Value: value,
					Err:   runErr,
				})
			})
			if err != nil {
				wg.Done()
				<-slots
				future.push(BatchResult[I, O]{
					Index: idx,
					Input: input,
					Err:   err,
				})
				completeBatchWithError(future, items, idx+1, err)
				return
			}
		}
	}()

	return future
}

func (p *Pool) worker() {
	defer p.workersWg.Done()

	for {
		select {
		case <-p.done:
			return
		case job := <-p.jobs:
			if job != nil {
				job()
			}
		}
	}
}

func (p *Pool) submit(job func()) error {
	if p == nil {
		return ErrNilPool
	}

	select {
	case <-p.done:
		return ErrPoolClosed
	case p.jobs <- job:
		return nil
	}
}

func (f *Future[T]) resolve(result Result[T]) {
	if f == nil || f.ch == nil {
		return
	}
	f.ch <- result
	close(f.ch)
}

func (f *BatchFuture[I, O]) push(result BatchResult[I, O]) {
	if f == nil || f.ch == nil {
		return
	}
	f.ch <- result
}

func (f *BatchFuture[I, O]) close() {
	if f == nil || f.ch == nil {
		return
	}
	close(f.ch)
}

func (f *BatchFuture[I, O]) record(result BatchResult[I, O]) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if result.Index < 0 || result.Index >= len(f.results) {
		return
	}
	f.results[result.Index] = result
}

func (f *BatchFuture[I, O]) snapshot() []BatchResult[I, O] {
	f.mu.Lock()
	defer f.mu.Unlock()

	cloned := make([]BatchResult[I, O], len(f.results))
	copy(cloned, f.results)
	return cloned
}

func newBatchFuture[I any, O any](size int) *BatchFuture[I, O] {
	bufferSize := size
	if bufferSize <= 0 {
		bufferSize = 1
	}

	return &BatchFuture[I, O]{
		ch:      make(chan BatchResult[I, O], bufferSize),
		results: make([]BatchResult[I, O], size),
	}
}

func completeBatchWithError[I any, O any](future *BatchFuture[I, O], items []I, start int, err error) {
	for i := start; i < len(items); i++ {
		future.push(BatchResult[I, O]{
			Index: i,
			Input: items[i],
			Err:   err,
		})
	}
}

func resolveWorkers(maxWorkers int, itemCount int, opts ...RunOption) int {
	cfg := runConfig{
		workers: maxWorkers,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.workers <= 0 {
		cfg.workers = maxWorkers
	}
	if cfg.workers > maxWorkers {
		cfg.workers = maxWorkers
	}
	if cfg.workers > itemCount {
		cfg.workers = itemCount
	}
	if cfg.workers <= 0 {
		cfg.workers = 1
	}

	return cfg.workers
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func defaultMaxWorkers() int {
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return 1
}
