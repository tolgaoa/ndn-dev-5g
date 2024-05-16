#!/bin/bash

iptables -t nat -A PREROUTING -p tcp --dport 8080 -j REDIRECT --to-port 11095
iptables -t nat -A OUTPUT -p tcp --dport 8080 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095

iptables -t nat -A PREROUTING -p udp --dport 443 -j REDIRECT --to-port 11096
iptables -t nat -A OUTPUT -p tcp --dport 443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096

# List all iptables rules.
iptables -t nat --list
