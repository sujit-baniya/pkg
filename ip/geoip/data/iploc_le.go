//go:build 386 || amd64 || amd64p32 || arm64 || ppc64le || mipsle || mips64le || mips64p32le
// +build 386 amd64 amd64p32 arm64 ppc64le mipsle mips64le mips64p32le

package data

import (
	_ "embed" // for ip data
)

//go:embed ipv4le.bin
var Ip4bin []byte

//go:embed ipv6le.gz
var Ip6bin []byte
