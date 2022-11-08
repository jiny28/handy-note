FROM golang:1.18.4-buster as builder
ENV TDENGINE_VERSION=2.6.0.6
ENV HOSTARCH=Linux-aarch64
RUN wget -c https://www.taosdata.com/assets-download/TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz \
   && tar xvf TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz \
   && cd TDengine-client-${TDENGINE_VERSION} \
   && ./install_client.sh \
   && cd ../ \
   && rm -rf TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz TDengine-client-${TDENGINE_VERSION}
RUN mkdir /app
COPY . /app
WORKDIR /app
ENV GOPROXY="https://goproxy.io,direct"
RUN go env
RUN go mod tidy
RUN go build -o app

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y wget
ENV TDENGINE_VERSION=2.6.0.6
ENV HOSTARCH=Linux-aarch64
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive
RUN wget -c https://www.taosdata.com/assets-download/TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz \
   && tar xvf TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz \
   && cd TDengine-client-${TDENGINE_VERSION} \
   && ./install_client.sh \
   && cd ../ \
   && rm -rf TDengine-client-${TDENGINE_VERSION}-${HOSTARCH}.tar.gz TDengine-client-${TDENGINE_VERSION}

RUN apt install -y tzdata \
    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/app /usr/bin/
CMD ["app"]