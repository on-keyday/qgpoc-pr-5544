#!/bin/bash
sysctl -w net.ipv4.ip_forward=1
# add ipip tunnel to server vip
ip tunnel add ipip0 mode ipip remote 10.220.0.2 local 10.220.0.254
ip link set ipip0 up
ip route add ${SERVER_VIP}/32 dev ipip0
# drop packets from server real ip to simulate drop by nat
nft add table inet quic_poc
nft add chain inet quic_poc forward \{ type filter hook forward priority 0 \; policy accept \; \}
nft add rule inet quic_poc forward ip saddr 10.220.0.2 drop
nft list ruleset
tail -f /dev/null
