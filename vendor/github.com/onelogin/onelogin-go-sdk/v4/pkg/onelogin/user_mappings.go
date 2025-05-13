package onelogin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	UserMappingsPath string = "api/2/mappings"
)

// ListUserMappings gets a list of all User Mappings
// Returns an array of UserMapping objects or an error
func (sdk *OneloginSDK) ListUserMappings(query *mod.UserMappingsQuery) ([]mod.UserMapping, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	var mappings []mod.UserMapping
	err = utl.CheckHTTPResponseAndUnmarshal(resp, &mappings)
	return mappings, err
}

// GetUserMapping gets a specific User Mapping by ID
// Returns a UserMapping object or an error
func (sdk *OneloginSDK) GetUserMapping(mappingID int32) (*mod.UserMapping, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, mappingID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	var mapping mod.UserMapping
	err = utl.CheckHTTPResponseAndUnmarshal(resp, &mapping)
	return &mapping, err
}

// CreateUserMapping creates a new User Mapping
// Returns the created UserMapping object or an error
func (sdk *OneloginSDK) CreateUserMapping(mapping mod.UserMapping) (*mod.UserMapping, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, mapping)
	if err != nil {
		return nil, err
	}
	var newMapping mod.UserMapping
	err = utl.CheckHTTPResponseAndUnmarshal(resp, &newMapping)
	return &newMapping, err
}

// UpdateUserMapping updates an existing User Mapping
// Returns the updated UserMapping object or an error
func (sdk *OneloginSDK) UpdateUserMapping(mappingID int32, mapping mod.UserMapping) (*mod.UserMapping, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, mappingID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, mapping)
	if err != nil {
		return nil, err
	}

	// Some endpoints return just {id: XYZ} instead of the full object
	// Check for this case and handle it by fetching the full object
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted {
		// Try to unmarshal the response
		var responseObj struct {
			ID int32 `json:"id"`
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		// Check if the response is just an ID
		err = json.Unmarshal(body, &responseObj)
		if err == nil && responseObj.ID > 0 {
			// If it's just an ID, retrieve the full object
			return sdk.GetUserMapping(responseObj.ID)
		}

		// Otherwise, create a new response with the original body for the usual unmarshaling
		resp = &http.Response{
			StatusCode: resp.StatusCode,
			Body:       io.NopCloser(bytes.NewBuffer(body)),
		}
	}

	var updatedMapping mod.UserMapping
	err = utl.CheckHTTPResponseAndUnmarshal(resp, &updatedMapping)
	return &updatedMapping, err
}

// DeleteUserMapping deletes a User Mapping by ID
// Returns nil on success or an error
func (sdk *OneloginSDK) DeleteUserMapping(mappingID int32) error {
	p, err := utl.BuildAPIPath(UserMappingsPath, mappingID)
	if err != nil {
		return err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return err
	}
	_, err = utl.CheckHTTPResponse(resp)
	return err
}

