FROM golang:1.19-alpine
RUN apk add --no-cache git make gcc musl-dev linux-headers
RUN git clone https://github.com/octaspace/go-octa.git
RUN cd go-octa && make geth && cp build/bin/geth /usr/local/bin
RUN mkdir /root/.ethereum
ENTRYPOINT [ "geth" ]