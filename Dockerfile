FROM golang:1.24-alpine as builder
WORKDIR /app

RUN apk add --no-cache --update \
    gcc \
    musl-dev \
    g++ \
    wget

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /app/twir_application ./cmd/main.go

FROM alpine:3.21
LABEL org.opencontainers.image.authors="Satont <satontworldwide@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/twirapp/language-processor"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Language Processor"
LABEL org.opencontainers.image.description="Simple HTTP Server for translate texts and detect languages"
LABEL org.opencontainers.image.vendor="TwirApp"

WORKDIR /app
ADD --chmod=004 https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.bin .
RUN apk add --no-cache \
    libstdc++ \
    libgcc

COPY --from=builder /app/twir_application /bin/twir_application
CMD ["/bin/twir_application", "-modelpath", "/app/lid.176.bin"]