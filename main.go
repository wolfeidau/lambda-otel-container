package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello")
}

func main() {
	ctx := context.Background()

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Create and start new OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
		// otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalf("failed to create new OTLP trace exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tp)

	handler := http.HandlerFunc(home)
	http.Handle("/", otelhttp.NewHandler(handler, "home"))

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
