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
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp["message"] = message

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

func localProxy[T any, R any](handler func(ctx context.Context, payload T) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t T
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		_ = decoder.Decode(&t)

		_, err := handler(context.Background(), t)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Start[T any, R any](ctx context.Context, handler func(ctx context.Context, payload T) (R, error), tracing bool) {
	if !IsRunningAsLambda() {
		log.Info().Msg("starting web proxy for local execution")
		proxyHandler := http.HandlerFunc(localProxy[T, R](handler))
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

			lambda.Start(otellambda.InstrumentHandler(handler, xrayconfig.WithRecommendedOptions(tp)...))
		} else {
			lambda.Start(handler)
		}
	}
}

func IsRunningAsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
