package onelogin

import (
	// "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	ReportsPath = "api/2/reports"
)

func (sdk *OneloginSDK) GetReports() (interface{}, error) {
	p, err := utl.BuildAPIPath(ReportsPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) RunReport(reportsID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(ReportsPath, reportsID, "run")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) RunReportInBackground(reportsID int, body interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(ReportsPath, reportsID, "run_background")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, body)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
