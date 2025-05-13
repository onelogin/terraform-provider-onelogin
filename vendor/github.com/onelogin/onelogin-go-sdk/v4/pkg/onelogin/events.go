package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	EventsPath string = "api/1/events"
)

func (sdk *OneloginSDK) ListEvents(query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(EventsPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetEvents(eventID int, query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(EventsPath, eventID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetEventTypes(query models.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(EventsPath, "types")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
