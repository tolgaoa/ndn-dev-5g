#!/bin/bash

# Clear any existing rules
iptables -t nat -F

# Forward incoming HTTP traffic on port 80 to the HTTP/1.1 proxy port 11095
iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095

# Forward incoming HTTPS traffic on port 443 to the HTTPS proxy port 11096
iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 11096
iptables -t nat -A OUTPUT -p tcp --dport 443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096

# List all iptables rules to verify
iptables -t nat --list

