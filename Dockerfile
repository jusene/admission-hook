FROM alpine

WORKDIR /data/project

ENV LANG=zh_CN.UTF-8

COPY admission-webhook ./admission-webhook

RUN apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

CMD ["./admission-webhook"]