package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	MFAPath string = "api/2/mfa/users"
)

// https://<subdomain>/api/2/mfa/users/<user_id>/factors to return a list of authentication factors that are available for user enrollment via API
func (sdk *OneloginSDK) GetAvailableMFAFactors(userID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "factors")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/registrations to initiate enrollment for user with a given authentication factor
func (sdk *OneloginSDK) EnrollMFAFactor(factor models.EnrollFactorRequest, userID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "registrations")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, factor)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/devices to return a list of authentication factors registered to a particular user for MFA
func (sdk *OneloginSDK) GetEnrolledFactor(userID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "devices")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/devices/<device_id> to remove an enrolled factor from a user
func (sdk *OneloginSDK) RemoveMFAFactor(userID, deviceID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "devices", deviceID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/verifications to trigger an SMS, Voice, Email or Push notification containing a One-Time Password
// or Magic Link that can be used to authenticate a user
func (sdk *OneloginSDK) ActivateMFAFactor(userID int, request models.ActivateFactorRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "verifications")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, request)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/:user_id/mfa_token to generate a temporary MFA token that can be used in place of other MFA tokens
// for a set time period. For example, use this token for account recovery when an MFA device has been lost
func (sdk *OneloginSDK) GenerateMFAToken(userID int, request models.GenerateMFATokenRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "mfa_token")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, request)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/registrations/<registration_id>  to verify enrollment for OneLogin SMS, OneLogin Email,
// OneLogin Protect and Authenticator authentication factors
func (sdk *OneloginSDK) VerifyMFAEnrollment(userID int, registrationID string, request models.VerifyEnrollmentFactorRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "registrations", registrationID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, request)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/registrations/<registration_id> to verify enrollment for OneLogin Voice
func (sdk *OneloginSDK) VerifyMFAEnrollmentGet(userID int, registrationID string) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "registrations", registrationID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/verifications/<verification_id> to verify an OTP code provided by SMS, Email, or Authenticator
func (sdk *OneloginSDK) VerifyAuthFactor(userID, verificationID int, request models.VerificationFactorRequest) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "verifications", verificationID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, request)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// https://<subdomain>/api/2/mfa/users/<user_id>/verifications/<verification_id> to verify completion of OneLogin Push or OneLogin Voice factors,
// or in cases where email is used as an authentication factor via Magic Link rather than OTP
func (sdk *OneloginSDK) VerifyAuthFactorGet(userID, verificationID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(MFAPath, userID, "verifications", verificationID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, nil)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
