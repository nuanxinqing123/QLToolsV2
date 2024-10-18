FROM  alpine:3.18.2 AS builder

ARG TARGETARCH

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /usr/src/QLToolsV2

# 安装项目必要环境
RUN \
  apk add --no-cache --update go go-bindata g++ ca-certificates tzdata

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .

# 打包项目文件
RUN \
  go build -ldflags '-linkmode external -s -w -extldflags "-static"' -o QLToolsV2-linux-$TARGETARCH
  

# FROM alpine:3.15
FROM ubuntu:22.10

MAINTAINER QLToolsV2 "nuanxinqing@gmail.com"

ARG TARGETARCH
ENV TARGET_ARCH=$TARGETARCH

WORKDIR /QLToolsV2

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/QLToolsV2/QLToolsV2-linux-$TARGETARCH /usr/src/QLToolsV2/docker-entrypoint.sh /usr/src/QLToolsV2/config/example.config.yaml ./

EXPOSE 1500

ENTRYPOINT ["sh", "docker-entrypoint.sh"]
