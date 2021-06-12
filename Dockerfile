FROM golang:alpine as builder
WORKDIR /go/src/penti
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-w -s -extldflags -static"
EXPOSE "9999"
VOLUME /config
ENTRYPOINT ["./penti"]
