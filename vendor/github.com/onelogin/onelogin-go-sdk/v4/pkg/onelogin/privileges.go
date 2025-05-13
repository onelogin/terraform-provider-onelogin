package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	PrivilegesPath string = "api/1/privileges"
)

func (sdk *OneloginSDK) ListPrivileges() (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreatePrivilege(privilege models.Privilege) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, privilege)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetPrivilege(privilegeID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdatePrivilege(privilegeID string, privilege models.Privilege) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, privilege)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeletePrivilege(privilegeID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetPrivilegeRoles(privilegeID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "roles")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) AddPrivilegeToRoles(privilegeID string, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "roles")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) RemoveRoleFromPrivilege(privilegeID string, roleID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "roles", roleID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetPrivilegeUsers(privilegeID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) AddPrivilegeToUsers(privilegeID string, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "users")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) RemoveUserFromPrevilege(privilegeID string, userID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(PrivilegesPath, privilegeID, "users", userID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
