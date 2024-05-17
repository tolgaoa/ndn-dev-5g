# Use the latest Ubuntu image for the base.
FROM ubuntu:focal

# Install the iptables command.
RUN apt-get update && \
    apt-get install -y iptables wget curl vim net-tools make

RUN wget https://dl.google.com/go/go1.22.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

COPY ./ /proxy/src/
WORKDIR /proxy/src/
RUN make build

RUN adduser proxyop -u 1005 --gecos "proxyop, , , " --disabled-password
#RUN chown proxyop:proxyop /tmp/streamsa/src
#RUN chmod -R u+w /streamsa/src

USER proxyop

ENV PATH="$PATH:/usr/local/go/bin"
WORKDIR /proxy/src/

#ENTRYPOINT [ "su", "proxyop", "/bin/bash", "-l", "-c" ]
CMD ["su", "proxyop", "-c", "./bin/proxy"]
#CMD ["sleep", "infinity"]
