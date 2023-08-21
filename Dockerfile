FROM golang:alpine as builder

WORKDIR /build

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

COPY . .

RUN apk update --no-cache \
    && apk upgrade \
    && apk add --no-cache bash \
            bash-doc \
            bash-completion \
    && apk add --no-cache tzdata \
    && rm -rf /var/cache/apk/* \
    && go mod download \
    && bash ./scripts/build-all.sh

FROM docker.io/epicmo/gugotik-basic:1.0 as prod

ENV TZ Asia/Shanghai

WORKDIR /work

RUN apk update --no-cache \
    && apk upgrade

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /build/output .