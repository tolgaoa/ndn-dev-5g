#!/bin/bash

clear_rules() {
    iptables -t nat -F
}

list_rules() {
    iptables -t nat --list
}

setup_http1() {
    echo "Setting up iptables for HTTP1 mode"
    iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
    iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095
}

setup_https() {
    echo "Setting up iptables for HTTPS mode"

	iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
	iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095
	iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 11096
	iptables -t nat -A OUTPUT -p tcp --dport 443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096

}

setup_http2() {
    echo "Setting up iptables for HTTP2 mode"
    iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
    iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095
    iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 11096
    iptables -t nat -A OUTPUT -p tcp --dport 443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096
}

setup_http3() {
    echo "Setting up iptables for HTTP3 mode"
    iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 11095
    iptables -t nat -A OUTPUT -p tcp --dport 80 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095
    iptables -t nat -A PREROUTING -p udp --dport 8443 -j REDIRECT --to-port 11096
    iptables -t nat -A OUTPUT -p udp --dport 8443 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11096
}

clear_qdisc() {
    tc qdisc del dev eth0 root || true
}

apply_packet_loss() {
    local packet_loss_rate=$1
	clear_qdisc
    echo "Applying packet loss: $packet_loss_rate"
    tc qdisc add dev eth0 root netem loss "$packet_loss_rate"
}

clear_rules

case "$OPERATION_MODE" in
    "HTTP1")
        setup_http1
        ;;
    "HTTPS")
        setup_https
        ;;
    "HTTP2")
        setup_http2
        ;;
    "HTTP3")
        setup_http3
        ;;
    *)
        echo "Invalid OPERATION_MODE: $OPERATION_MODE"
        exit 1
        ;;
esac

if [ -n "$PACKET_LOSS_RATE" ]; then
    apply_packet_loss "$PACKET_LOSS_RATE"
fi

list_rules
