# Enhanced Environment Configuration

[![Go Reference](https://pkg.go.dev/badge/github.com/legrch/envconfig.svg)](https://pkg.go.dev/github.com/legrch/envconfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/legrch/envconfig)](https://goreportcard.com/report/github.com/legrch/envconfig)
[![License](https://img.shields.io/github/license/legrch/envconfig)](LICENSE)
[![Release](https://img.shields.io/github/v/release/legrch/envconfig)](https://github.com/legrch/envconfig/releases)

A drop-in replacement for [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) that provides detailed error messages for missing environment variables.

## Overview

This package enhances the standard `envconfig` package by providing more detailed error messages when environment variables are missing. It's particularly useful in applications with complex configuration structures where identifying missing variables can be challenging.

## Features

- **Drop-in replacement** for `kelseyhightower/envconfig`
- **Detailed error messages** for missing environment variables
- **Identifies specific fields** that are missing their environment values
- **Simple API** that matches the original library
- **Zero additional dependencies** beyond the original library

## Installation

```bash
go get github.com/legrch/envconfig
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/legrch/envconfig"
)

type Config struct {
	Database struct {
		Host     string `envconfig:"DB_HOST" required:"true"`
		Port     int    `envconfig:"DB_PORT" required:"true"`
		User     string `envconfig:"DB_USER" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
		Name     string `envconfig:"DB_NAME" required:"true"`
	}
	Server struct {
		Host string `envconfig:"SERVER_HOST" required:"true"`
		Port int    `envconfig:"SERVER_PORT" required:"true"`
	}
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error processing environment variables: %v", err)
	}
	
	fmt.Printf("Configuration loaded successfully: %+v\n", cfg)
}
```

### Error Messages

When a required environment variable is missing, the error message includes both the field path and environment variable name:

```
missing required environment variable "DB_HOST" for field "Database.Host"
```

This makes it much easier to track down configuration issues in complex applications.

## Documentation

For detailed documentation and API reference, please visit:

- [Package Documentation](https://pkg.go.dev/github.com/legrch/envconfig)
- [Original envconfig Documentation](https://github.com/kelseyhightower/envconfig) - Base library this package enhances

## Related Documentation

### Project Documentation
- [Detailed Usage](pkg/envconfig/README.md) - More detailed usage examples and documentation

### Official Documentation
- [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) - Original envconfig library
- [Go Documentation](https://golang.org/doc/) - Official Go documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 