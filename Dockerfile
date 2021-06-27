FROM golang:alpine as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g' /etc/apk/repositories
RUN apk add build-base
WORKDIR /go/src/news
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' ./...

FROM scratch

COPY --from=builder /go/src/news/news /
COPY --from=builder /go/src/news/config/config.yaml.example /config/config.yaml
COPY --from=builder /go/src/news/templates /templates
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Shanghai
ENV ZONEINFO=/zoneinfo.zip

EXPOSE "9999"
ENTRYPOINT ["/news"]
