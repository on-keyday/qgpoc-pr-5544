package util

import (
	"log"
	"net"
	"os/exec"
)

func AssignVIPToLoopback(vip string) {
	out, err := exec.Command("ip", "addr", "add", vip+"/32", "dev", "lo").CombinedOutput()
	if err != nil {
		log.Fatal("Failed to assign VIP to loopback:", string(out), err)
	}
}

func AddIPIPTunnel(remote, local string) {
	out, err := exec.Command("ip", "tunnel", "add", "ipip0", "mode", "ipip", "remote", remote, "local", local).CombinedOutput()
	if err != nil {
		log.Fatal("Failed to add IPIP tunnel:", string(out), err)
	}
	out, err = exec.Command("ip", "link", "set", "ipip0", "up").CombinedOutput()
	if err != nil {
		log.Fatal("Failed to set IPIP tunnel up:", string(out), err)
	}
}

func AddIPRoute(dest, via string) {
	out, err := exec.Command("ip", "route", "add", dest, "via", via).CombinedOutput()
	if err != nil {
		log.Fatal("Failed to add IP route:", string(out), err)
	}
}

func DisableRPFilters() {
	out, err := exec.Command("sysctl", "-w", "net.ipv4.conf.all.rp_filter=0").Output()
	if err != nil {
		log.Fatal("Failed to disable rp_filter:", string(out), err)
	}
	out, err = exec.Command("sysctl", "-w", "net.ipv4.conf.default.rp_filter=0").Output()
	if err != nil {
		log.Fatal("Failed to disable rp_filter:", string(out), err)
	}
	// listening interface names by net.Interfaces()
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Failed to get network interfaces:", err)
	}
	for _, iface := range ifaces {
		out, err = exec.Command("sysctl", "-w", "net.ipv4.conf."+iface.Name+".rp_filter=0").Output()
		if err != nil {
			log.Fatalf("Failed to disable rp_filter for %s: %s %v", iface.Name, string(out), err)
		}
	}
}
