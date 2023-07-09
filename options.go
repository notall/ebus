package ebus

// Option represents the optional function.
type Option func(opts *_options)

func loadOptions(options ...Option) *_options {
	opts := new(_options)
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}
	return opts
}

// _options contains all options which will be applied when instantiating an ants pool.
type _options struct {
	// Logger is the customized logger for logging info, if it is not set,
	// default logger is used.
	Logger Logger
	// ThreadPool is the customized pool for executing handler, if it is not set,
	// default thread pool is used.
	ThreadPool ThreadPool
}

// WithLogger sets up a customized logger.
func WithLogger(logger Logger) Option {
	return func(opts *_options) {
		opts.Logger = logger
	}
}

// WithThreadPool sets up a customized ThreadPool.
func WithThreadPool(threadPool ThreadPool) Option {
	return func(opts *_options) {
		opts.ThreadPool = threadPool
	}
}
