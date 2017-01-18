package strings

import (
	"testing"
)

func TestIsIPv4(t *testing.T) {
	ip := "127.0.0.1"
	if !IsIPv4(ip) {
		t.Fatal(ip, "not true")
	}
	ip = "0.0.0.0"
	if !IsIPv4(ip) {
		t.Fatal(ip, "not true")
	}
	ip = "255.255.255.255"
	if !IsIPv4(ip) {
		t.Fatal(ip, "not true")
	}

	ip = "127.0.0.256"
	if IsIPv4(ip) {
		t.Fatal(ip, "true")
	}
	ip = "127.a.0.1"
	if IsIPv4(ip) {
		t.Fatal(ip, "true")
	}
	ip = "127.00.0.1"
	if IsIPv4(ip) {
		t.Fatal(ip, "true")
	}
	ip = "127.0.0.0.1"
	if IsIPv4(ip) {
		t.Fatal(ip, "true")
	}

	ip = ""
	if IsIPv4(ip) {
		t.Fatal(ip, "true")
	}
}
