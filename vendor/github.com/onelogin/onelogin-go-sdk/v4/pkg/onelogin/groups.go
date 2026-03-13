package onelogin

import (
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	GroupsPath   = "api/1/groups"
	GroupsV2Path = "api/2/groups"
)

func (sdk *OneloginSDK) GetGroups() (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetGroupByID(groupID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsPath, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateGroup(group *mod.Group) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, group)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateGroup(groupID int, group *mod.Group) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, group)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteGroup(groupID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
