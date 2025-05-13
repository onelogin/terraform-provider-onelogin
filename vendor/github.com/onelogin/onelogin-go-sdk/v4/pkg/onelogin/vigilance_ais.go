package onelogin

import (
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	RiskPath = "api/2/risk"
)

func (sdk *OneloginSDK) TrackEvent(requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "events")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetRiskScore(requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "verify")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateRule(rule mod.Rule) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "rules")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, rule)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) ListRules() (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "rules")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetRule(ruleID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateRule(ruleID string, rule mod.Rule) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, rule)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteRule(ruleID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "rules", ruleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetScoreSummary() (interface{}, error) {
	p, err := utl.BuildAPIPath(RiskPath, "scores")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
