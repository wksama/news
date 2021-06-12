FROM golang:alpine as builder
WORKDIR /go/src/penti
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-w -s -extldflags -static"

FROM scratch

COPY --from=builder /go/src/penti/penti /
COPY --from=builder /go/src/penti/config/config.yaml.example /config/config.yaml
COPY --from=builder /go/src/penti/templates /templates

EXPOSE "9999"
ENTRYPOINT ["/penti"]
