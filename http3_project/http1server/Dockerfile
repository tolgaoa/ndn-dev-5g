FROM ubuntu:focal

# Install the iptables command.
RUN apt-get update && \
    apt-get install -y iptables wget curl vim net-tools make

RUN wget https://dl.google.com/go/go1.22.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

COPY ./ /app/src/
WORKDIR /app/src/
RUN make build

#ENTRYPOINT [ "su", "proxyop", "/bin/bash", "-l", "-c" ]
CMD ["./bin/http1server"]
#CMD ["sleep", "infinity"]
