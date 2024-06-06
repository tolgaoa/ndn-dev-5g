FROM ubuntu:focal

RUN apt-get update && \
    apt-get install -y iptables wget curl vim net-tools make

RUN wget https://dl.google.com/go/go1.22.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && rm go1.22.3.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

COPY ./ /proxy/src/
WORKDIR /proxy/src/
RUN make build

RUN adduser --disabled-password --gecos '' --uid 1005 proxyop

USER proxyop

ENV PATH="$PATH:/usr/local/go/bin"

WORKDIR /proxy/src/

# Specify the command to run your application.
CMD ["./bin/proxy"]
#CMD ["sleep", "infinity"]

