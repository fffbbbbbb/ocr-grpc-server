FROM golang:alpine AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ocr-grpc-server

FROM alpine


RUN apk update \ 
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk add --no-cache tesseract-ocr

WORKDIR /app

COPY --from=builder /build /app

ENTRYPOINT ["/app/ocr-grpc-server"]