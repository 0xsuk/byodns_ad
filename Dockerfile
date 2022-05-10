#For raspberry pi 
FROM ubuntu:16.04
MAINTAINER 0xsuk
RUN apt-get update
RUN apt-get install -y git redis-server tar wget build-essential

# install go
RUN mkdir -p /root/go/src/github.com/0xsuk
WORKDIR /root/go/src/github.com/0xsuk
RUN git clone https://github.com/0xsuk/byodns.git
WORKDIR /root/go/src/github.com/0xsuk/byodns
RUN wget https://golang.org/dl/go1.17.1.linux-armv6l.tar.gz -O go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> /root/.bashrc
# we dno't sourced .bashrc
ENV PATH="$PATH:/usr/local/go/bin"
ENV GOPATH="/root/go"
ENV PATH="$PATH:/root/go/bin"

# setup
RUN mkdir -p /etc/byodns
RUN cp default/domains.db /etc/byodns
RUN cp default/config.ini /etc/byodns

RUN echo "ALMOST THERE!!!"


RUN go build
RUN go install

EXPOSE 53 53/udp
EXPOSE 80 80/tcp

ENTRYPOINT redis-server & byodns
