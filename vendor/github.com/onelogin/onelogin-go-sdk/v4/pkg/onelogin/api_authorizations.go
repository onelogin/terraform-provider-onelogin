package onelogin

import (
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	APIAuthPath string = "api/2/api_authorizations"
)

// ListAuthServers
func (sdk *OneloginSDK) GetAuthServers(queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// GetAuthServersByID
func (sdk *OneloginSDK) GetAuthServerByID(authID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Create Authorization Server
func (sdk *OneloginSDK) CreateAuthServer(authServer *mod.AuthServer) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, authServer)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// UpdateAuthServerByID
func (sdk *OneloginSDK) UpdateAuthServer(authID int, authServer *mod.AuthServer) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, authServer)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// DeleteAuthServerById
func (sdk *OneloginSDK) DeleteAuthServer(authID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Claim related endpoints
// List Access Token Claims
func (sdk *OneloginSDK) GetAuthClaims(authID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "claims")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateAuthServerClaim(authID int, claim mod.AccessTokenClaim) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "claims")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, claim)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateAuthClaim(authID int, claimID int, claim mod.AccessTokenClaim) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "claims", claimID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, claim)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteAuthClaim(authID, claimID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "claims", claimID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Scopes related endpoints
func (sdk *OneloginSDK) GetAuthServerScopes(authID int, queryParams mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "scopes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, queryParams)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateAuthServerScope(authID int, scope mod.Scope) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "scopes")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, scope)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateAuthServerScope(authID, scopeID int, scope mod.Scope) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "scopes", scopeID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, scope)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteAuthServerScope(authID, scopeID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "scopes", scopeID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Client App related endpoints
func (sdk *OneloginSDK) GetClientApps(authID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "clients")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateClientApp(authID int, clientApp mod.ClientAppRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "clients")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, clientApp)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateClientApp(authID, clientID int, clientApp mod.ClientAppRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "clients", clientID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, clientApp)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteClientApp(authID, clientID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(APIAuthPath, authID, "clients", clientID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
