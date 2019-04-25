# docker build -t shuxs/fproxy:latest . && docker push shuxs/fproxy:latest
FROM shuxs/alpine:latest

WORKDIR /app

ENV api="http://serviceName:servicePort/" \
    www="www"

ADD dist/fproxy-min /usr/bin/fproxy

CMD ["fproxy"]
