package smarthookenvs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errEnvVarsV2Context = "envVar environment variables v1 service"

// V1Service holds the information needed to interface with a repository
type V1Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V1Service {
	return &V1Service{
		Endpoint:     fmt.Sprintf("%s/api/2/hooks/envs", host),
		Repository:   repo,
		ErrorContext: errEnvVarsV2Context,
	}
}

// Query retrieves all the envVars from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all envVars
func (svc *V1Service) Query(query *SmartHookEnvVarQuery) ([]EnvVar, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var envVars []EnvVar
	for _, bytes := range resp {
		var unmarshalled []EnvVar
		json.Unmarshal(bytes, &unmarshalled)
		envVars = append(envVars, unmarshalled...)
	}
	return envVars, nil
}

// GetOne retrieves the envVar by id and returns it
func (svc *V1Service) GetOne(id string) (*EnvVar, error) {
	out := EnvVar{}
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return &out, err
	}

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}

	json.Unmarshal(resp[0], &out)
	return &out, nil
}

// Create takes a envVar without an id and attempts to use the parameters to create it
// in the API. Modifies the envVar in place, or returns an error if one occurs
func (svc *V1Service) Create(envVar *EnvVar) (*EnvVar, error) {
	out := EnvVar{}
	if envVar.Name == nil || envVar.Value == nil {
		return &out, errors.New("Name and Value are both required")
	}

	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    envVar,
	})
	if err != nil {
		return &out, err
	}

	json.Unmarshal(resp, &out)
	return &out, nil
}

// Update takes a envVar and an id and attempts to use the parameters to update it
// in the API. Returns a new EnvVar object, or returns an error if one occurs
func (svc *V1Service) Update(envVar *EnvVar) (*EnvVar, error) {
	out := EnvVar{}
	if envVar.Name != nil && envVar.ID == nil { // give a name but no id, we'll try and find it for you
		possible, err := svc.Query(nil)
		if err != nil {
			return &out, errors.New("unable to find by ID or Name")
		}
		for _, p := range possible {
			if *p.Name == *envVar.Name {
				envVar.ID = p.ID
				break
			}
		}
	}
	if envVar.ID == nil {
		return &out, errors.New("no ID or Name given")
	}
	if envVar.Value == nil {
		return &out, errors.New("value is required")
	}

	id := *envVar.ID
	envVar.ID = nil
	envVar.Name = nil
	envVar.CreatedAt = nil
	envVar.UpdatedAt = nil

	resp, err := svc.UpdateRaw(id, envVar)
	if err != nil {
		return &out, err
	}

	json.Unmarshal(resp, &out)
	return &out, nil
}

// UpdateRaw takes a envVar and an id and attempts to use the parameters to update it
// in the API. Modifies the envVar in place, or returns an error if one occurs
func (svc *V1Service) UpdateRaw(id string, envVar interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    envVar,
	})
}

// Destroy deletes the envVar with the given id, and if successful, it returns nil
func (svc *V1Service) Destroy(id string) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}
