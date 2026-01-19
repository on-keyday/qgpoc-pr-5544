# Poc of quic-go connection migration problem

see also https://github.com/quic-go/quic-go/pull/5544

## How to setup

1. Run
```
# enable ipip tunnel if not enabled
$ modprobe ipip
$ docker compose up -d
```

2. After launched, replace client ip address
```
# network interface name depends on runtime...
$ docker compose exec client ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host proto kernel_lo
       valid_lft forever preferred_lft forever
2: eth0@if2783: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 6e:fd:c4:d9:80:e7 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.210.0.2/24 brd 10.210.0.255 scope global eth0
       valid_lft forever preferred_lft forever
3: eth1@if2784: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether c2:b5:16:23:02:21 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.200.0.2/24 brd 10.200.0.255 scope global eth1
       valid_lft forever preferred_lft forever
$ docker compose exec client ip route replace 10.250.1.1/32 via 10.210.0.254 dev eth0
```

3. Can see endless Path Probing

Server VIP: 10.250.0.2
Server Real IP: 10.220.0.2
Client on Path1: 10.200.0.2
Client on Path2: 10.210.0.2

```
$ docker compose exec router1 tcpdump -i any -n
tcpdump: WARNING: any: That device doesn't support promiscuous mode
(Promiscuous mode not supported on the "any" device)
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on any, link-type LINUX_SLL2 (Linux cooked v2), snapshot length 262144 bytes
13:43:27.596782 ipip0 Out IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 56
13:43:27.596790 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 56
13:43:27.596973 eth2  In  IP 10.220.0.2.8889 > 10.210.0.2.39701: UDP, length 1200
13:43:27.622372 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:27.622392 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:28.571091 eth1  In  IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:28.571124 ipip0 Out IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:28.571139 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:28.597393 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 35
13:43:28.597439 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 35
13:43:28.597922 eth1  In  IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 57
13:43:28.597944 ipip0 Out IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 57
13:43:28.597961 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 57
13:43:28.598461 eth2  In  IP 10.220.0.2.8889 > 10.210.0.2.39701: UDP, length 1200
13:43:28.623944 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:28.623975 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:29.571405 eth1  In  IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:29.571434 ipip0 Out IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:29.571447 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 32
13:43:29.598081 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:29.598099 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 33
13:43:29.599255 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 25
13:43:29.599266 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.39701: UDP, length 25
13:43:29.599437 eth1  In  IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 57
13:43:29.599451 ipip0 Out IP 10.210.0.2.39701 > 10.250.1.1.8889: UDP, length 57
```

# How to apply patch
1. edit go.mod 
```diff
-// replace github.com/quic-go/quic-go => github.com/on-keyday/quic-go v0.0.0-20260118200636-5a7fa253b928
+   replace github.com/quic-go/quic-go => github.com/on-keyday/quic-go v0.0.0-20260118200636-5a7fa253b928
```
2.  run `go mod tidy && docker compose down && docker compose up -d --build`
 - its may be ok to run only compose up, but sometimes cache causes unexpected behavior...

After patch applied, behavior is like below.
```
(omitted)
14:00:09.493171 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.200.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:09.520327 eth2  In  IP 10.250.1.1.8889 > 10.200.0.2.48934: UDP, length 32
14:00:09.520368 eth0  Out IP 10.250.1.1.8889 > 10.200.0.2.48934: UDP, length 32
14:00:09.982743 eth2  Out ARP, Request who-has 10.220.0.2 tell 10.220.0.254, length 28
14:00:09.982788 eth2  In  ARP, Reply 10.220.0.2 is-at e2:89:75:15:0b:8c, length 28
14:00:10.493084 eth1  B   ARP, Request who-has 10.210.0.254 tell 10.210.0.2, length 28
14:00:10.493121 eth1  Out ARP, Reply 10.210.0.254 is-at de:af:49:69:37:dd, length 28
14:00:10.493134 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:10.493160 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:10.493175 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:10.493731 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1200
14:00:10.493765 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1200
14:00:10.495230 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 40
14:00:10.495259 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 40
14:00:10.495271 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 40
14:00:10.496150 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 32
14:00:10.496185 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 32
14:00:11.492926 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:11.492970 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:11.492988 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:11.494059 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1366
14:00:11.494087 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1366
14:00:11.521672 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 32
14:00:11.521720 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 32
14:00:11.527953 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:11.527994 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:11.528010 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.493043 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.493078 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.493094 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.493825 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1409
14:00:12.493865 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1409
14:00:12.519109 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:12.519150 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:12.522206 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.522241 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:12.522257 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.493041 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.493054 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.493062 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.494416 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1430
14:00:13.494429 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1430
14:00:13.519957 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:13.519976 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:13.520458 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.520473 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.520481 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:13.566712 eth0  Out ARP, Request who-has 10.200.0.2 tell 10.200.0.254, length 28
14:00:13.566775 eth0  In  ARP, Reply 10.200.0.2 is-at 42:ff:15:65:19:1a, length 28
14:00:14.493842 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:14.493882 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:14.493898 eth2  Out IP 10.220.0.254 > 10.220.0.2: IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 32
14:00:14.494604 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1441
14:00:14.494639 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 1441
14:00:14.521017 eth2  In  IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:14.521062 eth1  Out IP 10.250.1.1.8889 > 10.210.0.2.48934: UDP, length 33
14:00:14.522026 eth1  In  IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 33
14:00:14.522068 ipip0 Out IP 10.210.0.2.48934 > 10.250.1.1.8889: UDP, length 33
```



