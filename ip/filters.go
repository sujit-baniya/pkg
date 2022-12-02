package ip

import (
	"github.com/gofiber/fiber/v2"
	"verify/utils/ip/geoip"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

//Config for Filter. Allow supersedes Block for IP checks
//across all matching subnets, whereas country checks use the
//latest Allow/Block setting.
//IPs can be IPv4 or IPv6 and can optionally contain subnet
//masks (e.g. /24). Note however, determining if a given IP is
//included in a subnet requires a linear scan so is less performant
//than looking up single IPs.
//
//This could be improved with cidr range prefix tree.
type Config struct {
	ErrorHandler fiber.Handler
	//explicitly allowed IPs
	AllowedIPs []string
	//explicitly blocked IPs
	BlockedIPs []string
	// Store IP on context local
	IPContextKey string
	//explicitly allowed country ISO codes
	AllowedCountries []string
	//explicitly blocked country ISO codes
	BlockedCountries []string
	//block by default (defaults to allow)
	BlockByDefault bool
	// TrustProxy enable check request IP from proxy
	TrustProxy bool
	// Logger enables logging, printing using the provided interface
	Logger interface {
		Printf(format string, v ...interface{})
	}
	// These fields currently have no effect
	IPDB         []byte
	IPDBPath     string
	IPDBNoFetch  bool
	IPDBFetchURL string
}

type Filter struct {
	opts Config
	//mut protects the below
	//rw since writes are rare
	mut            sync.RWMutex
	defaultAllowed bool
	ips            map[string]bool
	codes          map[string]bool
	subnets        []*subnet
}

type subnet struct {
	str     string
	ipNet   *net.IPNet
	allowed bool
}

var filter = &Filter{}

//NewFilter constructs Filter instance without downloading DB.
func NewFilter(cfg ...Config) func(c *fiber.Ctx) error {
	var opts Config
	if len(cfg) > 0 {
		opts = cfg[0]
	}
	if opts.Logger == nil {
		//disable logging by default
		opts.Logger = log.New(ioutil.Discard, "", 0)
	}
	filter = &Filter{
		opts:           opts,
		ips:            map[string]bool{},
		codes:          map[string]bool{},
		defaultAllowed: !opts.BlockByDefault,
	}
	for _, ip := range opts.BlockedIPs {
		filter.BlockIP(ip)
	}
	for _, ip := range opts.AllowedIPs {
		filter.AllowIP(ip)
	}
	for _, code := range opts.BlockedCountries {
		filter.BlockCountry(code)
	}
	for _, code := range opts.AllowedCountries {
		filter.AllowCountry(code)
	}
	if opts.IPContextKey == "" {
		opts.IPContextKey = "ip"
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":   true,
				"message": fiber.StatusServiceUnavailable,
			})
		}
	}
	return func(c *fiber.Ctx) error {
		var remoteIP string
		rIP := c.Locals(opts.IPContextKey)
		if rIP != nil {
			remoteIP = rIP.(string)
		} else {
			remoteIP = geoip.FromRequest(c)
			c.Locals(opts.IPContextKey, remoteIP)
		}
		allowed := filter.Allowed(remoteIP)
		//special case localhost ipv4
		if !allowed && remoteIP == "::1" && filter.Allowed("127.0.0.1") {
			allowed = true
		}
		if !allowed {
			return opts.ErrorHandler(c)
		}
		//success!
		return c.Next()
	}
}

func (f *Filter) AllowIP(ip string) bool {
	return f.ToggleIP(ip, true)
}

func (f *Filter) BlockIP(ip string) bool {
	return f.ToggleIP(ip, false)
}

func (f *Filter) ToggleIP(str string, allowed bool) bool {
	//check if has subnet
	if ip, nt, err := net.ParseCIDR(str); err == nil {
		// containing only one ip? (no bits masked)
		if n, total := nt.Mask.Size(); n == total {
			f.mut.Lock()
			f.ips[ip.String()] = allowed
			f.mut.Unlock()
			return true
		}
		//check for existing
		f.mut.Lock()
		found := false
		for _, subnet := range f.subnets {
			if subnet.str == str {
				found = true
				subnet.allowed = allowed
				break
			}
		}
		if !found {
			f.subnets = append(f.subnets, &subnet{
				str:     str,
				ipNet:   nt,
				allowed: allowed,
			})
		}
		f.mut.Unlock()
		return true
	}
	//check if plain ip (/32)
	if ip := net.ParseIP(str); ip != nil {
		f.mut.Lock()
		f.ips[ip.String()] = allowed
		f.mut.Unlock()
		return true
	}
	return false
}