// DryRunUserMapping performs a dry run of a User Mapping against specific user IDs
// Returns the results of the dry run or an error
func (sdk *OneloginSDK) DryRunUserMapping(mappingID int32, userIDs []int32) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, mappingID, "dryrun")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, userIDs)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// ListUserMappingConditions gets a list of all available conditions for User Mappings
// Returns the list of conditions or an error
func (sdk *OneloginSDK) ListUserMappingConditions() (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, "conditions")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// ListUserMappingConditionOperators gets a list of operators available for a condition
// Returns the list of operators or an error
func (sdk *OneloginSDK) ListUserMappingConditionOperators(conditionValue string) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, "conditions", conditionValue, "operators")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// ListUserMappingConditionValues gets a list of values available for a condition
// Returns the list of values or an error
func (sdk *OneloginSDK) ListUserMappingConditionValues(conditionValue string) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, "conditions", conditionValue, "values")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// ListUserMappingActions gets a list of all available actions for User Mappings
// Returns the list of actions or an error
func (sdk *OneloginSDK) ListUserMappingActions() (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, "actions")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// ListUserMappingActionValues gets a list of values available for an action
// Returns the list of values or an error
func (sdk *OneloginSDK) ListUserMappingActionValues(actionValue string) (interface{}, error) {
	p, err := utl.BuildAPIPath(UserMappingsPath, "actions", actionValue, "values")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// BulkSortUserMappings updates the order of multiple User Mappings at once
// Returns nil on success or an error
func (sdk *OneloginSDK) BulkSortUserMappings(mappingIDs []int32) error {
	p, err := utl.BuildAPIPath(UserMappingsPath, "sort")
	if err != nil {
		return err
	}
	resp, err := sdk.Client.Put(&p, mappingIDs)
	if err != nil {
		return err
	}
	_, err = utl.CheckHTTPResponse(resp)
	return err
}

// Legacy functions kept for backwards compatibility
func (sdk *OneloginSDK) ListMappings() (interface{}, error) {
	mappings, err := sdk.ListUserMappings(nil)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

func (sdk *OneloginSDK) GetMapping(mappingID int) (interface{}, error) {
	// Check for potential integer overflow when converting to int32
	if mappingID > 2147483647 || mappingID < -2147483648 {
		return nil, fmt.Errorf("mapping ID %d is outside the range of int32", mappingID)
	}
	return sdk.GetUserMapping(int32(mappingID))
}

func (sdk *OneloginSDK) CreateMapping(mapping mod.UserMapping) (interface{}, error) {
	return sdk.CreateUserMapping(mapping)
}

func (sdk *OneloginSDK) UpdateMapping(mappingID int, mapping mod.UserMapping) (interface{}, error) {
	// Check for potential integer overflow when converting to int32
	if mappingID > 2147483647 || mappingID < -2147483648 {
		return nil, fmt.Errorf("mapping ID %d is outside the range of int32", mappingID)
	}
	return sdk.UpdateUserMapping(int32(mappingID), mapping)
}

func (sdk *OneloginSDK) DeleteMapping(mappingID int) (interface{}, error) {
	// Check for potential integer overflow when converting to int32
	if mappingID > 2147483647 || mappingID < -2147483648 {
		return nil, fmt.Errorf("mapping ID %d is outside the range of int32", mappingID)
	}
	err := sdk.DeleteUserMapping(int32(mappingID))
	return nil, err
}

func (sdk *OneloginSDK) DryrunMapping(mappingID int, userIds []int) (interface{}, error) {
	// Check for potential integer overflow when converting to int32
	if mappingID > 2147483647 || mappingID < -2147483648 {
		return nil, fmt.Errorf("mapping ID %d is outside the range of int32", mappingID)
	}
	
	// Convert int slice to int32 slice with overflow checking
	userIDs32 := make([]int32, len(userIds))
	for i, id := range userIds {
		if id > 2147483647 || id < -2147483648 {
			return nil, fmt.Errorf("user ID %d is outside the range of int32", id)
		}
		userIDs32[i] = int32(id)
	}
	
	// Safe to convert since we've already checked the range above
	// #nosec G115 - conversion is safe because we checked the range
	mappingID32 := int32(mappingID)
	return sdk.DryRunUserMapping(mappingID32, userIDs32)
}

func (sdk *OneloginSDK) ListConditions() (interface{}, error) {
	return sdk.ListUserMappingConditions()
}

func (sdk *OneloginSDK) ListConditionOperators(conditionValue string) (interface{}, error) {
	return sdk.ListUserMappingConditionOperators(conditionValue)
}

func (sdk *OneloginSDK) ListConditionValues(conditionValue string) (interface{}, error) {
	return sdk.ListUserMappingConditionValues(conditionValue)
}

func (sdk *OneloginSDK) ListActions() (interface{}, error) {
	return sdk.ListUserMappingActions()
}

func (sdk *OneloginSDK) ListActionValues(actionValue string) (interface{}, error) {
	return sdk.ListUserMappingActionValues(actionValue)
}

func (sdk *OneloginSDK) BulkSortMappings(mappingIDs []int) (interface{}, error) {
	// Convert int slice to int32 slice with overflow checking
	mappingIDs32 := make([]int32, len(mappingIDs))
	for i, id := range mappingIDs {
		if id > 2147483647 || id < -2147483648 {
			return nil, fmt.Errorf("mapping ID %d is outside the range of int32", id)
		}
		mappingIDs32[i] = int32(id)
	}
	err := sdk.BulkSortUserMappings(mappingIDs32)
	return nil, err
}
