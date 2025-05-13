package onelogin

import (
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
