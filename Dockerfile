FROM golang:1.21

WORKDIR /usr/src/money

# 设置代理，防止后面因为墙的原因下载go.mod依赖信息超时
RUN go env -w GOPROXY=https://goproxy.cn,direct

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/money .
RUN go build -v -o /usr/local/bin/workers ./workers/workers.go

CMD ["money"]