// package geoip provides fastest Geolocation Country library for Go.
//
//	package main
//
//	import (
//		"fmt"
//		"net"
//		"github.com/phuslu/iploc"
//	)
//
//	func main() {
//		fmt.Printf("%s", iploc.Country(net.ParseIP("2001:4860:4860::8888")))
//	}
//
//	// Output: US
package geoip

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/sujit-baniya/pkg/ip/geoip/data"
	"io"
	"net"
	"os"
	"reflect"
	"unsafe"
)

// Version is iplocation database version.
const Version = "v1.0.20211029"

var ip4uint []uint32
var ip6uint []uint64

func init() {
	// ipv4
	ip4uint = *(*[]uint32)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&data.Ip4bin[0])),
		Len:  len(data.Ip4bin) / 4,
		Cap:  len(data.Ip4bin) / 4,
	}))

	// ipv6
	if os.Getenv("IPLOC_IPV4ONLY") == "" {
		r, _ := gzip.NewReader(bytes.NewReader(data.Ip6bin))
		data.Ip6bin, _ = io.ReadAll(r)
		ip6uint = *(*[]uint64)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(&data.Ip6bin[0])),
			Len:  len(data.Ip6bin) / 8,
			Cap:  len(data.Ip6bin) / 8,
		}))
	} else {
		ip6uint = []uint64{0, 0}
	}
}

// Country return ISO 3166-1 alpha-2 country code of IP.
func countryByIP(ip net.IP) []byte {
	if ip == nil {
		return nil
	}

	if ip4 := ip.To4(); ip4 != nil {
		// ipv4
		n := binary.BigEndian.Uint32(ip4)
		i, j := 0, len(ip4uint)
		_ = ip4uint[j-1]
		for i < j {
			h := (i + j) >> 1
			if ip4uint[h] > n {
				j = h
			} else {
				i = h + 1
			}
		}
		return data.Ip4txt[i*2-2 : i*2]
	}
	// ipv6
	high := binary.BigEndian.Uint64(ip)
	low := binary.BigEndian.Uint64(ip[8:])
	i, j := 0, len(ip6uint)
	_ = ip6uint[j-1]
	for i < j {
		h := (i + j) >> 1 & ^1
		n := ip6uint[h]
		if n > high || (n == high && ip6uint[h+1] > low) {
			j = h
		} else {
			i = h + 2
		}
	}
	return data.Ip6txt[i-2 : i]
}

func Country(ip string) string {
	return string(countryByIP(net.ParseIP(ip)))
}

func CountryByIP(ip net.IP) string {
	return string(countryByIP(ip))
}

// IsReservedIPv4 detects a net.IP is a reserved address, return false if IPv6
func IsReservedIPv4(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	_ = ip4[3]
	switch ip4[0] {
	case 10:
		return true
	case 100:
		return ip4[1] >= 64 && ip4[1] <= 127
	case 127:
		return true
	case 169:
		return ip4[1] == 254
	case 172:
		return ip4[1] >= 16 && ip4[1] <= 31
	case 192:
		switch ip4[1] {
		case 0:
			switch ip4[2] {
			case 0, 2:
				return true
			}
		case 18, 19:
			return true
		case 51:
			return ip4[2] == 100
		case 88:
			return ip4[2] == 99
		case 168:
			return true
		}
	case 203:
		return ip4[1] == 0 && ip4[2] == 113
	case 224:
		return true
	case 240:
		return true
	}
	return false
}
