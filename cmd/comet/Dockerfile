FROM golang:1.19-buster as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
WORKDIR /go/work
COPY . .
VOLUME /go/work/log
RUN go build main.go
EXPOSE 8080
EXPOSE 8081
EXPOSE 8083
ENTRYPOINT ["./main"]