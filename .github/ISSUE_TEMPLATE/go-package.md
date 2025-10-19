---
name: Go Package Issue
about: Report an issue with the Go package
title: '[GO] '
labels: ['go-package', 'bug']
assignees: ''

---

## Go Package Issue

**Package Version:**
<!-- Which version of the package are you using? -->
```
go list -m github.com/myerscode/aws-meta
```

**Go Version:**
<!-- What version of Go are you using? -->
```
go version
```

**Issue Description:**
<!-- A clear and concise description of what the issue is -->

**Code Example:**
<!-- Please provide a minimal code example that reproduces the issue -->
```go
package main

import (
    "fmt"
    "github.com/myerscode/aws-meta/pkg/partitions"
    "github.com/myerscode/aws-meta/pkg/services"
)

func main() {
    // Your code here
}
```

**Expected Behavior:**
<!-- What did you expect to happen? -->

**Actual Behavior:**
<!-- What actually happened? -->

**Error Output:**
<!-- If there's an error, please include the full error message -->
```
```

**Additional Context:**
<!-- Add any other context about the problem here -->

**Environment:**
- OS: <!-- e.g., Ubuntu 20.04, macOS 12.0, Windows 10 -->
- Architecture: <!-- e.g., amd64, arm64 -->
- Go modules enabled: <!-- yes/no -->