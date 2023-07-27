
FROM alpine:latest as builder

apk add go

RUN mkdir -p /tmp/http-echo
COPY . /tmp/http-echo/

RUN go build /tmp/http-echo

FROM alpine:latest

COPY --from=builder /tmp/http-echo/http-echo /

ENTRYPOINT ("/http-echo")

