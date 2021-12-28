FROM golang:alpine as builder

#RUN sed -i 's/deb.debian.org/mirrors.cloud.tencent.com/g' /etc/apt/sources.list
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update gcc musl-dev
WORKDIR /go/src/news
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go mod tidy
#RUN CGO_ENABLED=1 GOOS=linux go build
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .
RUN cp config.yaml.example config.yaml

FROM scratch

COPY --from=builder /go/src/news/news /
COPY --from=builder /go/src/news/config.yaml.example /config.yaml
COPY --from=builder /go/src/news/templates /templates
COPY --from=builder /go/src/news/data /data
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Shanghai
ENV ZONEINFO=/zoneinfo.zip

EXPOSE "9999"
ENTRYPOINT ["/news"]
