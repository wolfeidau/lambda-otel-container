# lambda-otel-container

This project demonstrates the use of [OpenTelemetry](https://opentelemetry.io/) with an [AWS Lambda](https://aws.amazon.com/lambda/) function written in [Go](https://go.dev) and running in a container.

# Overview

This project is built using the [opentelemetry-lambda](https://github.com/open-telemetry/opentelemetry-lambda) via an extension to send traces to [honeycomb](https://honeycomb.io). 

Rather than using the lambda layer, I am building [collector](https://github.com/open-telemetry/opentelemetry-lambda/tree/main/collector) from source in the [Dockerfile](./Dockerfile).

Outside of these changes the project is the same as the examples on [github.com/awslabs/aws-lambda-go-api-proxy](https://github.com/awslabs/aws-lambda-go-api-proxy), have a look at the [main.goo](./main.go) for more details.

# Deployment

To deploy this project I have included a [Makefile](./Makefile) with a `deploy` target, you should have a read over this to understand what it does.

It assumes you have a few environment variables set, I use [direnv](https://direnv.net/) which can be configured to load these from a `.envrc` file.

```
#!/usr/bin/env bash

export HONEYCOMB_API_KEY=<YOUR API KEY HERE>
export HONEYCOMB_DATASET=dev
export HONEYCOMB_ENDPOINT=api.honeycomb.io:443
```


# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).