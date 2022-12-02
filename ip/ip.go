package ip

import (
	"github.com/gofiber/fiber/v2"
	"net"
	"time"
	"verify/utils/ip/geoip"
)

//Country is a simple IP-country code lookup.
//Returns an empty string when cannot determine country.
func Country(ip string) string {
	return geoip.Country(ip)
}

//CountryByNetIP is a simple IP-country code lookup.
//Returns an empty string when cannot determine country.
func CountryByNetIP(ip net.IP) string {
	return geoip.CountryByIP(ip)
}

func Detect(c *fiber.Ctx) error {
	ip := FromRequest(c)
	c.Locals("ip", ip)
	c.Locals("ip_country", Country(ip))
	return c.Next()
}

func FromRequest(c *fiber.Ctx) string {
	return geoip.FromRequest(c)
}

func ChangeTimezone(dt time.Time, timezone string) time.Time {
	loc, _ := time.LoadLocation(timezone)
	newTime := dt.In(loc)
	return newTime
}
