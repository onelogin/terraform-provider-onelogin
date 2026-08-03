package onelogin

import (
	"context"

	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

var (
	RolePath string = "api/2/roles"
)

func (sdk *OneloginSDK) CreateRole(role *mod.Role) (interface{}, error) {
	return sdk.CreateRoleWithContext(context.Background(), role)
}

// CreateRoleWithContext creates a role using the provided context
func (sdk *OneloginSDK) CreateRoleWithContext(ctx context.Context, role *mod.Role) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PostWithContext(ctx, &p, role)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetRoles retrieves a list of roles
func (sdk *OneloginSDK) GetRoles(queryParams mod.Queryable) (interface{}, error) {
	return sdk.GetRolesWithContext(context.Background(), queryParams)
}

// GetRolesWithContext retrieves a list of roles using the provided context
func (sdk *OneloginSDK) GetRolesWithContext(ctx context.Context, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetRolesWithPagination retrieves a list of roles with pagination info from response headers
func (sdk *OneloginSDK) GetRolesWithPagination(queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	return sdk.GetRolesWithPaginationAndContext(context.Background(), queryParams)
}

// GetRolesWithPaginationAndContext retrieves a list of roles with pagination info using the provided context
func (sdk *OneloginSDK) GetRolesWithPaginationAndContext(ctx context.Context, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	p, err := utl.BuildAPIPath(RolePath)
	if err != nil {
		return nil, nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, nil, err
	}
	return utl.CheckHTTPResponseWithPagination(resp)
}

func (sdk *OneloginSDK) GetRoleByID(id int, queryParams mod.Queryable) (interface{}, error) {
	return sdk.GetRoleByIDWithContext(context.Background(), id, queryParams)
}

// GetRoleByIDWithContext retrieves a role by ID using the provided context
func (sdk *OneloginSDK) GetRoleByIDWithContext(ctx context.Context, id int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateRole(id int, role *mod.Role) (interface{}, error) {
	return sdk.UpdateRoleWithContext(context.Background(), id, role)
}

// UpdateRoleWithContext updates a role using the provided context
func (sdk *OneloginSDK) UpdateRoleWithContext(ctx context.Context, id int, role *mod.Role) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PutWithContext(ctx, &p, role)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteRole(id int) (interface{}, error) {
	return sdk.DeleteRoleWithContext(context.Background(), id)
}

// DeleteRoleWithContext deletes a role using the provided context
func (sdk *OneloginSDK) DeleteRoleWithContext(ctx context.Context, id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.DeleteWithContext(ctx, &p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// was ListRoleUsers
func (sdk *OneloginSDK) GetRoleUsers(roleID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetRoleUsersWithPagination retrieves a role's users with pagination info from response headers
func (sdk *OneloginSDK) GetRoleUsersWithPagination(roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	return sdk.GetRoleUsersWithPaginationAndContext(context.Background(), roleID, queryParams)
}

// GetRoleUsersWithPaginationAndContext retrieves a role's users with pagination info using the provided context.
// This endpoint is paginated: callers should follow PaginationInfo.AfterCursor until it is empty to
// retrieve every user, otherwise only the first page is returned.
func (sdk *OneloginSDK) GetRoleUsersWithPaginationAndContext(ctx context.Context, roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "users")
	if err != nil {
		return nil, nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, nil, err
	}
	return utl.CheckHTTPResponseWithPagination(resp)
}

func (sdk *OneloginSDK) AddRoleUsers(roleID int, users []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, users)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// was removeRoleUsers
func (sdk *OneloginSDK) DeleteRoleUsers(roleID int, users []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.DeleteWithBody(&p, users)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetRoleAdmins(roleID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "admins")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetRoleAdminsWithPagination retrieves a role's admins with pagination info from response headers
func (sdk *OneloginSDK) GetRoleAdminsWithPagination(roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	return sdk.GetRoleAdminsWithPaginationAndContext(context.Background(), roleID, queryParams)
}

// GetRoleAdminsWithPaginationAndContext retrieves a role's admins with pagination info using the provided context.
// This endpoint is paginated: callers should follow PaginationInfo.AfterCursor until it is empty to
// retrieve every admin, otherwise only the first page is returned.
func (sdk *OneloginSDK) GetRoleAdminsWithPaginationAndContext(ctx context.Context, roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "admins")
	if err != nil {
		return nil, nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, nil, err
	}
	return utl.CheckHTTPResponseWithPagination(resp)
}

func (sdk *OneloginSDK) AddRoleAdmins(roleID int, admins []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "admins")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, admins)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// was removeRoleAdmins
func (sdk *OneloginSDK) DeleteRoleAdmins(roleID int, admins []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "admins")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.DeleteWithBody(&p, admins)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetRoleApps(roleID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "apps")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetRoleAppsWithPagination retrieves a role's apps with pagination info from response headers
func (sdk *OneloginSDK) GetRoleAppsWithPagination(roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	return sdk.GetRoleAppsWithPaginationAndContext(context.Background(), roleID, queryParams)
}

// GetRoleAppsWithPaginationAndContext retrieves a role's apps with pagination info using the provided context.
// This endpoint is paginated: callers should follow PaginationInfo.AfterCursor until it is empty to
// retrieve every app, otherwise only the first page is returned.
func (sdk *OneloginSDK) GetRoleAppsWithPaginationAndContext(ctx context.Context, roleID int, queryParams mod.Queryable) (interface{}, *mod.PaginationInfo, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "apps")
	if err != nil {
		return nil, nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, nil, err
	}
	return utl.CheckHTTPResponseWithPagination(resp)
}

// was setRoleApps
func (sdk *OneloginSDK) UpdateRoleApps(roleID int, apps []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(RolePath, roleID, "apps")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, apps)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
