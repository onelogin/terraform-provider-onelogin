---
layout: "onelogin"
page_title: "OneLogin: SDK Versions"
sidebar_current: "docs-onelogin-guide-sdk-versions"
description: |-
  Information about the OneLogin Go SDK versions used by this provider
---

# OneLogin SDK Versions

This Terraform provider uses the OneLogin Go SDK to interact with the OneLogin API. Understanding the SDK version can help diagnose compatibility issues and understand feature availability.

## Current SDK Version

As of version 0.8.1 of the provider, we use OneLogin Go SDK v4.5.0.

## SDK Version History

| Provider Version | SDK Version | Key Features |
|------------------|-------------|--------------|
| 0.8.1 | v4.5.0 | Context support for role operations, improved role updates |
| 0.7.0 | v4.4.0 | User mappings support |
| < 0.7.0 | < v4.4.0 | Earlier functionality |

## SDK v4.5.0 Features

SDK v4.5.0 introduces several important improvements:

1. **Context Support for Role Operations**:
   - All role operations now accept a context parameter (`ctx context.Context`)
   - This enables better request handling, cancellation, and deadline management
   - Methods include: `CreateRoleWithContext`, `GetRolesWithContext`, `GetRoleByIDWithContext`, `UpdateRoleWithContext`, `DeleteRoleWithContext`

2. **Role Model Changes**:
   - `ID` field changed from `int64` to `*int32`
   - `Name` field is now a pointer type
   - Added `Admins`, `Apps`, and `Users` as int32 slices
   - Improved pointer handling for optional fields

3. **Method Signature Updates**:
   - `UpdateRole` now requires a pointer parameter
   - Removed unused `queryParams` from `DeleteRole`

## Using Context-Aware Methods

If you're developing custom integrations with the OneLogin API using the Go SDK directly, you can take advantage of the context-aware methods:

```go
import (
    "context"
    "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
    "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

func example() {
    client := onelogin.OneloginSDK{}
    ctx := context.Background()
    
    // Create a role with context
    role := &models.Role{
        Name: "Example Role",
    }
    
    result, err := client.CreateRoleWithContext(ctx, role)
    if err != nil {
        // Handle error
    }
    
    // Use the result
}
```

The Terraform provider automatically uses these context-aware methods to ensure optimal integration with Terraform's context management.