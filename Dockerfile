FROM golang:1.23-alpine AS build_base
RUN apk add --no-cache git make
WORKDIR /src/lambda-otel-container
ARG APP_VERSION

COPY . /src/lambda-otel-container
WORKDIR /src/lambda-otel-container
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 installsuffix=cgo go build -ldflags "-s -w -X 'main.version=$APP_VERSION'" -trimpath -o bootstrap .

RUN git clone https://github.com/open-telemetry/opentelemetry-lambda /src/opentelemetry-lambda
WORKDIR /src/opentelemetry-lambda/collector
RUN BUILDTAGS="lambdacomponents.custom,lambdacomponents.receiver.otlp,lambdacomponents.processor.all,lambdacomponents.exporter.otlp,lambdacomponents.connector.spanmetrics" make build

FROM alpine:3.21
RUN apk add ca-certificates
ENV OTEL_TRACES_SAMPLER=always_on
ENV OTEL_LAMBDA_DISABLE_AWS_CONTEXT_PROPAGATION=true
COPY --from=build_base /src/lambda-otel-container/bootstrap /app/bootstrap
COPY --from=build_base /src/opentelemetry-lambda/collector/build/extensions/collector /opt/extensions/collector
COPY adot-config.yaml /opt/collector-config/config.yaml
ENTRYPOINT ["/app/bootstrap"]