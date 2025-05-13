package onelogin

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	UserPathV1 string = "api/1/users"
	UserPathV2 string = "api/2/users"
)

// Users V2
// was ListUsers
func (sdk *OneloginSDK) GetUsers(query mod.Queryable) (interface{}, error) {
	return sdk.GetUsersWithContext(context.Background(), query)
}

// GetUsersWithContext retrieves users using the provided context
func (sdk *OneloginSDK) GetUsersWithContext(ctx context.Context, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2)
	if err != nil {
		return nil, err
	}
	// Validate query parameters
	validators := query.GetKeyValidators()
	if !utl.ValidateQueryParams(query, validators) {
		return nil, errors.New("invalid query parameters")
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetUsersWithPagination retrieves users with pagination information
func (sdk *OneloginSDK) GetUsersWithPagination(query mod.Queryable) (*mod.PagedResponse, error) {
	return sdk.GetUsersWithPaginationWithContext(context.Background(), query)
}

// GetUsersWithPaginationWithContext retrieves users with pagination information using the given context
func (sdk *OneloginSDK) GetUsersWithPaginationWithContext(ctx context.Context, query mod.Queryable) (*mod.PagedResponse, error) {
	p, err := utl.BuildAPIPath(UserPathV2)
	if err != nil {
		return nil, err
	}

	// Validate query parameters
	validators := query.GetKeyValidators()
	if !utl.ValidateQueryParams(query, validators) {
		return nil, errors.New("invalid query parameters")
	}

	// Make the API request with context
	resp, err := sdk.Client.GetWithContext(ctx, &p, query)
	if err != nil {
		return nil, err
	}

	// Extract data from response
	data, err := utl.CheckHTTPResponse(resp)
	if err != nil {
		return nil, err
	}

	// Extract pagination information from headers
	pagination := mod.PaginationInfo{
		Cursor:       resp.Header.Get("Cursor"),
		AfterCursor:  resp.Header.Get("After-Cursor"),
		BeforeCursor: resp.Header.Get("Before-Cursor"),
	}

	// Try to parse total pages and current page
	if totalPages := resp.Header.Get("Total-Pages"); totalPages != "" {
		if i, err := strconv.Atoi(totalPages); err == nil {
			pagination.TotalPages = i
		}
	}

	if currentPage := resp.Header.Get("Current-Page"); currentPage != "" {
		if i, err := strconv.Atoi(currentPage); err == nil {
			pagination.CurrentPage = i
		}
	}

	if totalCount := resp.Header.Get("Total-Count"); totalCount != "" {
		if i, err := strconv.Atoi(totalCount); err == nil {
			pagination.TotalCount = i
		}
	}

	// Combine data and pagination info
	return &mod.PagedResponse{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func (sdk *OneloginSDK) GetUserByID(id int, queryParams mod.Queryable) (interface{}, error) {
	return sdk.GetUserByIDWithContext(context.Background(), id, queryParams)
}

// GetUserByIDWithContext retrieves a user by ID using the provided context
func (sdk *OneloginSDK) GetUserByIDWithContext(ctx context.Context, id int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateUser(user mod.User) (interface{}, error) {
	return sdk.CreateUserWithContext(context.Background(), user)
}

// CreateUserWithContext creates a user using the provided context
func (sdk *OneloginSDK) CreateUserWithContext(ctx context.Context, user mod.User) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PostWithContext(ctx, &p, user)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateUser(id int, user mod.User) (interface{}, error) {
	return sdk.UpdateUserWithContext(context.Background(), id, user)
}

// UpdateUserWithContext updates a user using the provided context
func (sdk *OneloginSDK) UpdateUserWithContext(ctx context.Context, id int, user mod.User) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PutWithContext(ctx, &p, user)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteUser(id int) (interface{}, error) {
	return sdk.DeleteUserWithContext(context.Background(), id)
}

// DeleteUserWithContext deletes a user using the provided context
func (sdk *OneloginSDK) DeleteUserWithContext(ctx context.Context, id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.DeleteWithContext(ctx, &p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetUserApps(id int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, id, "apps")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetCustomAttributes() (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetCustomAttributeByID(id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateCustomAttributes(requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// CreateCustomAttribute creates a new custom attribute with the specified name and shortname.
// This helper method properly wraps the name and shortname in the required "user_field" object
// as expected by the OneLogin API.
func (sdk *OneloginSDK) CreateCustomAttribute(name, shortname string) (interface{}, error) {
	return sdk.CreateCustomAttributeWithContext(context.Background(), name, shortname)
}

// CreateCustomAttributeWithContext creates a new custom attribute with the provided context
func (sdk *OneloginSDK) CreateCustomAttributeWithContext(ctx context.Context, name, shortname string) (interface{}, error) {
	// Use map[string]string for the inner map to ensure consistent types
	payload := map[string]interface{}{
		"user_field": map[string]string{
			"name":      name,
			"shortname": shortname,
		},
	}

	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PostWithContext(ctx, &p, payload)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateCustomAttributes(id int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// UpdateCustomAttribute updates an existing custom attribute with the specified name and shortname.
// This helper method properly wraps the name and shortname in the required "user_field" object
// as expected by the OneLogin API.
func (sdk *OneloginSDK) UpdateCustomAttribute(id int, name, shortname string) (interface{}, error) {
	return sdk.UpdateCustomAttributeWithContext(context.Background(), id, name, shortname)
}

// UpdateCustomAttributeWithContext updates an existing custom attribute with the provided context
func (sdk *OneloginSDK) UpdateCustomAttributeWithContext(ctx context.Context, id int, name, shortname string) (interface{}, error) {
	// Use map[string]string for the inner map to ensure consistent types
	payload := map[string]interface{}{
		"user_field": map[string]string{
			"name":      name,
			"shortname": shortname,
		},
	}

	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PutWithContext(ctx, &p, payload)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteCustomAttributes(id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV2, "custom_attributes", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetUsersModels(query mod.Queryable) ([]mod.User, error) {
	p, err := utl.BuildAPIPath(UserPathV2)
	if err != nil {
		return nil, err
	}

	// Validate query parameters
	validators := query.GetKeyValidators()
	if !utl.ValidateQueryParams(query, validators) {
		return nil, errors.New("invalid query parameters")
	}

	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}

	tmp, err := utl.CheckHTTPResponse(resp)
	if err != nil {
		return nil, err
	}

	var users []mod.User
	tmpBytes, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tmpBytes, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Users V1

func (sdk *OneloginSDK) GetUserRoles(id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, id, "roles")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) AddRolesForUser(userID int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, userID, "add_roles")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) RemoveRolesForUser(userID int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, userID, "remove_roles")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdatePasswordInsecure(id int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, "set_password_clear_text", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdatePasswordSecure(id int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, "set_password_using_salt", id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) SetCustomAttributes(userID int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, userID, "set_custom_attributes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) SetUserState(userID, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, userID, "set_state")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) LogOutUser(userID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, userID, "logout")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) LockUserAccount(id int, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserPathV1, id, "lock_user")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
