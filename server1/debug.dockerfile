FROM golang:1.21.5-alpine

RUN apk add --no-cache go git
RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.0
ENV PATH="/root/go/bin:${PATH}"
RUN which dlv

WORKDIR /app
COPY . .
RUN go mod download
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV MODE=debug
ENV RABBITMQ_HOST=rabbitmq
RUN go build -gcflags="all=-N -l" -o main ./cmd/server

EXPOSE 8080 2345

CMD ["dlv", "exec", "./main", "--headless=true", "--accept-multiclient" ,"--listen=:2345", "--api-version=2"]
