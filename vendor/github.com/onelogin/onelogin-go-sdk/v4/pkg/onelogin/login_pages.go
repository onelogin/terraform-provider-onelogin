package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	LoginPath = "api/1/login"
)

func (sdk *OneloginSDK) CreateSessionLoginToken(requestBody models.CreateSessionLoginRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(LoginPath, "auth")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) VerifyFactor(requestBody models.VerifyFactorRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(LoginPath, "verify_factor")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
