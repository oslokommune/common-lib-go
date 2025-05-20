module github.com/oslokommune/common-lib-go/aws

go 1.23.0

require (
	github.com/aws/aws-sdk-go-v2 v1.32.3
	github.com/aws/aws-sdk-go-v2/config v1.27.33
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.2
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.7.37
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.7
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.39.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.8
	github.com/aws/aws-sdk-go-v2/service/ecr v1.36.3
	github.com/aws/aws-sdk-go-v2/service/ecs v1.45.4
	github.com/aws/aws-sdk-go-v2/service/lambda v1.58.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.61.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.8
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.7
	github.com/aws/aws-sdk-go-v2/service/ssm v1.52.8
	github.com/aws/aws-sdk-go-v2/service/xray v1.27.7
	github.com/aws/smithy-go v1.22.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-contrib/logger v1.1.2
	github.com/gin-gonic/gin v1.10.0
	github.com/oslokommune/common-lib-go/httpcomm v0.2.3
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.9.0
	github.com/swaggest/openapi-go v0.2.53
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.54.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.54.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.54.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0
	go.opentelemetry.io/contrib/propagators/aws v1.29.0
	go.opentelemetry.io/otel v1.29.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0
	go.opentelemetry.io/otel/sdk v1.29.0
	go.opentelemetry.io/otel/trace v1.29.0
)

require (
	github.com/bytedance/sonic v1.12.2 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/swaggest/jsonschema-go v0.3.72 // indirect
	github.com/swaggest/refl v1.3.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.54.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/arch v0.10.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.32 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.13 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.18
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.25.8
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.22.8
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.8
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.7 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws v0.54.0
	golang.org/x/sys v0.30.0 // indirect
)
