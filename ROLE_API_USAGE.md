# OneLogin Role API Usage Guide

## Overview

This document outlines the proper usage of the OneLogin Role API endpoints when working with the OneLogin Go SDK v4.5.1, especially when managing role users. It provides guidance based on our testing and findings to ensure that the Terraform provider correctly interacts with the OneLogin API.

## Key API Behavior Findings

Our testing revealed important details about how the OneLogin API behaves when managing role users:

1. **Base role APIs do return user information**
   - Both the `GetRoles` and `GetRoleByID` APIs return the array of user IDs in the role
   - This contradicts our initial hypothesis that role users might be stored separately

2. **Standard role update with empty users array does not remove users**
   - Despite the API returning users in the role object, when updating a role with `UpdateRoleWithContext` and providing an empty users array (`[]`), users are NOT removed from the role
   - This inconsistency is important to understand for proper implementation

3. **Specialized user management APIs work correctly**
   - `AddRoleUsers` correctly adds users to a role
   - `DeleteRoleUsers` correctly removes users from a role
   - These specialized APIs should be used for user management operations

## Best Practices for SDK Usage

### Creating Roles

When creating a role with users:
1. Create the role using `CreateRoleWithContext` with all properties, including users
2. The API correctly processes users during initial creation

```go
// Create the role with all properties
roleResp, err := client.CreateRoleWithContext(ctx, &models.Role{
    Name: &roleName,
    Apps: appIDs,
    Admins: adminIDs,
    Users: userIDs, // Include users in the initial creation
})
```

### Reading Roles

When reading role data:
1. Use `GetRoleByIDWithContext` to get all role information, including users
2. The API returns the correct list of user IDs in the role
3. Optionally use `GetRoleUsers` if you need detailed user information (email, name, etc.)

```go
// Get role information including user IDs
roleResp, err := client.GetRoleByIDWithContext(ctx, roleID, nil)
// Handle error and process roleResp
```

### Updating Roles

When updating a role with user changes:
1. Use `UpdateRoleWithContext` for non-user fields (name, apps, admins)
2. For user changes, use the specialized user management APIs:
   a. Use `AddRoleUsers` to add new users (in new list but not in old list)
   b. Use `DeleteRoleUsers` to remove users (in old list but not in new list)
3. Do not rely on sending an empty users array to remove users - it won't work

```go
// Handle user changes
if usersChanged {
    // Add new users
    if len(usersToAdd) > 0 {
        _, err := client.AddRoleUsers(roleID, usersToAdd)
        // Handle error
    }
    
    // Remove old users
    if len(usersToRemove) > 0 {
        _, err := client.DeleteRoleUsers(roleID, usersToRemove)
        // Handle error
    }
}

// Update other role properties
_, err := client.UpdateRoleWithContext(ctx, roleID, &models.Role{
    Name: &roleName,
    Apps: appIDs,
    Admins: adminIDs,
    Users: userIDs, // Including users is fine, but won't affect user membership
})
```

### Deleting Roles

When deleting a role:
1. Use `DeleteRoleWithContext` to delete the role
2. No need to remove users first as they will be automatically disassociated

```go
_, err := client.DeleteRoleWithContext(ctx, roleID)
// Handle error
```

## Technical Details

### API Endpoints

The SDK uses these OneLogin API endpoints:

- `POST /api/2/roles` - Create a role
- `GET /api/2/roles` - List all roles
- `GET /api/2/roles/{id}` - Get role information
- `PUT /api/2/roles/{id}` - Update a role
- `DELETE /api/2/roles/{id}` - Delete a role
- `GET /api/2/roles/{id}/users` - Get detailed user information for a role
- `POST /api/2/roles/{id}/users` - Add users to a role
- `DELETE /api/2/roles/{id}/users` - Remove users from a role

### Model Structure

The OneLogin SDK uses the following model for roles:

```go
// Role represents the Role resource in OneLogin
type Role struct {
    ID     *int32  `json:"id,omitempty"`
    Name   *string `json:"name,omitempty"`
    Admins []int32 `json:"admins"`
    Apps   []int32 `json:"apps"`
    Users  []int32 `json:"users"`
}
```

However, our testing revealed that the `Users` field is only partially processed during updates. When updating a role, emptying the users array does not remove users. The specialized user management APIs must be used for this purpose.

## Conclusion

The OneLogin API has an inconsistency in how it handles role users during updates. While it correctly displays and processes users during creation and reading, it does not properly process empty user arrays during updates. To ensure correct user management:

1. Always use the specialized `AddRoleUsers` and `DeleteRoleUsers` APIs for adding and removing users
2. Compare old and new user sets to determine which users to add or remove
3. Do not rely on empty arrays in `UpdateRoleWithContext` to remove users

This approach ensures that your Terraform provider can correctly manage role users without encountering errors or unexpected behavior.