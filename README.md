# whoisthere

whoisthere is a Model Context Protocol (MCP) server written in Go that enables Large Language Models (LLMs) to check domain name availability and discover available Top-Level Domains (TLDs).

## Tools

This server exposes two main tools:

### DomainAvailable
Checks if a specific domain name with a TLD is available for registration.

- **Input**:
  - `domain` (string): The domain to check (e.g., "example.com"). Can include http/https prefix.
- **Output**:
  - `Available` (boolean): True if available.
  - `Domain` (string): The checked domain.

### AvailableDomainTLDFinder
Given a domain name stem (without TLD), checks availability across multiple TLDs to find open registration options.

- **Input**:
  - `domain` (string): The domain stem to check (e.g., "example" for "example.com"). Should not include a TLD.
  - `only_popular` (boolean, optional): If true, only checks popular TLDs (com, net, org, etc.).
  - `only_country` (boolean, optional): If true, only checks country-code TLDs (us, uk, ca, etc.).
- **Output**:
  - `AvailableDomains` (array of strings): List of full domain names available for registration.
  - `PopularDomains` (array of strings): List of available domains from the "popular" category.

## Installation & Usage

### Prerequisites
- Go 1.23 or higher
- Make

### Build
To build the project:

```bash
make build
```

### Docker

You can also build and run the server using Docker.

**Build the image:**
```bash
docker build -t whoisthere .
```

**Run the container:**
```bash
docker run -p 8080:8080 whoisthere
```

You can also use the pre-built image from GitHub Container Registry:

```bash
docker run -p 8080:8080 ghcr.io/jaysongiroux/whoisthere:latest
```

### Run
To run the server locally (without Docker):

```bash
make run
```

By default, the server listens on `localhost:8080`.

### Configuration
You can configure the host and port using either command-line flags or environment variables.

**Using Flags:**
```bash
./bin/whoisthere --host localhost:9000
```

**Using Environment Variables:**
```bash
HOST=localhost:9000 ./bin/whoisthere
```

## Integration
Add this server to your MCP client configuration (e.g., `mcp.json`):

```json
{
  "mcpServers": {
    "whoisthere": {
      "url": "http://localhost:8080/sse"
    }
  }
}
```
