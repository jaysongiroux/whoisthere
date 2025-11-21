package utils

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func GetAllTLDs(domain string, onlyCountry bool, onlyPopular bool) []string {
	// returns a lit of the domain with all possible tlds
	domains := []string{}

	if onlyCountry {
		for _, tld := range COUNTRY_TLDS {
			domains = append(domains, domain+"."+tld)
		}
	}

	if onlyPopular {
		for _, tld := range POPULAR_TLDS {
			domains = append(domains, domain+"."+tld)
		}
	}

	if !onlyCountry && !onlyPopular {
		for _, tld := range POPULAR_TLDS {
			domains = append(domains, domain+"."+tld)
		}
		for _, tld := range COUNTRY_TLDS {
			domains = append(domains, domain+"."+tld)
		}
	}
	return domains
}

func HasTLD(domain string) bool {
	domain = strings.TrimPrefix(strings.ToLower(domain), ".")
	if domain == "" || !strings.Contains(domain, ".") {
		return false
	}
	tld, icann := publicsuffix.PublicSuffix(domain)
	// icann=true means it's on the ICANN-managed list
	// if icann=false but tld != domain, it's a private suffix (still valid)
	return icann || tld != domain
}

func GetTLD(domain string) string {
	domain = strings.TrimPrefix(strings.ToLower(domain), ".")
	tld, _ := publicsuffix.PublicSuffix(domain)
	return tld
}

func IsPopularTLD(domain string) bool {
	tld := GetTLD(domain)
	for _, t := range POPULAR_TLDS {
		if tld == t {
			return true
		}
	}
	return false
}

func CleanAndValidateDomain(input string) (string, error) {
	input = strings.TrimSpace(input)

	// Reject emails
	if strings.Contains(input, "@") {
		return "", errors.New("invalid URL: looks like an email")
	}

	// Try parsing as a URL first
	if strings.Contains(input, "://") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.New("invalid URL")
		}
		input = u.Host
	}

	// If URL lacked scheme, attempt parsing with dummy scheme
	if !strings.Contains(input, "://") &&
		(strings.Contains(input, "/") || strings.Contains(input, ":")) {
		if u, err := url.Parse("https://" + input); err == nil && u.Host != "" {
			input = u.Host
		}
	}

	// Strip port if present
	host, _, err := net.SplitHostPort(input)
	if err == nil {
		input = host
	}

	// At this point input should be a bare domain
	input = strings.ToLower(strings.TrimSpace(input))

	if !domainRegex.MatchString(input) {
		return "", errors.New("invalid domain format")
	}

	// Validate label structure
	labels := strings.Split(input, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return "", errors.New("invalid domain label length")
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return "", errors.New("invalid domain: labels cannot start/end with hyphens")
		}
	}

	// last label must be a valid TLD (>=2 chars)
	if len(labels[len(labels)-1]) < 2 {
		return "", errors.New("invalid TLD")
	}

	return input, nil
}
