package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jaysongiroux/whoisthere/utils"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type AvailableDomainTLDInput struct {
	Domain      string `json:"domain"                 jsonschema:"the domain to check availability for. Should not start with http or https. DO not include TLD. Will include a list of all available domains with TLD. TLD options can be refined with only_popular or only_country"`
	OnlyPopular bool   `json:"only_popular,omitempty" jsonschema:"if true, only check popular TLDs"`
	OnlyCountry bool   `json:"only_country,omitempty" jsonschema:"if true, only check country TLDs"`
}

type SingleDomainAvailableInput struct {
	Domain string `json:"domain" jsonschema:"the domain to check availability for. Can start with http:// or https:// or be empty. TLD is optional, if not provided, all TLDs will be checked but can be refined by only_popular or only_country"`
}

type SingleDomainAvailableOutput struct {
	Available bool
	Domain    string
}

type AvailableDomainTLDOutput struct {
	AvailableDomains []string `json:"List of Domain that are available to be registered"`
	PopularDomains   []string `json:"List of popular for the supplied domain that are available to be registered"`
}

type Config struct {
	Host string
}

func checkDomainAvailabilityWhoIs(domain string) bool {
	/**
	* Check if a domain is available via WHOIS
	* returns true if domain is available
	* returns false if domain is not available
	 */
	start := time.Now()
	raw, err := whois.Whois(domain)
	end := time.Now()
	fmt.Printf("Time taken to get WHOIS data: %v\n", end.Sub(start))
	if err != nil {
		return false
	}

	// Check for common "not found" patterns
	lower := strings.ToLower(raw)
	notFoundPatterns := []string{
		"no match for",
		"not found",
		"no entries found",
		"no data found",
		"domain not found",
		"no information available",
		"status: free",
		"status: available",
	}

	for _, pattern := range notFoundPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	// Try parsing
	parsed, err := whoisparser.Parse(raw)
	if errors.Is(err, whoisparser.ErrNotFoundDomain) {
		return true
	}
	if err != nil {
		return false
	}

	// Check for expired status flags
	expiredStatuses := []string{
		"pendingdelete",
		"redemptionperiod",
		"expired",
	}
	for _, s := range parsed.Domain.Status {
		statusLower := strings.ToLower(s)
		for _, exp := range expiredStatuses {
			if strings.Contains(statusLower, exp) {
				return true
			}
		}
	}

	// Check if expiry date has passed
	if parsed.Domain.ExpirationDate != "" {
		expirationDate, err := utils.ParseExpirationDate(parsed.Domain.ExpirationDate)
		if err != nil {
			return false
		}
		t := expirationDate.In(time.UTC)
		if t.Before(time.Now()) {
			return true
		}
	}

	return false
}

func CheckAvailableDomainTLD(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AvailableDomainTLDInput,
) (
	*mcp.CallToolResult,
	AvailableDomainTLDOutput,
	error,
) {
	hasTld := utils.HasTLD(input.Domain)
	if hasTld {
		return nil, AvailableDomainTLDOutput{}, errors.New(
			"domain must not contain a TLD. This is used to find available TLDs for a domain",
		)
	}

	// add a temporary TLD to the end in order to validate the domain
	_, err := utils.CleanAndValidateDomain(input.Domain + ".com")
	if err != nil {
		return nil, AvailableDomainTLDOutput{}, err
	}

	// these can be empty, if they are not provided they are set to false
	onlyCountry := input.OnlyCountry
	onlyPopular := input.OnlyPopular

	domains := utils.GetAllTLDs(input.Domain, onlyCountry, onlyPopular)

	type result struct {
		domain    string
		available bool
	}

	results := make(chan result, len(domains))
	var wg sync.WaitGroup
	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			available := checkDomainAvailabilityWhoIs(domain)

			results <- result{domain, available}
		}(domain)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	//	format result into structured output
	output := AvailableDomainTLDOutput{
		AvailableDomains: []string{},
		PopularDomains:   []string{},
	}

	for result := range results {
		if result.available {
			output.AvailableDomains = append(output.AvailableDomains, result.domain)

			tld := utils.GetTLD(result.domain)
			isPopular := utils.IsPopularTLD(tld)
			if isPopular {
				output.PopularDomains = append(output.PopularDomains, result.domain)
			}
		}

	}
	return &mcp.CallToolResult{}, output, nil
}

func CheckDomainAvailability(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SingleDomainAvailableInput,
) (
	*mcp.CallToolResult,
	SingleDomainAvailableOutput,
	error,
) {
	hasTld := utils.HasTLD(input.Domain)
	if !hasTld {
		return nil, SingleDomainAvailableOutput{}, errors.New(
			"domain must contain a TLD. Please try again with a TLD included or call AvailableDomainTLDFinder to find available TLDs for a domain",
		)
	}

	if input.Domain == "" {
		return nil, SingleDomainAvailableOutput{}, errors.New("domain cannot be empty")
	}

	validatedDomain, err := utils.CleanAndValidateDomain(input.Domain)
	if err != nil {
		return nil, SingleDomainAvailableOutput{}, err
	}
	input.Domain = validatedDomain

	available := checkDomainAvailabilityWhoIs(input.Domain)

	return &mcp.CallToolResult{}, SingleDomainAvailableOutput{
		Available: available,
		Domain:    input.Domain,
	}, nil
}

func getArgs() Config {
	/**
	 * gets args from terminal or env variable
	 * Flags:
	 *  host: host to run server on (default :8080)
	 *
	 * Hierarchy
	 *  env > flag
	 */
	defaultHost := "localhost:8080"

	// Use flags to get command-line arguments
	var host string

	flag.StringVar(&host, "host", defaultHost, "Host to run server on (default: 'localhost:8080')")
	flag.Parse()

	envHost := os.Getenv("HOST")
	if envHost != "" {
		host = envHost
	}

	return Config{
		Host: host,
	}
}

func main() {
	config := getArgs()

	// Create a server with a single tool.
	log.Print("Starting MCP Server...")

	server := mcp.NewServer(
		&mcp.Implementation{Name: "whoisthere", Version: "v1.0.0"},
		&mcp.ServerOptions{
			HasTools: true,
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "DomainAvailable", Description: "Check if a domain is available"},
		CheckDomainAvailability,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "AvailableDomainTLDFinder",
			Description: "Given a string without a TLD, find all available TLDs for a potential domain",
		},
		CheckAvailableDomainTLD,
	)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	log.Printf("MCP server listening on %s", config.Host)

	// nolint:gosec
	if err := http.ListenAndServe(config.Host, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
