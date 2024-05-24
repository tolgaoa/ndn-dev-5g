#!/bin/bash

# Forward incoming HTTP traffic on port 80 to the HTTP/1.1 proxy port 11095
iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095

# Forward incoming HTTPS traffic on port 443 to the HTTP/2 proxy port 11096
iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 11096
iptables -t nat -A OUTPUT -p tcp --dport 443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096

# Forward incoming HTTP/3 traffic on port 8443 to the HTTP/3 proxy port 11096 (QUIC uses UDP)
iptables -t nat -A PREROUTING -p udp --dport 8443 -j REDIRECT --to-port 11096
iptables -t nat -A OUTPUT -p udp --dport 8443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096

# List all iptables rules to verify
iptables -t nat --list

