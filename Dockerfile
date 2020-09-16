FROM alpine:3.7
COPY ./linux_clickhouse_sinker /root/linux_clickhouse_sinker
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata
COPY ./start.sh /root/start.sh
ENTRYPOINT sh /root/start.sh
