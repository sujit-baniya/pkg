package email

var freeDomain = map[string]bool{
}

func IsFreeDomain(domain string) bool {
	if _, ok := freeDomain[domain]; ok {
		return true
	}
	return false
}