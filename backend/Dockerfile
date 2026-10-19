FROM golang:1.25-alpine AS build

ENV GOPROXY=https://mirrors.aliyun.com/goproxy,direct \
    GOSUMDB=sum.golang.google.cn

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/ai-for-api ./cmd/server

FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app && apk add --no-cache ca-certificates tzdata wget
WORKDIR /app
COPY --from=build /out/ai-for-api /app/ai-for-api
USER app
EXPOSE 8080
ENTRYPOINT ["/app/ai-for-api"]
