
FROM alpine:latest as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add go

RUN mkdir -p /tmp/http-echo
COPY . /tmp/http-echo/

RUN cd /tmp/http-echo; go build

FROM alpine:latest

COPY --from=builder /tmp/http-echo/http-echo /

ENTRYPOINT ("/http-echo")

