package xnet

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

func IpBetween(from, to, test net.IP) (bool, error) {
	if from == nil || to == nil || test == nil {
		return false, fmt.Errorf("ip input is nil")
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := test.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		return false, fmt.Errorf("ip did not convert to a 16 byte")
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true, nil
	}
	return false, nil
}

func IpBetweenStr(from, to, test string) (bool, error) {
	return IpBetween(net.ParseIP(from), net.ParseIP(to), net.ParseIP(test))
}

// IsInterIp
// 10.0.0.0/8：10.0.0.0～10.255.255.255
// 172.16.0.0/12：172.16.0.0～172.31.255.255
// 192.168.0.0/16：192.168.0.0～192.168.255.255
func IsInterIp(ip string) (bool, error) {
	ok, err := IpBetweenStr("10.0.0.0", "10.255.255.255", ip)
	if err != nil {
		return false, err
	}

	if !ok {
		ok, err = IpBetweenStr("172.16.0.0", "172.31.255.255", ip)
		if err != nil {
			return false, err
		}

		if !ok {
			ok, err = IpBetweenStr("192.168.0.0", "192.168.255.255", ip)
			if err != nil {
				return false, err
			}
		}
	}

	return ok, nil
}

func GetInterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println(ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}

	/*
		for _, addr := range addrs {
			//fmt.Printf("Inter %v\n", addr)
			ip := addr.String()
			if "10." == ip[:3] {
				return strings.Split(ip, "/")[0], nil
			} else if "172." == ip[:4] {
				return strings.Split(ip, "/")[0], nil
			} else if "196." == ip[:4] {
				return strings.Split(ip, "/")[0], nil
			} else if "192." == ip[:4] {
				return strings.Split(ip, "/")[0], nil
			}

		}
	*/

	return "", errors.New("no inter ip")
}

// GetExterIp 获取首个外网ip v4
func GetExterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// fmt.Printf("Inter %v\n", addr)
		ips := addr.String()
		idx := strings.LastIndex(ips, "/")
		if idx == -1 {
			continue
		}
		ipv := net.ParseIP(ips[:idx])
		if ipv == nil {
			continue
		}

		ipv4 := ipv.To4()
		if ipv4 == nil {
			// ipv6
			continue
		}
		ip := ipv4.String()

		// if "10." != ip[:3] && "172." != ip[:4] && "196." != ip[:4] && "127." != ip[:4] {
		//	return ip, nil
		// }
		ok, _ := IsInterIp(ip)
		if !ok && !ipv.IsLoopback() {
			return ip, nil
		}

	}

	return "", errors.New("no exter ip")
}

// GetListenAddr
// 不指定host使用内网host
// 指定了就使用指定的，不管指定的是0.0.0.0还是内网或者外网
func GetListenAddr(a string) (string, error) {
	addrTcp, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		return "", err
	}

	addr := addrTcp.String()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	if len(host) == 0 {
		return GetServAddr(addrTcp)
	}

	return addr, nil
}

func GetServAddr(a net.Addr) (string, error) {
	addr := a.String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		host = "0.0.0.0"
	}

	ip := net.ParseIP(host)

	if ip == nil {
		return "", fmt.Errorf("ParseIP error, host: %s", host)
	}
	/*
		fmt.Println("ADDR TYPE", ip,
			"IsGlobalUnicast",
			ip.IsGlobalUnicast(),
			"IsInterfaceLocalMulticast",
			ip.IsInterfaceLocalMulticast(),
			"IsLinkLocalMulticast",
			ip.IsLinkLocalMulticast(),
			"IsLinkLocalUnicast",
			ip.IsLinkLocalUnicast(),
			"IsLoopback",
			ip.IsLoopback(),
			"IsMulticast",
			ip.IsMulticast(),
			"IsUnspecified",
			ip.IsUnspecified(),
		)
	*/

	raddr := addr
	if ip.IsUnspecified() {
		// 没有指定ip的情况下，使用内网地址
		inerip, err := GetInterIp()
		if err != nil {
			return "", err
		}

		raddr = net.JoinHostPort(inerip, port)
	}

	// slog.Tracef("ServAddr --> addr:[%s] ip:[%s] host:[%s] port:[%s] raddr[%s]", addr, ip, host, port, raddr)

	return raddr, nil
}

// IpAddrFromRemoteAddr
// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func IpAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func IpAddrPort(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return ""
	}
	return s[idx+1:]
}
