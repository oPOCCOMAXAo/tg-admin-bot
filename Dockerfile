FROM golang:1.24-alpine3.20 AS builder
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go \
	build \ 
	-trimpath \
	-o /bin/app \
	cmd/app/main.go

FROM alpine:3.20
RUN apk add --no-cache tzdata
COPY --from=builder /bin/app /bin/app
RUN mkdir -p /data
ENV SERVER_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/bin/app"]
