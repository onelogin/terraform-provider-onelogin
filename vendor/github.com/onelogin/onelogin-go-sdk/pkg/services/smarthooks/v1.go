package smarthooks

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errSmartHooksV2Context = "smarthooks v1 service"

// V1Service holds the information needed to interface with a repository
type V1Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V1Service {
	return &V1Service{
		Endpoint:     fmt.Sprintf("%s/api/2/hooks", host),
		Repository:   repo,
		ErrorContext: errSmartHooksV2Context,
	}
}

type SmartHookWriteRequest struct {
	*SmartHook
	EnvVars []string `json:"env_vars"`
}

// Query retrieves all the smarthooks from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all smarthooks
func (svc *V1Service) Query(query *SmartHookQuery) ([]SmartHook, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var smarthooks []SmartHook
	for _, bytes := range resp {
		var unmarshalled []SmartHook
		json.Unmarshal(bytes, &unmarshalled)
		smarthooks = append(smarthooks, unmarshalled...)
	}
	return smarthooks, nil
}

// GetOne retrieves the smarthook by id and returns it
func (svc *V1Service) GetOne(id string) (*SmartHook, error) {
	out := SmartHook{}
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

// Create takes a smarthook without an id and attempts to use the parameters to create it
// in the API. Modifies the smarthook in place, or returns an error if one occurs
func (svc *V1Service) Create(smarthook *SmartHook) (*SmartHook, error) {
	out := SmartHook{}
	if smarthook.Function == nil {
		return &out, errors.New("No Function Definition Given")
	}

	envVars := make([]string, len(smarthook.EnvVars))
	for i, v := range smarthook.EnvVars {
		envVars[i] = *v.Name
	}
	writeRequest := SmartHookWriteRequest{
		SmartHook: smarthook,
		EnvVars:   envVars,
	}

	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    writeRequest,
	})
	if err != nil {
		return &out, err
	}

	json.Unmarshal(resp, &out)
	return &out, nil
}

// Update takes a smarthook and an id and attempts to use the parameters to update it
// in the API. Returns a new SmartHook object, or returns an error if one occurs
func (svc *V1Service) Update(smarthook *SmartHook) (*SmartHook, error) {
	out := SmartHook{}
	if smarthook.ID == nil {
		return &out, errors.New("No ID Given")
	}
	if smarthook.Function == nil {
		return &out, errors.New("No Function Definition Given")
	}

	envVars := make([]string, len(smarthook.EnvVars))
	for i, v := range smarthook.EnvVars {
		envVars[i] = *v.Name
	}
	writeRequest := SmartHookWriteRequest{
		SmartHook: smarthook,
		EnvVars:   envVars,
	}

	id := *smarthook.ID
	smarthook.ID = nil
	resp, err := svc.UpdateRaw(id, writeRequest)
	if err != nil {
		return &out, err
	}

	json.Unmarshal(resp, &out)
	return &out, nil
}

// UpdateRaw takes a smarthook and an id and attempts to use the parameters to update it
// in the API. Returns the raw response bytes, or returns an error if one occurs
func (svc *V1Service) UpdateRaw(id string, smarthook interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    smarthook,
	})
}

// Destroy deletes the smarthook with the given id, and if successful, it returns nil
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
