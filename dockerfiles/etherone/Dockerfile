FROM golang:1.19-alpine
RUN apk add --no-cache git make gcc musl-dev linux-headers
RUN git clone https://github.com/etheronechain/go-ethereum.git
RUN git clone https://github.com/etheronechain/Genesis.git  
RUN cd go-ethereum && make geth && cp build/bin/geth /usr/local/bin
RUN mkdir /root/.ethereum
RUN cd Genesis && cp genesis.json /root/.ethereum/genesis.json
RUN geth --datadir /root/.ethereum init /root/.ethereum/genesis.json
