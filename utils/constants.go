package utils

import "regexp"

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+$`)

var POPULAR_TLDS = []string{
	"com", "net", "org", "io", "dev", "app", "co",
}

var COUNTRY_TLDS = []string{
	"us", "uk", "ca", "au", "de", "fr", "it", "es", "nl", "jp", "kr", "cn", "in",
	"br", "mx", "ar", "cl", "co", "pe", "ru", "pl", "cz", "ch", "at", "se", "no",
	"dk", "fi", "be", "pt", "gr", "tr", "za", "eg", "ma", "ng", "ke", "co.uk",
}

var NEW_TLDS = []string{
	"xyz", "me", "info", "biz", "ai", "shop", "store", "online",
	"info", "blog", "us", "tech",
}

var ALL_TLDS = append(POPULAR_TLDS, append(COUNTRY_TLDS, NEW_TLDS...)...)
