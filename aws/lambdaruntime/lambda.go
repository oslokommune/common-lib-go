package lambdaruntime

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
)

func init() {
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
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackFieldName = "stacktrace"

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

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp["message"] = message

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

func localProxy[T any, R any](handler func(payload T) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t T
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		_ = decoder.Decode(&t)

		_, err := handler(t)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func localProxyWithContext[T any, R any](ctx context.Context, handler func(ctx context.Context, payload T) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t T
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		_ = decoder.Decode(&t)

		_, err := handler(ctx, t)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Start[T any, R any](handler func(payload T) (R, error)) {
	if !IsRunningAsLambda() {
		log.Info().Msg("starting web proxy for local execution")
		proxyHandler := http.HandlerFunc(localProxy[T, R](handler))
		http.Handle("/", proxyHandler)
		log.Fatal().Err(http.ListenAndServe(":8080", nil))
	} else {
		lambda.Start(handler)
	}
}

func StartWithContext[T any, R any](ctx context.Context, createHandler func(context.Context) func(context.Context, T) (R, error), tracing bool) {
	if !IsRunningAsLambda() {
		log.Info().Msg("starting web proxy for local execution")
		proxyHandler := http.HandlerFunc(localProxyWithContext[T, R](ctx, createHandler(ctx)))
		http.Handle("/", proxyHandler)
		log.Fatal().Err(http.ListenAndServe(":8080", nil))
	} else {
		if tracing {
			tp, err := xrayconfig.NewTracerProvider(ctx)
			if err != nil {
				log.Panic().Err(err).Msg("Error creating trace provider")
			}

			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(xray.Propagator{})

			defer func(ctx context.Context) {
				err := tp.Shutdown(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Error shutting down tracer provider")
				}
			}(ctx)

			lambda.Start(otellambda.InstrumentHandler(createHandler(ctx), xrayconfig.WithRecommendedOptions(tp)...))
		} else {
			lambda.Start(createHandler(ctx))
		}
	}
}

func IsRunningAsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
