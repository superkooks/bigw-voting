package util

import "net"

// IsPublicIP returns whether the ip is routable globally
func IsPublicIP(in string) bool {
	inIP, err := net.ResolveIPAddr("ip4", in)
	if err != nil {
		return false
	}

	ip := inIP.IP

	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}

	return false
}
