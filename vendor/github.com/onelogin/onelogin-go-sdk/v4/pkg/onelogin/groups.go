package onelogin

import (
	"context"

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

func (sdk *OneloginSDK) GetGroupByIDV2(groupID int) (interface{}, error) {
	return sdk.GetGroupByIDV2WithContext(context.Background(), groupID)
}

func (sdk *OneloginSDK) GetGroupByIDV2WithContext(ctx context.Context, groupID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.GetWithContext(ctx, &p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateGroup(group *mod.Group) (interface{}, error) {
	return sdk.CreateGroupWithContext(context.Background(), group)
}

func (sdk *OneloginSDK) CreateGroupWithContext(ctx context.Context, group *mod.Group) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PostWithContext(ctx, &p, group)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateGroup(groupID int, group *mod.Group) (interface{}, error) {
	return sdk.UpdateGroupWithContext(context.Background(), groupID, group)
}

func (sdk *OneloginSDK) UpdateGroupWithContext(ctx context.Context, groupID int, group *mod.Group) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.PutWithContext(ctx, &p, group)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteGroup(groupID int) (interface{}, error) {
	return sdk.DeleteGroupWithContext(context.Background(), groupID)
}

func (sdk *OneloginSDK) DeleteGroupWithContext(ctx context.Context, groupID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(GroupsV2Path, groupID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.DeleteWithContext(ctx, &p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
