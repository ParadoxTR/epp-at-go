# EPP-AT-GO

[![Go Reference](https://pkg.go.dev/badge/github.com/ParadoxTR/epp-at-go.svg)](https://pkg.go.dev/github.com/ParadoxTR/epp-at-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ParadoxTR/epp-at-go)](https://goreportcard.com/report/github.com/ParadoxTR/epp-at-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust Go library designed specifically for nic.at EPP (Extensible Provisioning Protocol) operations. This module provides comprehensive support for Austrian (.at) domain management through the nic.at registry system. Built for domain registrars and registry operators who need reliable, type-safe EPP client functionality tailored to nic.at's specific requirements and extensions.


## Installation

```bash
go get github.com/ParadoxTR/epp-at-go
```

## Quick Start

```go
package main

import (
    "log"
    "time"
    
    "github.com/ParadoxTR/epp-at-go/epp"
)

func main() {
    // Configure your EPP connection
    config := epp.Config{
        Hostname: "epp.nic.at",
        Port:     700,
        Username: "your-username",
        Password: "your-password",
        Timeout:  30 * time.Second,
    }
    
    // Create and connect
    client := epp.NewClient(config)
    defer client.Close()
    
    if err := client.Connect(); err != nil {
        log.Fatal("Connection failed:", err)
    }
    
    if err := client.Login(); err != nil {
        log.Fatal("Login failed:", err)
    }
    defer client.Logout()
    
    // Check domain availability
    domains := []string{"example.at", "test.at"}
    response, err := client.CheckDomain(domains)
    if err != nil {
        log.Fatal("Domain check failed:", err)
    }
    
    for _, domain := range response.ResData.ChkData.Names {
        available := domain.Name.Available == "1"
        log.Printf("Domain %s: available=%t", domain.Name.Name, available)
    }
}
```

## Core Features

### Domain Management

```go
// Check multiple domains at once
response, err := client.CheckDomain([]string{"example.at", "test.at"})

// Register a new domain
domain := epp.Domain{
    Name:        "example.at",
    Period:      &epp.Period{Unit: "y", Value: 1},
    Nameservers: []string{"ns1.example.com", "ns2.example.com"},
    Registrant:  "contact-123",
    Contacts: []epp.DomainContact{
        {Type: "admin", ID: "admin-contact"},
        {Type: "tech", ID: "tech-contact"},
    },
    AuthInfo: "secure-auth-code",
}
createResp, err := client.CreateDomain(domain)

// Transfer a domain
transferResp, err := client.TransferRequestDomain("example.at", "auth-code")
```

### Contact Management

```go
// Create a contact (Austrian EPP style)
contact := epp.Contact{
    ID: "AUTO", // Austrian EPP auto-generates IDs
    PostalInfo: epp.ContactPostalInfo{
        Type: "int",
        Name: "Max Mustermann",
        Addr: epp.ContactAddr{
            Street: []string{"Musterstraße 123"},
            City:   "Wien",
            PC:     "1010",
            CC:     "AT",
        },
    },
    Voice: "+43.15551234567",
    Email: "max@example.at",
    AuthInfo: epp.ContactAuthInfo{Pw: "contact-password"},
    Type: "privateperson", // Austrian extension
}
createResp, err := client.CreateContact(contact)
```

### DNSSEC Operations

```go
// Add DNSSEC during domain creation
dnssecExt := &epp.DNSSECExtension{
    KeyData: []epp.DNSSECKeyData{
        {
            Flags:     257,
            Protocol:  3,
            Algorithm: 8,
            PubKey:    "your-base64-encoded-key",
        },
    },
}
// Use in domain create/update operations
```

## Error Handling

The library provides comprehensive error handling with EPP-specific codes:

```go
response, err := client.CheckDomain([]string{"example.at"})
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

// Always check EPP result codes
if response.Result.Code != "1000" {
    log.Printf("EPP error %s: %s", response.Result.Code, response.Result.Msg)
    return
}
```

### Common EPP Result Codes

| Code | Description |
|------|-------------|
| 1000 | Command completed successfully |
| 1001 | Command completed successfully; action pending |
| 2001 | Command syntax error |
| 2003 | Required parameter missing |
| 2302 | Object exists |
| 2303 | Object does not exist |

## Configuration

### Environment Variables

For development, you can use environment variables:

```bash
export EPP_HOST=epp.nic.at
export EPP_USERNAME=your-username
export EPP_PASSWORD=your-password
```

### TLS Configuration

The library automatically handles TLS connections with proper certificate validation. For testing environments, you may need to adjust TLS settings.

## Examples

Check the `examples/` directory for complete working examples:

- `basic_usage.go` - Comprehensive example covering all major operations
- Domain registration workflow
- Contact management
- Error handling patterns

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build ./...
```

### Linting

```bash
golangci-lint run
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


## Support

- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/ParadoxTR/epp-at-go/issues)
- **Documentation**: Full API documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/ParadoxTR/epp-at-go)
- **Examples**: See the `examples/` directory for working code samples

---

**Note**: This library is not officially affiliated with nic.at or any domain registry. It's an independent implementation of the EPP protocol with Austrian extensions support.