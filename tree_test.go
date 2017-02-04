package main

import (
	"net"
	"testing"
)

func TestTree(t *testing.T) {
	tree := New()

	for i, n := range []string{
		"192.168.0.0/12",
		"192.168.0.0/24",
		"192.175.0.0/30",
		"2001:db8:1234::/48",
	} {
		_, network, _ := net.ParseCIDR(n)
		tree.Insert(network, i)
	}

	for _, tc := range []struct {
		ip       string
		expected int
		nothing  bool
	}{
		{ip: "", nothing: true},
		{ip: "1.2.3.4", nothing: true},
		{ip: "192.168.0.1", expected: 1},       // most specific route is 192.168.0.0/24
		{ip: "192.175.0.32", expected: 0},      // most specific route is 192.168.0.0/12
		{ip: "192.175.0.1", expected: 2},       // most specific route is 192.175.0.0/30
		{ip: "2001:db8:1234::15", expected: 3}, // most specific route is 2001:db8:1234::/48
	} {
		ip := net.ParseIP(tc.ip)
		v := tree.Lookup(ip)
		if v == nil {
			if tc.nothing {
				continue
			}
			t.Fatalf("no valid network found: expected %d", tc.expected)
		}
		if i := v.(int); v != tc.expected {
			t.Fatalf("unexpected value returned: got %d, expected %d", i, tc.expected)
		}
	}
}
