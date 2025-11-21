package utils

import (
	"strings"
	"testing"
)

func TestCleanAndValidateDomain(t *testing.T) {
	goodDomains := []map[string]interface{}{{
		"input":  "google.com",
		"output": "google.com",
	}, {
		"input":  "google.com/",
		"output": "google.com",
	}, {
		"input":  "bad.com////",
		"output": "bad.com",
	}, {
		"input":  "http://google.com/",
		"output": "google.com",
	}}

	badDomains := []string{"bad", "^&*("}

	for _, domain := range goodDomains {
		t.Run(domain["input"].(string), func(t *testing.T) {
			cleaned, err := CleanAndValidateDomain(domain["input"].(string))
			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}

			if cleaned != domain["output"].(string) {
				t.Errorf("Expected cleaned domain %s, but got %s", domain["input"].(string), cleaned)
			}
		})
	}

	for _, domain := range badDomains {
		t.Run(domain, func(t *testing.T) {
			cleaned, err := CleanAndValidateDomain(domain)
			if err == nil {
				t.Errorf("Expected an error, but got none for domain %s", domain)
			}

			if cleaned != "" {
				t.Errorf("Expected cleaned domain to be %s, but got %s", domain, cleaned)
			}
		})
	}
}

func TestIsPopularTLD(t *testing.T) {
	popularTLDs := []string{"a.com", "a.org", "a.net"}
	unpopularTLDs := []string{"a.com.br", "a.org.br", "a.net.br"}

	for _, tld := range popularTLDs {
		t.Run(tld, func(t *testing.T) {
			if !IsPopularTLD(tld) {
				t.Errorf("Expected %s to be popular, but got false", tld)
			}
		})
	}

	for _, tld := range unpopularTLDs {
		t.Run(tld, func(t *testing.T) {
			if IsPopularTLD(tld) {
				t.Errorf("Expected %s to be unpopular, but got true", tld)
			}
		})
	}
}

func TestGetTLD(t *testing.T) {
	domains := []map[string]interface{}{
		{
			"domain":       "google.com",
			"tld":          "com",
			"should_be_ok": true,
		}, {
			"domain":       "google.com.br",
			"tld":          "com.br",
			"should_be_ok": true,
		},
	}

	for _, domain := range domains {
		t.Run(domain["domain"].(string), func(t *testing.T) {
			tld := GetTLD(domain["domain"].(string))
			if tld != domain["tld"].(string) {
				t.Errorf("Expected TLD %s, but got %s", domain["tld"].(string), tld)
			}
		})
	}
}

func TestHasTLD(t *testing.T) {
	domains := []map[string]interface{}{
		{
			"domain":       "google.com",
			"should_be_ok": true,
		}, {
			"domain":       "google.com.br",
			"should_be_ok": true,
		}, {
			"domain":       "google",
			"should_be_ok": false,
		}, {
			"domain":       "google.com.br.com",
			"should_be_ok": true,
		}, {
			"domain":       "google.com.br.com.br",
			"should_be_ok": true,
		},
	}

	for _, domain := range domains {
		t.Run(domain["domain"].(string), func(t *testing.T) {
			hasTld := HasTLD(domain["domain"].(string))
			if hasTld != domain["should_be_ok"].(bool) {
				t.Errorf("Expected HasTLD(%s) to be %v, but got %v", domain["domain"].(string), domain["should_be_ok"].(bool), hasTld)
			}
		})
	}
}

func TestGetAllTLDs(t *testing.T) {
	domains := []map[string]interface{}{
		{
			"domain":       "google",
			"only_country": false,
			"only_popular": true,
		},
		{
			"domain":       "google",
			"only_country": true,
			"only_popular": false,
		},
		{
			"domain":       "google",
			"only_country": true,
			"only_popular": true,
		},
		{
			"domain":       "google",
			"only_country": false,
			"only_popular": false,
		},
	}

	for _, domain := range domains {

		t.Run(domain["domain"].(string), func(t *testing.T) {
			allTLDs := GetAllTLDs(domain["domain"].(string), false, false)

			for _, tld := range allTLDs {
				if !domain["only_popular"].(bool) {
					// check if there are tlds from country in list
					if strings.HasSuffix(tld, ".us") {
						t.Errorf("Expected %s to not have a country TLD", tld)
					}
				}

				if !domain["only_country"].(bool) {
					// check if there are tlds from popular in list
					if strings.HasSuffix(tld, ".com") {
						t.Errorf("Expected %s to not have a popular TLD", tld)
					}
				}

				if domain["only_popular"].(bool) && domain["only_country"].(bool) {
					// check if there are tlds from popular and country in list
					if strings.HasSuffix(tld, ".us") || strings.HasSuffix(tld, ".com") {
						t.Errorf("Expected %s to not have a popular or country TLD", tld)
					}
				}

				if !domain["only_popular"].(bool) && !domain["only_country"].(bool) {
					// check if there are tlds from popular and country in list
					if strings.HasSuffix(tld, ".us") || strings.HasSuffix(tld, ".com") {
						t.Errorf("Expected %s to not have a popular or country TLD", tld)
					}
				}

			}
		})

	}
}
