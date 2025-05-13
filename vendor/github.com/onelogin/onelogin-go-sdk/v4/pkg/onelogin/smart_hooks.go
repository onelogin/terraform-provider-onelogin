package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	SmartHooksPath string = "api/2/hooks"
)

func (sdk *OneloginSDK) ListHooks(query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetHook(hookID string, query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, hookID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetHookLogs(hookID string, query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, hookID, "logs")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateHook(hook models.SmartHook) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, hook)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateSmartHook(hookID string, hook models.SmartHook) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, hookID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, hook)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteHook(hookID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, hookID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) ListEnvironmentVariables() (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, "envs")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetEnvironmentVariable(envVarID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, "envs", envVarID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateEnvironmentVariable(requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, "envs")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateEnvironmentVariable(envVarID string, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, "envs", envVarID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteEnvironmentVariable(envVarID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(SmartHooksPath, "envs", envVarID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
