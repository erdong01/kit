package util

import "net"

func Ip2long(ip string) uint64 {
	if ip == "" {
		return 0
	}
	ipByte := net.ParseIP(string(ip))
	a := uint64(ipByte[12])
	b := uint64(ipByte[13])
	c := uint64(ipByte[14])
	d := uint64(ipByte[15])
	return uint64(a<<24 | b<<16 | c<<8 | d)
}

func Long2ip(ip uint64) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}
