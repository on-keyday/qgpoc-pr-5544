# Poc of quic-go connection migration problem

see also https://github.com/quic-go/quic-go/pull/5544

## How to setup

1. Run
```
# enable ipip if not enabled
$ modprobe ipip
$ docker compose up -d
```

2. After launched, replace client ip address
```
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


After patch applied, behavior is like below.
```

```



