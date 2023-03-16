FROM golang:1.20.2-alpine3.17 as builder

WORKDIR /app

COPY go.* ./

RUN apk update && apk add --no-cache git

RUN go mod download

COPY *.go .

RUN go build -v -o ranutification-app


FROM alpine:3.17 as runner


RUN adduser -D runner | mkdir -p /app
RUN chown -R runner:runner $HOME
RUN chown -R runner:runner /app

USER runner


# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/ranutification-app /app/ranutification-app

EXPOSE 8080

CMD ["/app/ranutification-app"]
