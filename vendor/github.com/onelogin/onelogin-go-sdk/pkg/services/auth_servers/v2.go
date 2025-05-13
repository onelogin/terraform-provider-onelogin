package authservers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errAuthServersV2Context = "auth_servers v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V2Service {
	return &V2Service{
		Endpoint:     fmt.Sprintf("%s/api/2/api_authorizations", host),
		Repository:   repo,
		ErrorContext: errAuthServersV2Context,
	}
}

// Query retrieves all the auth_servers from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all auth_servers
func (svc *V2Service) Query(query *AuthServerQuery) ([]AuthServer, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var auth_servers []AuthServer
	for _, bytes := range resp {
		var unmarshalled []AuthServer
		json.Unmarshal(bytes, &unmarshalled)
		auth_servers = append(auth_servers, unmarshalled...)
	}

	return auth_servers, nil
}

// GetOne retrieves the authServer by id and returns it
func (svc *V2Service) GetOne(id int32) (*AuthServer, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}
	var authServer AuthServer
	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	json.Unmarshal(resp[0], &authServer)
	return &authServer, nil
}

// Create takes a authServer without an id and attempts to use the parameters to create it
// in the API. Modifies the authServer in place, or returns an error if one occurs
func (svc *V2Service) Create(authServer *AuthServer) error {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    authServer,
	})
	if err != nil {
		return err
	}
	respObj := map[string]int32{}
	json.Unmarshal(resp, &respObj)
	authServer.ID = oltypes.Int32(respObj["id"])
	return nil
}

// Update takes a authServer and an id and attempts to use the parameters to update it
// in the API. Modifies the authServer in place, or returns an error if one occurs
func (svc *V2Service) Update(authServer *AuthServer) error {
	if authServer.ID == nil {
		return errors.New("No ID Given")
	}
	_, err := svc.UpdateRaw(*authServer.ID, authServer)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRaw takes a authServer and an id and attempts to use the parameters to update it
// in the API. Returns the raw response, or an error if one occurs
func (svc *V2Service) UpdateRaw(id int32, authServer interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    authServer,
	})
}

// Destroy deletes the authServer with the given id, and if successful, it returns nil
func (svc *V2Service) Destroy(id int32) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}
