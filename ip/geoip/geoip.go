package geoip

import (
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var cidrs []*net.IPNet

func init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}
}

var fetchIPFromString = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
var possibleHeaders = []string{
	"X-Original-Forwarded-For",
	"X-Forwarded-For",
	"X-Real-Ip",
	"X-Client-Ip",
	"Forwarded-For",
	"Forwarded",
	"Remote-Addr",
	"Client-Ip",
	"CF-Connecting-IP",
}

func isPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}
	if ipAddress.IsLoopback() || ipAddress.IsLinkLocalUnicast() || ipAddress.IsLinkLocalMulticast() {
		return true, nil
	}

	for i := range cidrs {
		if cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
}

// FromRequest determine user ip
func FromRequest(c *fiber.Ctx) string {
	var headerValue []byte
	if c.App().Config().ProxyHeader != "" && c.App().Config().ProxyHeader != "*" {
		headerValue = []byte(c.IP())
		if len(headerValue) <= 3 {
			headerValue = []byte("0.0.0.0")
		}
		return string(fetchIPFromString.Find(headerValue))
	}
	if c.App().Config().ProxyHeader == "*" {
		for _, headerName := range possibleHeaders {
			headerValue = c.Request().Header.Peek(headerName)
			if len(headerValue) > 3 {
				// Check list of IP in X-Forwarded-For and return the first global address
				for _, address := range strings.Split(string(headerValue), ",") {
					address = strings.TrimSpace(address)
					isPrivate, err := isPrivateAddress(address)
					if !isPrivate && err == nil {
						return string(fetchIPFromString.Find([]byte(address)))
					}
				}
				return string(fetchIPFromString.Find(headerValue))
			}
		}
	}
	headerValue = []byte(c.Context().RemoteIP().String())
	if len(headerValue) <= 3 {
		headerValue = []byte("0.0.0.0")
	}
	return string(fetchIPFromString.Find(headerValue))
}
