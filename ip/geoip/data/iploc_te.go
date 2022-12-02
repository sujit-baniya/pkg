package data

import _ "embed"

//go:embed ipv4.txt
var Ip4txt []byte

//go:embed ipv6.txt
var Ip6txt []byte
