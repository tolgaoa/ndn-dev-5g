FROM ubuntu:focal

RUN apt-get update && \
    apt-get install -y iptables wget curl vim net-tools make

RUN wget https://dl.google.com/go/go1.22.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && rm go1.22.3.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

COPY ./ /client/src/
WORKDIR /client/src/
RUN make build

# Specify the command to run your application.
#CMD ["./bin/proxy"]
CMD ["sleep", "infinity"]

