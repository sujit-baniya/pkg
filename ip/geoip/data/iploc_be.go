//go:build arm || ppc64be || mipsbe || mips64be || mips64p32be
// +build arm ppc64be mipsbe mips64be mips64p32be

package data

import (
	_ "embed" // for ip data
)

//go:embed ipv4be.bin
var Ip4bin []byte

//go:embed ipv6be.gz
var Ip6bin []byte
