package onelogin

import (
	"context"
	"fmt"
	"strconv"
	
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	AppPath string = "api/2/apps"
)

func (sdk *OneloginSDK) CreateApp(app mod.App) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, app)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// List Apps
func (sdk *OneloginSDK) GetApps(queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Get an App
func (sdk *OneloginSDK) GetAppByID(id int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateApp(id int, app mod.App) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, app)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Delete an App Parameter
func (sdk *OneloginSDK) DeleteAppParameter(id int, parameterID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, id, "parameters", parameterID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteApp(id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, id)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetAppUsers(appID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetAppUsersWithPagination retrieves users for an app with optional pagination parameters
func (sdk *OneloginSDK) GetAppUsersWithPagination(appID int, queryParams *mod.AppUserQuery) (*mod.PagedResponse, error) {
	return sdk.GetAppUsersWithPaginationAndContext(context.Background(), appID, queryParams)
}

// GetAppUsersWithPaginationAndContext retrieves users for an app with optional query parameters using the provided context
// Returns pagination metadata along with the user data
func (sdk *OneloginSDK) GetAppUsersWithPaginationAndContext(ctx context.Context, appID int, queryParams *mod.AppUserQuery) (*mod.PagedResponse, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "users")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}
	
	// Validate query parameters if provided
	if queryParams != nil {
		validators := queryParams.GetKeyValidators()
		if !utl.ValidateQueryParams(queryParams, validators) {
			return nil, fmt.Errorf("invalid query parameters: validation failed")
		}
	}
	
	// Make the API request with context
	resp, err := sdk.Client.GetWithContext(ctx, &p, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get app users: %w", err)
	}

	// Extract data from response
	data, err := utl.CheckHTTPResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to process response: %w", err)
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

// App Rules APIs
// list all rules for an app
func (sdk *OneloginSDK) GetAppRules(appID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// get rule by ruleId
func (sdk *OneloginSDK) GetAppRuleByID(appID int, ruleID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateAppRule(appID int, appRule mod.AppRule) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, appRule)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateAppRule(appID, ruleID int, appRule mod.AppRule, queryParams map[string]string) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, appRule)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteAppRule(appID, ruleID int, queryParams map[string]string) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) ListAppRulesConditions(appID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules", "conditions")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetAppRuleOperators(appID int, ruleConditionValue string, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appID, "rules", "conditions", ruleConditionValue, "operators")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetAppRuleConditionValues(appId int, conditionValue string) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appId, "rules", "conditions", conditionValue, "values")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetAppRuleActions(appId int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appId, "rules", "actions")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) ListAppRulesActionValues(appId int, actionValue string) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appId, "rules", "actions", actionValue, "values")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) BulkSortAppRules(appId int, ruleIDs []int) (interface{}, error) {
	p, err := utl.BuildAPIPath(AppPath, appId, "rules", "sort")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, ruleIDs)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
