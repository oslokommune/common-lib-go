module github.com/oslokommune/common-lib-go/aws

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.26.2
	github.com/aws/aws-sdk-go-v2/config v1.27.14
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.16
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.7.16
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.38.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.35.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.32.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.41.9
	github.com/aws/aws-sdk-go-v2/service/lambda v1.54.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.54.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.28.8
	github.com/aws/aws-sdk-go-v2/service/sns v1.29.6
	github.com/aws/aws-sdk-go-v2/service/ssm v1.50.2
	github.com/aws/aws-sdk-go-v2/service/xray v1.25.6
	github.com/aws/smithy-go v1.20.2
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-contrib/logger v1.1.2
	github.com/gin-gonic/gin v1.10.0
	github.com/oslokommune/common-lib-go/httpcomm v0.2.2
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.9.0
	github.com/swaggest/openapi-go v0.2.50
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.51.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.51.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.51.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.opentelemetry.io/contrib/propagators/aws v1.26.0
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0
	go.opentelemetry.io/otel/sdk v1.26.0
	go.opentelemetry.io/otel/trace v1.26.0
)

require (
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/swaggest/jsonschema-go v0.3.70 // indirect
	github.com/swaggest/refl v1.3.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.51.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240515191416-fc5f0ca64291 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.19
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.23.8
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.20.6
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.32.1
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.24.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.8 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws v0.51.0
	golang.org/x/sys v0.20.0 // indirect
)
