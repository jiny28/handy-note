FROM golang:1.17.6-buster as builder
ENV TDENGINE_VERSION=2.4.0.0
RUN wget -c https://www.taosdata.com/assets-download/TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && tar xvf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && cd TDengine-client-${TDENGINE_VERSION} \
   && ./install_client.sh \
   && cd ../ \
   && rm -rf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz TDengine-client-${TDENGINE_VERSION}
WORKDIR /usr/src/app/
ENV GOPROXY="https://goproxy.io,direct"
COPY ./main.go ./go.mod ./go.sum /usr/src/app/
RUN go env
RUN go mod tidy
RUN go build

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y wget
ENV TDENGINE_VERSION=2.4.0.0
RUN wget -c https://www.taosdata.com/assets-download/TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && tar xvf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && cd TDengine-client-${TDENGINE_VERSION} \
   && ./install_client.sh \
   && cd ../ \
   && rm -rf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz TDengine-client-${TDENGINE_VERSION}

## add your application next, eg. go, build it in builder stage, copy the binary to the runtime
COPY --from=builder /usr/src/app/app /usr/bin/
CMD ["app"]