func (f *Filter) AllowCountry(code string) {
	f.ToggleCountry(code, true)
}

func (f *Filter) BlockCountry(code string) {
	f.ToggleCountry(code, false)
}

//ToggleCountry alters a specific country setting
func (f *Filter) ToggleCountry(code string, allowed bool) {

	f.mut.Lock()
	f.codes[code] = allowed
	f.mut.Unlock()
}

//ToggleDefault alters the default setting
func (f *Filter) ToggleDefault(allowed bool) {
	f.mut.Lock()
	f.defaultAllowed = allowed
	f.mut.Unlock()
}

//Allowed returns if a given IP can pass through the filter
func (f *Filter) Allowed(ipStr string) bool {
	return f.NetAllowed(net.ParseIP(ipStr))
}

//NetAllowed returns if a given net.IP can pass through the filter
func (f *Filter) NetAllowed(ip net.IP) bool {
	//invalid ip
	if ip == nil {
		return false
	}
	//read lock entire function
	//except for db access
	f.mut.RLock()
	defer f.mut.RUnlock()
	//check single ips
	allowed, ok := f.ips[ip.String()]
	if ok {
		return allowed
	}
	//scan subnets for any allow/block
	blocked := false
	for _, subnet := range f.subnets {
		if subnet.ipNet.Contains(ip) {
			if subnet.allowed {
				return true
			}
			blocked = true
		}
	}
	if blocked {
		return false
	}
	//check country codes
	code := geoip.CountryByIP(ip)
	if code != "" {
		if allowed, ok := f.codes[code]; ok {
			return allowed
		}
	}
	//use default setting
	return f.defaultAllowed
}

//Blocked returns if a given IP can NOT pass through the filter
func (f *Filter) Blocked(ip string) bool {
	return !f.Allowed(ip)
}

//NetBlocked returns if a given net.IP can NOT pass through the filter
func (f *Filter) NetBlocked(ip net.IP) bool {
	return !f.NetAllowed(ip)
}

func (f *Filter) IPToCountry(ip string) string {
	return geoip.Country(ip)
}

func (f *Filter) NetIPToCountry(ip net.IP) string {
	return geoip.CountryByIP(ip)
}

func AllowIP(ip string) bool {
	return filter.AllowIP(ip)
}

func BlockIP(ip string) bool {
	return filter.BlockIP(ip)
}

func ToggleIP(str string, allowed bool) bool {
	return filter.ToggleIP(str, allowed)
}

func AllowCountry(code string) {
	filter.AllowCountry(code)
}

func BlockCountry(code string) {
	filter.BlockCountry(code)
}

//ToggleCountry alters a specific country setting
func ToggleCountry(code string, allowed bool) {
	filter.ToggleCountry(code, allowed)
}

//ToggleDefault alters the default setting
func ToggleDefault(allowed bool) {
	filter.ToggleDefault(allowed)
}

//Allowed returns if a given IP can pass through the filter
func Allowed(ipStr string) bool {
	return filter.Allowed(ipStr)
}

//NetAllowed returns if a given net.IP can pass through the filter
func NetAllowed(ip net.IP) bool {
	return filter.NetAllowed(ip)
}

//Blocked returns if a given IP can NOT pass through the filter
func Blocked(ip string) bool {
	return filter.Blocked(ip)
}

//NetBlocked returns if a given net.IP can NOT pass through the filter
func NetBlocked(ip net.IP) bool {
	return filter.NetBlocked(ip)
}

func IPToCountry(ip string) string {
	return filter.IPToCountry(ip)
}

func NetIPToCountry(ip net.IP) string {
	return filter.NetIPToCountry(ip)
}
