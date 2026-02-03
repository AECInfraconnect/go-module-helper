package jaeger

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// LogrusAdapter adapts logrus to jaeger.Logger interface
type LogrusAdapter struct{}

func (l LogrusAdapter) Error(msg string) {
	logrus.Error(msg)
}

func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}

// Option configures the tracer
type Option func(*options)

type options struct {
	logger jaeger.Logger
}

// WithStdLogger uses standard logger
func WithStdLogger() Option {
	return func(o *options) {
		o.logger = jaeger.StdLogger
	}
}

// WithNullLogger uses null logger (no output)
func WithNullLogger() Option {
	return func(o *options) {
		o.logger = jaeger.NullLogger
	}
}

// WithLogrusLogger uses logrus adapter
func WithLogrusLogger() Option {
	return func(o *options) {
		o.logger = LogrusAdapter{}
	}
}

// Init returns an instance of Jaeger Tracer
func Init(serviceName string, opts ...Option) (opentracing.Tracer, io.Closer) {
	o := &options{
		logger: jaeger.NullLogger, // default
	}
	for _, opt := range opts {
		opt(o)
	}

	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	// Override with environment variables
	if envCfg, err := config.FromEnv(); err == nil {
		if envCfg.Sampler != nil {
			cfg.Sampler = envCfg.Sampler
		}
		if envCfg.Reporter != nil {
			cfg.Reporter = envCfg.Reporter
		}
	}

	cfg.ServiceName = serviceName

	tracer, closer, err := cfg.NewTracer(config.Logger(o.logger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
