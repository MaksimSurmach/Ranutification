FROM golang:1.21.4-alpine3.17 as builder

ENV GO111MODULE=on

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.* ./

RUN go mod download

COPY . ./

RUN go build -v -o ranutification

FROM golang:1.21.4-alpine3.17 as runner

LABEL MAINTAINER Author <maxsurm@gmail.com>

RUN adduser -D runner | mkdir -p /app
RUN chown -R runner:runner $HOME
RUN chown -R runner:runner /app

USER runner

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/ranutification /app/ranutification

EXPOSE 8080

CMD ["/app/ranutification"]
