package roles

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errRolesV1Context = "roles v1 service"

// V1Service holds the information needed to interface with a repository
type V1Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v1.
func New(repo services.Repository, host string) *V1Service {
	return &V1Service{
		Endpoint:     fmt.Sprintf("%s/api/2/roles", host),
		Repository:   repo,
		ErrorContext: errRolesV1Context,
	}
}

// Query retrieves all the roles from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all roles
func (svc *V1Service) Query(query *RoleQuery) ([]Role, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var roles []Role
	for _, bytes := range resp {
		var unmarshalled []Role
		json.Unmarshal(bytes, &unmarshalled)
		roles = append(roles, unmarshalled...)
	}
	return roles, nil
}

// GetOne retrieves the role by id and returns it
func (svc *V1Service) GetOne(id int32) (*Role, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}
	var role Role

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}

	json.Unmarshal(resp[0], &role)
	return &role, nil
}

// Create takes a role without an id and attempts to use the parameters to create it
// in the API. Modifies the role in place, or returns an error if one occurs
func (svc *V1Service) Create(role *Role) error {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    role,
	})
	if err != nil {
		return err
	}
	json.Unmarshal(resp, role)
	return nil
}

// Update takes a role and an id and attempts to use the parameters to update it
// in the API. Modifies the role in place, or returns an error if one occurs
func (svc *V1Service) Update(role *Role) error {
	if role.ID == nil {
		return errors.New("No ID Given")
	}

	id := *role.ID
	role.ID = nil
	resp, err := svc.UpdateRaw(id, role)

	if err != nil {
		return err
	}

	json.Unmarshal(resp, role)
	return nil
}

// UpdateRaw takes a role and an id and attempts to use the parameters to update it
// in the API. Returned the raw response bytes, or returns an error if one occurs
func (svc *V1Service) UpdateRaw(id int32, role interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    role,
	})
}

// Destroy deletes the role with the given id, and if successful, it returns nil
func (svc *V1Service) Destroy(id int32) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}
