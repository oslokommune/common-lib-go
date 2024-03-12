package ginruntime

import (
	"context"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/oslokommune/common-lib-go/aws/ginruntime/openapi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type GinEngine struct {
	ctx        context.Context
	engine     *gin.Engine
	tp         *trace.TracerProvider
	propagator propagation.TextMapPropagator
	openapi    *openapi.OpenAPI
	onShutdown []func()
}

func New(ctx context.Context, options ...Option) *GinEngine {

	configureLogging()

	// Creates a router without any middleware by default
	engine := gin.New()

	// Do not encode path
	engine.UseRawPath = true

	// Global middleware
	engine.Use(ErrorHandler())

	var rxURL = regexp.MustCompile(`^/*`)
	// Use zerolog for logging and turn off access logging for all paths
	engine.Use(logger.SetLogger(logger.WithSkipPathRegexps(rxURL)))

	// CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"authorization", "content-type"}
	corsConfig.AddAllowMethods("OPTIONS")
	engine.Use(cors.New(corsConfig))

	// Recover from panics
	engine.Use(RecoveryMiddleware)

	e := &GinEngine{ctx, engine, nil, nil, nil, make([]func(), 0)}
	for _, option := range options {
		if option.openapi != nil {
			e.enableOpenAPI(option.openapi)
		}
	}
	return e
}

func (e *GinEngine) OnShutdown(f func()) {
	e.onShutdown = append(e.onShutdown, f)
}

func (e *GinEngine) shutdownCallbacks() {
	for _, f := range e.onShutdown {
		defer f()
	}
	e.onShutdown = []func(){}
}

func IsRunningAsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func configureLogging() {
	// Loglevel defaults to info, can be overriden
	logLevel := zerolog.InfoLevel
	if l := os.Getenv("LOG_LEVEL"); l != "" {
		switch l {
		case "DEBUG":
			logLevel = zerolog.DebugLevel
		case "INFO":
			logLevel = zerolog.InfoLevel
		case "ERROR":
			logLevel = zerolog.ErrorLevel
		case "TRACE":
			logLevel = zerolog.TraceLevel
		}
	}
	zerolog.SetGlobalLevel(logLevel)

	// Setup logger to correspond with logstash format
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = StackTraceMarshaller
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackFieldName = "stack_trace"

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}

	// Add default fields to std logger
	log.Logger = log.With().Str("app_label", os.Getenv("APP_LABEL")).Caller().Logger()

	// If running locally, do not format as json
	if !IsRunningAsLambda() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
