FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./configs/api/ ./configs/api/
COPY ./docs ./docs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/api

FROM scratch

COPY ./docs ./docs
COPY ./configs/api/ ./configs/api/
COPY --from=builder /app/api /usr/bin/
EXPOSE 8080

ENTRYPOINT ["api"]
