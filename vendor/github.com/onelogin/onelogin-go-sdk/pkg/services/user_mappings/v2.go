package usermappings

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/onelogin/onelogin-go-sdk/internal/customerrors"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
	"github.com/onelogin/onelogin-go-sdk/pkg/utils"

	"sync"
)

const errUserMappingsV2Context = "user mappings v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
	LegalValuesService     services.SimpleQuery
}

// New creates the new svc service v2.
func New(repo services.Repository, legalValues services.SimpleQuery, host string) *V2Service {
	return &V2Service{
		Endpoint:           fmt.Sprintf("%s/api/2/mappings", host),
		Repository:         repo,
		ErrorContext:       errUserMappingsV2Context,
		LegalValuesService: legalValues,
	}
}

// Query retrieves all the userMappings from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all userMappings
func (svc *V2Service) Query(query *UserMappingsQuery) ([]UserMapping, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "userMappinglication/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var userMappings []UserMapping
	for _, bytes := range resp {
		var unmarshalled []UserMapping
		json.Unmarshal(bytes, &unmarshalled)
		if len(userMappings) == 0 {
			userMappings = unmarshalled
		} else {
			userMappings = append(userMappings, unmarshalled...)
		}
	}

	return userMappings, nil
}

// GetOne retrieves the user mapping by id, and if successful, it returns
// the http response and the pointer to the user mapping.
func (svc *V2Service) GetOne(id int32) (*UserMapping, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}

	var mapping UserMapping

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	json.Unmarshal(resp[0], &mapping)

	return &mapping, nil
}

// Update takes a user mapping and an id and attempts to use the parameters to update it
// in the API. Modifies the user mapping in place, or returns an error if one occurs
func (svc *V2Service) Update(mapping *UserMapping) error {
	if mapping.ID == nil {
		return errors.New("No ID Given")
	}
	validationErr := validateMappingValues(mapping, svc.LegalValuesService)
	if validationErr != nil {
		return validationErr
	}
	id := *mapping.ID
	mapping.ID = nil
	resp, err := svc.UpdateRaw(id, mapping)
	if err != nil {
		return err
	}

	var mappingID map[string]int
	json.Unmarshal(resp, &mappingID)

	mapping.ID = oltypes.Int32(int32(mappingID["id"]))

	return nil
}

// UpdateRaw takes a user mapping and an id and attempts to use the parameters to update it
// in the API. Returns the raw response bytes or an error.
func (svc *V2Service) UpdateRaw(id int32, mapping interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    mapping,
	})
}

// Create creates a new user mapping, and if successful, it returns
// the http response and the pointer to the user mapping.
func (svc *V2Service) Create(mapping *UserMapping) error {
	validationErr := validateMappingValues(mapping, svc.LegalValuesService)
	if validationErr != nil {
		return validationErr
	}
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    mapping,
	})

	if err != nil {
		return err
	}
	var mappingID map[string]int
	json.Unmarshal(resp, &mappingID)

	mapping.ID = oltypes.Int32(int32(mappingID["id"]))

	return nil
}

// Destroy deletes the user mapping for the id, and if successful, it returns nil
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

func validateMappingValues(mapping *UserMapping, svc services.SimpleQuery) error {
	legalValRequests := map[string][]string{}
	legalValRequests["mappings/conditions"] = []string{}
	legalValRequests["mappings/actions"] = []string{}
	for _, condition := range mapping.Conditions {
		legalValRequests[fmt.Sprintf("mappings/conditions/%s/values", *condition.Source)] = []string{}
		legalValRequests[fmt.Sprintf("mappings/conditions/%s/operators", *condition.Source)] = []string{}
	}
	for _, action := range mapping.Actions {
		legalValRequests[fmt.Sprintf("mappings/actions/%s/values", *action.Action)] = []string{}
	}

	var (
		wg    sync.WaitGroup
		mutex = &sync.Mutex{}
	)
	for reqURL := range legalValRequests {
		wg.Add(1)
		go func(reqURL string, legalValRequest map[string][]string) {
			defer wg.Done()
			legalValResp := []map[string]string{}
			err := svc.Query(reqURL, &legalValResp)
			if err != nil {
				log.Println("Problem validating mapping", reqURL, err)
			}
			legalVals := make([]string, len(legalValResp))
			for i, legalVal := range legalValResp {
				legalVals[i] = legalVal["value"]
			}
			mutex.Lock()
			legalValRequests[reqURL] = legalVals
			mutex.Unlock()
		}(reqURL, legalValRequests)
	}
	wg.Wait()

	errorMsgs := make([]error, 0)
	for _, condition := range mapping.Conditions {
		if len(legalValRequests["mappings/conditions"]) > 0 {
			err := utils.OneOf(fmt.Sprintf("%s.conditions.source", *mapping.Name), *condition.Source, legalValRequests["mappings/conditions"])
			if err != nil {
				log.Println("Illegal value given for condition source")
				errorMsgs = append(errorMsgs, err)
			}
		}
		if len(legalValRequests[fmt.Sprintf("mappings/conditions/%s/values", *condition.Source)]) > 0 {
			err := utils.OneOf(fmt.Sprintf("%s.conditions.value", *mapping.Name), *condition.Value, legalValRequests[fmt.Sprintf("mappings/conditions/%s/values", *condition.Source)])
			if err != nil {
				log.Println("Illegal value given for condition value")
				errorMsgs = append(errorMsgs, err)
			}
		}
		if len(legalValRequests[fmt.Sprintf("mappings/conditions/%s/operators", *condition.Source)]) > 0 {
			err := utils.OneOf(fmt.Sprintf("%s.conditions.operator", *mapping.Name), *condition.Operator, legalValRequests[fmt.Sprintf("mappings/conditions/%s/operators", *condition.Source)])
			if err != nil {
				log.Println("Illegal value given for condition operator")
				errorMsgs = append(errorMsgs, err)
			}
		}
	}

	for _, action := range mapping.Actions {
		if len(legalValRequests["mappings/actions"]) > 0 {
			err := utils.OneOf(fmt.Sprintf("%s.actions.action", *mapping.Name), *action.Action, legalValRequests["mappings/actions"])
			if err != nil {
				log.Println("Illegal value given for action")
				errorMsgs = append(errorMsgs, err)
			}
		}
		for _, val := range action.Value {
			if len(legalValRequests[fmt.Sprintf("mappings/actions/%s/values", *action.Action)]) > 0 {
				err := utils.OneOf(fmt.Sprintf("%s.actions.values", *mapping.Name), val, legalValRequests[fmt.Sprintf("mappings/actions/%s/values", *action.Action)])
				if err != nil {
					log.Println("Illegal value given for action value")
					errorMsgs = append(errorMsgs, err)
				}
			}
		}
	}
	return customerrors.StackErrors(errorMsgs)
}
