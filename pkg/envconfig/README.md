# Enhanced Environment Configuration

A drop-in replacement for `github.com/kelseyhightower/envconfig` that provides detailed error messages for missing environment variables.

## Overview

This package enhances the standard `envconfig` package by providing more detailed error messages when environment variables are missing. It's particularly useful in applications with complex configuration structures where identifying missing variables can be challenging.

## Usage

### Prerequisites
- Go 1.18 or later
- `github.com/kelseyhightower/envconfig`

### Examples

```go
package main

import (
    "fmt"
    "log"

    "github.com/legrch/envconfig"
)

type Config struct {
    Database struct {
        Host     string `envconfig:"HOST" required:"true"`
        Port     int    `envconfig:"PORT" required:"true"`
        User     string `envconfig:"USER" required:"true"`
        Password string `envconfig:"PASSWORD" required:"true"`
        Name     string `envconfig:"NAME" required:"true"`
    } `envconfig:"DB"`
    Server struct {
        Host string `envconfig:"HOST" required:"true"`
        Port int    `envconfig:"PORT" required:"true"`
    } `envconfig:"SERVER"`
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

This makes it much easier to track down which configuration fields are missing their environment values.

## Feature Comparison

| Feature                             | Standard envconfig | Enhanced envconfig |
|-------------------------------------|-------------------|-------------------|
| Basic environment variable parsing  | ✅                | ✅                |
| Required field validation           | ✅                | ✅                |
| Field path in error messages        | ❌                | ✅                |
| Multiple missing variables in report | ❌                | ✅                |
| API compatibility                   | -                 | 100%              |

## Best Practices

1. Use meaningful environment variable names that reflect their purpose
2. Group related variables with prefixes
3. Document all environment variables in your application README
4. Use required:"true" for mandatory configuration
5. Provide default values for optional configuration

## Related Documentation

- [Original envconfig](https://github.com/kelseyhightower/envconfig) - The base library this package enhances
- [Go Environment Variables](https://golang.org/pkg/os/#Getenv) - Go's standard library for environment variables

## License

Same as the project license. 