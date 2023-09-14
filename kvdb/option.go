package kvdb

type Option func(*options)

type options struct {
	dbFileName string
}

var defaultOptions = options{
	dbFileName: "kvdb.data",
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newOptions(opts ...Option) *options {
	options := &options{
		dbFileName: defaultOptions.dbFileName,
	}
	options.apply(opts...)
	return options
}

func WithDBFileName(name string) Option {
	return func(o *options) {
		o.dbFileName = name
	}
}
