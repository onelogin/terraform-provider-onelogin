package onelogin

import (
	"context"
	
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	SelfRegistrationProfilesPath string = "api/2/self_registration_profiles"
)

// GetSelfRegistrationProfiles retrieves all self-registration profiles
func (sdk *OneloginSDK) GetSelfRegistrationProfiles(query mod.Queryable) (interface{}, error) {
	return sdk.GetSelfRegistrationProfilesWithContext(context.Background(), query)
}

// GetSelfRegistrationProfilesWithContext retrieves all self-registration profiles using the provided context
func (sdk *OneloginSDK) GetSelfRegistrationProfilesWithContext(ctx context.Context, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.GetWithContext(ctx, &p, query)
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// GetSelfRegistrationProfile retrieves a specific self-registration profile by ID
func (sdk *OneloginSDK) GetSelfRegistrationProfile(id int) (interface{}, error) {
	return sdk.GetSelfRegistrationProfileWithContext(context.Background(), id)
}

// GetSelfRegistrationProfileWithContext retrieves a specific self-registration profile by ID using the provided context
func (sdk *OneloginSDK) GetSelfRegistrationProfileWithContext(ctx context.Context, id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath, id)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.GetWithContext(ctx, &p, nil)
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// CreateSelfRegistrationProfile creates a new self-registration profile
func (sdk *OneloginSDK) CreateSelfRegistrationProfile(profile mod.SelfRegistrationProfile) (interface{}, error) {
	return sdk.CreateSelfRegistrationProfileWithContext(context.Background(), profile)
}

// CreateSelfRegistrationProfileWithContext creates a new self-registration profile using the provided context
func (sdk *OneloginSDK) CreateSelfRegistrationProfileWithContext(ctx context.Context, profile mod.SelfRegistrationProfile) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.PostWithContext(ctx, &p, map[string]mod.SelfRegistrationProfile{
		"self_registration_profile": profile,
	})
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// UpdateSelfRegistrationProfile updates an existing self-registration profile
func (sdk *OneloginSDK) UpdateSelfRegistrationProfile(id int, profile mod.SelfRegistrationProfile) (interface{}, error) {
	return sdk.UpdateSelfRegistrationProfileWithContext(context.Background(), id, profile)
}

// UpdateSelfRegistrationProfileWithContext updates an existing self-registration profile using the provided context
func (sdk *OneloginSDK) UpdateSelfRegistrationProfileWithContext(ctx context.Context, id int, profile mod.SelfRegistrationProfile) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath, id)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.PutWithContext(ctx, &p, map[string]mod.SelfRegistrationProfile{
		"self_registration_profile": profile,
	})
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// DeleteSelfRegistrationProfile deletes a self-registration profile
func (sdk *OneloginSDK) DeleteSelfRegistrationProfile(id int) (interface{}, error) {
	return sdk.DeleteSelfRegistrationProfileWithContext(context.Background(), id)
}

// DeleteSelfRegistrationProfileWithContext deletes a self-registration profile using the provided context
func (sdk *OneloginSDK) DeleteSelfRegistrationProfileWithContext(ctx context.Context, id int) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath, id)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.DeleteWithContext(ctx, &p)
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// CreateSelfRegistrationProfileField creates a new field for a self-registration profile
func (sdk *OneloginSDK) CreateSelfRegistrationProfileField(profileID int, customAttributeID int) (interface{}, error) {
	return sdk.CreateSelfRegistrationProfileFieldWithContext(context.Background(), profileID, customAttributeID)
}

// CreateSelfRegistrationProfileFieldWithContext creates a new field for a self-registration profile using the provided context
func (sdk *OneloginSDK) CreateSelfRegistrationProfileFieldWithContext(ctx context.Context, profileID int, customAttributeID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath, profileID, "self_registration_profile_fields")
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.PostWithContext(ctx, &p, map[string]int{
		"custom_attribute_id": customAttributeID,
	})
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}

// DeleteSelfRegistrationProfileField deletes a field from a self-registration profile
func (sdk *OneloginSDK) DeleteSelfRegistrationProfileField(profileID int, fieldID int) (interface{}, error) {
	return sdk.DeleteSelfRegistrationProfileFieldWithContext(context.Background(), profileID, fieldID)
}

// DeleteSelfRegistrationProfileFieldWithContext deletes a field from a self-registration profile using the provided context
func (sdk *OneloginSDK) DeleteSelfRegistrationProfileFieldWithContext(ctx context.Context, profileID int, fieldID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(SelfRegistrationProfilesPath, profileID, "self_registration_profile_fields", fieldID)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.Client.DeleteWithContext(ctx, &p)
	if err != nil {
		return nil, err
	}
	
	return utl.CheckHTTPResponse(resp)
}
