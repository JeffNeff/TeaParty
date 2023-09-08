FROM golang:1.16-alpine

ARG GIT_SOURCE
ARG GIT_CHECKOUT_REF

RUN apk add --no-cache make gcc musl-dev linux-headers git

# Checkout latest version on Feb 8th 2021
RUN cd / && git clone https://github.com/bnb-chain/bsc.git\
    && cd ./bsc &&  make geth

RUN mv /bsc/build/bin/geth /usr/local/bin/

EXPOSE 30311 30311/udp