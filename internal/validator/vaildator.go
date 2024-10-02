package validator

import "net/netip"

func IsIPInSubnet(ipaddr string, subnetCidr string) (bool, error) {
	network, err := netip.ParsePrefix(subnetCidr)
	if err != nil {
		return false, err
	}

	ip, err := netip.ParseAddr(ipaddr)
	if err != nil {
		return false, err
	}

	return network.Contains(ip), nil
}
