package email

var disposableDomains = map[string]bool{
}

func isDisposable(domain string) bool {
	if _, ok := disposableDomains[domain]; ok {
		return true
	}
	return false
}