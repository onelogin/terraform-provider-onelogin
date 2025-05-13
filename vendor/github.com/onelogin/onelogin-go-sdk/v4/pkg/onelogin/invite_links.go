package onelogin

import (
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	InvitesPath = "api/1/invites"
)

func (sdk *OneloginSDK) GenerateInviteLink(invite mod.Invite) (interface{}, error) {
	p, err := utl.BuildAPIPath(InvitesPath, "get_invite_link")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, invite)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) SendInviteLink(invite mod.Invite) (interface{}, error) {
	p, err := utl.BuildAPIPath(InvitesPath, "send_invite_link")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, invite)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
