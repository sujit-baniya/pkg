package ip

import (
	"context"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/pkg/ip/geoip"
	"net"
	"time"
)

// Country is a simple IP-country code lookup.
// Returns an empty string when cannot determine country.
func Country(ip string) string {
	return geoip.Country(ip)
}

// CountryByNetIP is a simple IP-country code lookup.
// Returns an empty string when cannot determine country.
func CountryByNetIP(ip net.IP) string {
	return geoip.CountryByIP(ip)
}

func Detect(ctx context.Context, c *frame.Context) {
	ip := FromRequest(c)
	c.Set("ip", ip)
	c.Set("ip_country", Country(ip))
	c.Next(ctx)
}

func FromRequest(c *frame.Context) string {
	return geoip.FromRequest(c)
}

func ChangeTimezone(dt time.Time, timezone string) time.Time {
	loc, _ := time.LoadLocation(timezone)
	newTime := dt.In(loc)
	return newTime
}
