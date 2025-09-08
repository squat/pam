[![Go Reference](https://pkg.go.dev/badge/github.com/squat/pam.svg)](https://pkg.go.dev/github.com/squat/pam)
[![Go Report Card](https://goreportcard.com/badge/github.com/squat/pam)](https://goreportcard.com/report/github.com/squat/pam)
[![Build Status](https://github.com/squat/pam/actions/workflows/pre-commit.yml/badge.svg?branch=main)](https://github.com/squat/pam/actions/workflows/pre-commit.yml)

# PAM

PAM provides bindings for writing custom PAM service modules in Golang.
This allows you to easily write modules that:
* authenticate users using your desired mechanism, e.g. OIDC;
* enforce custom account validation logic, e.g. the given user has a certain permission in a database, or the given SSH public key exists in a certain object storage bucket;
* etc.

## Getting Started

Implementing a custom PAM service module requires implementing the [`ServiceModule` interface](/service_module.go) and registering the implementation with the `pamc` package in your main package:
```go
package main

import "github.com/squat/pam/c"

func init() {
    pamc.Register(...)
}

func main() {}
```

Once implemented, the module can be built with:
```shell
go build -buildmode=c-shared -o pam_go.so
```

Take a look at the [examples directory](/examples) for implementation details.
