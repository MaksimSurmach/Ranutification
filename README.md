# Ranutification

Ranutification is a lightweight notification service that allows you to send notifications to Telegram users via a simple webhook.

## Introduction

Ranutification is designed to receive POST or GET requests at a webhook endpoint and relay the request body to a specified Telegram user. This makes it easy to integrate notifications into your applications or services.

## Getting Started

You can get Ranutification from DockerHub using the `maxsurm/ranutif:latest` image. To configure it, make sure to set the `TELEGRAM_API_TOKEN` environment variable with your Telegram bot token.

Alternatively, you can build Ranutification from source.

### Prerequisites

To build from source, make sure to download the project's Go module dependencies:

```go mod download```


## Usage

To run Ranutification, execute the following command:

```go run main.go```

Alternatively, you can build the binary file with the following command:

```go build main.go```

Then, run the binary file to start the service.

## Contributing

We welcome contributions to Ranutification. To contribute, follow these steps:

1. Fork the project.
2. Create a feature branch for your changes.
3. Send us a pull request.

## License

Ranutification is released under the MIT License. For more details, see the LICENSE.md file.