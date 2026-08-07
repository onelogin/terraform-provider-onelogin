package appconfigurationschema

import (
	"fmt"
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// CustomConfigurationOpenId is a wrapper around ConfigurationOpenId that allows
// omitting timeout fields when they are not explicitly set, to avoid overriding
// API defaults with 0 values during updates.
type CustomConfigurationOpenId struct {
	RedirectURI                   string  `json:"redirect_uri,omitempty"`
	PostLogoutRedirectURI         *string `json:"post_logout_redirect_uri,omitempty"`
	IncludeAmrClaims              *bool   `json:"include_amr_claims,omitempty"`
	LoginURL                      string  `json:"login_url,omitempty"`
	OidcApplicationType           int     `json:"oidc_application_type,omitempty"`
	TokenEndpointAuthMethod       int     `json:"token_endpoint_auth_method,omitempty"`
	AccessTokenExpirationMinutes  *int    `json:"access_token_expiration_minutes,omitempty"`
	RefreshTokenExpirationMinutes *int    `json:"refresh_token_expiration_minutes,omitempty"`
}

func validSignatureAlgorithm(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"})
}

func getString(v interface{}) string {
	if st, notNil := v.(string); notNil {
		return st
	}
	return ""
}

func getInt(v interface{}) (int, error) {
	var (
		n   int
		err error
	)
	if st, notNil := v.(string); notNil {
		// Handle empty string as unset (return 0 without error)
		if st == "" {
			return 0, nil
		}
		if n, err = strconv.Atoi(st); err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
}

// intPtr creates a pointer to an int value for cleaner syntax when creating pointers
func intPtr(val int) *int {
	return &val
}

// handleTimeoutField processes a timeout field value and returns a pointer if the value is valid
func handleTimeoutField(s map[string]interface{}, fieldName string) (*int, error) {
	if val, exists := s[fieldName]; exists {
		if strVal, ok := val.(string); ok && strVal != "" {
			if timeoutVal, err := getInt(val); err != nil {
				return nil, err
			} else {
				return intPtr(timeoutVal), nil
			}
		}
		// If empty string or not provided, return nil (will be omitted from JSON)
	}
	return nil, nil
}

// Inflate takes a map of interfaces and uses the fields to construct
// either a ConfigurationOpenId or ConfigurationSAML instance.
func Inflate(s map[string]interface{}) (interface{}, error) {
	var err error
	var configType string

	// Determine if this is OpenID or SAML based on fields
	if _, ok := s["redirect_uri"]; ok {
		configType = "openid"
	} else if _, ok := s["signature_algorithm"]; ok {
		configType = "saml"
	}

	if configType == "openid" {
		customOidc := CustomConfigurationOpenId{}

		// Set OIDC fields
		customOidc.RedirectURI = getString(s["redirect_uri"])
		customOidc.LoginURL = getString(s["login_url"])

		// Present-but-empty means "remove every URI", which is a different
		// request from leaving the attribute out and is why the field is a
		// pointer. Only an absent key leaves the app's URIs alone.
		if val, exists := s["post_logout_redirect_uri"]; exists {
			uris := getString(val)
			customOidc.PostLogoutRedirectURI = &uris
		}

		// A pointer for the same reason: only an absent key leaves the app's
		// setting alone, so a configuration that turns the claim off sends
		// false rather than being indistinguishable from not mentioning it.
		//
		// The configuration block is a map of strings, so this arrives as
		// "true" or "false". Anything else is rejected rather than ignored:
		// dropping it would leave the key in the configuration and out of
		// state, which is a diff every plan reports and no apply can settle.
		// A typo should say so once, not forever.
		if val, exists := s["include_amr_claims"]; exists {
			parsed, parseErr := strconv.ParseBool(getString(val))
			if parseErr != nil {
				return nil, fmt.Errorf("include_amr_claims must be true or false, got %q", getString(val))
			}
			customOidc.IncludeAmrClaims = &parsed
		}

		// Handle timeout fields specially - only set them if explicitly provided and non-empty
		// This prevents overriding existing API values with 0 when fields are not specified
		if customOidc.RefreshTokenExpirationMinutes, err = handleTimeoutField(s, "refresh_token_expiration_minutes"); err != nil {
			return nil, err
		}

		if customOidc.AccessTokenExpirationMinutes, err = handleTimeoutField(s, "access_token_expiration_minutes"); err != nil {
			return nil, err
		}

		// Convert string to int for these required fields
		if customOidc.OidcApplicationType, err = getInt(s["oidc_application_type"]); err != nil {
			return nil, err
		}
		if customOidc.TokenEndpointAuthMethod, err = getInt(s["token_endpoint_auth_method"]); err != nil {
			return nil, err
		}

		return customOidc, nil
	} else if configType == "saml" {
		// Instead of using the limited ConfigurationSAML struct, create a generic map
		// that passes through all SAML configuration fields provided by the user.
		// The OneLogin API supports many more SAML fields than the SDK struct defines.
		outSaml := make(map[string]interface{})

		// Copy all provided fields to the output map
		for key, value := range s {
			if value != nil && value != "" {
				// Handle special field type conversions
				switch key {
				case "certificate_id":
					// Convert certificate_id to int as expected by API
					if certId, err := getInt(value); err != nil {
						return nil, err
					} else if certId != 0 {
						outSaml[key] = certId
					}
				default:
					// Pass through all other fields as-is
					outSaml[key] = value
				}
			}
		}

		return outSaml, nil
	}

	// Return an empty map if we can't determine the type
	return map[string]interface{}{}, nil
}

// RetainManaged narrows a flattened configuration to the keys the practitioner
// is actually managing.
//
// configuration is a TypeMap, and a map has no way to mark an individual key
// Computed. A key that is in state and not in the configuration is a removal
// every plan proposes and no apply can settle, because the next read puts it
// straight back. The API returns its own defaults regardless of what was sent:
// an OIDC app always comes back with include_amr_claims false and
// oidc_application_type 0, and a SAML app comes back with a few dozen keys, so
// recording everything it sends leaves anyone who did not spell out the whole
// object with a permanent diff.
//
// The prior state's key set is the practitioner's. On create it is the planned
// configuration, on update it is the new one, and on refresh it is whatever the
// last apply wrote. A key outside it is not managed by this resource, so its
// value is not recorded and drift in it is not drift in anything the
// configuration declares. An empty prior state is an import, where there is no
// configuration to respect yet and the whole object is written.
func RetainManaged(prior interface{}, flattened map[string]interface{}) map[string]interface{} {
	managed, ok := prior.(map[string]interface{})
	if !ok || len(managed) == 0 {
		return flattened
	}

	out := make(map[string]interface{}, len(managed))
	for key := range managed {
		// Absent from the response rather than nil: a key the API does not
		// report back is one this provider cannot claim a value for.
		if val, found := flattened[key]; found {
			out[key] = val
		}
	}
	return out
}

// FlattenOIDC takes an instance of ConfigurationOpenId and returns a map of interface{}
func FlattenOIDC(config models.ConfigurationOpenId) map[string]interface{} {
	tfOut := map[string]interface{}{}

	// Add non-empty fields
	if config.RedirectURI != "" {
		tfOut["redirect_uri"] = config.RedirectURI
	}

	// Keeps an empty value, matching Flatten below. The pointer draws the same
	// line the type assertion draws there: nil is an app that never had URIs,
	// while a pointer to "" is one whose URIs were cleared, and state has to
	// keep the latter or the plan never goes quiet.
	if config.PostLogoutRedirectURI != nil {
		tfOut["post_logout_redirect_uri"] = *config.PostLogoutRedirectURI
	}

	// The configuration block is a map of strings, so the bool is rendered as
	// one. Absent stays absent: a nil pointer is an app the API said nothing
	// about, which is not the same as one that has the claim switched off.
	if config.IncludeAmrClaims != nil {
		tfOut["include_amr_claims"] = strconv.FormatBool(*config.IncludeAmrClaims)
	}

	if config.LoginURL != "" {
		tfOut["login_url"] = config.LoginURL
	}

	// Terraform typeMap wants all strings so we convert int to string here
	if config.RefreshTokenExpirationMinutes != 0 {
		tfOut["refresh_token_expiration_minutes"] = strconv.FormatInt(int64(config.RefreshTokenExpirationMinutes), 10)
	}

	// Always written, unlike the timeouts below. Zero is a real value for both
	// of these -- application type 0 is Web and auth method 0 is basic -- so
	// skipping it kept the most common configuration out of state entirely,
	// and the plan never settled.
	tfOut["oidc_application_type"] = strconv.FormatInt(int64(config.OidcApplicationType), 10)
	tfOut["token_endpoint_auth_method"] = strconv.FormatInt(int64(config.TokenEndpointAuthMethod), 10)

	if config.AccessTokenExpirationMinutes != 0 {
		tfOut["access_token_expiration_minutes"] = strconv.FormatInt(int64(config.AccessTokenExpirationMinutes), 10)
	}

	return tfOut
}

// FlattenSAML takes an instance of ConfigurationSAML and returns a map of interface{}
func FlattenSAML(config models.ConfigurationSAML) map[string]interface{} {
	tfOut := map[string]interface{}{}

	// Add provider_arn if it exists
	if config.ProviderArn != nil {
		tfOut["provider_arn"] = config.ProviderArn
	}

	// Add other SAML fields
	if config.SignatureAlgorithm != "" {
		tfOut["signature_algorithm"] = config.SignatureAlgorithm
	}

	if config.CertificateID != 0 {
		tfOut["certificate_id"] = config.CertificateID
	}

	return tfOut
}

// Flatten takes a generic configuration map and returns a map of interface{}
func Flatten(config map[string]interface{}) map[string]interface{} {
	// If config is empty, return an empty map to ensure consistency
	if len(config) == 0 {
		return map[string]interface{}{}
	}

	tfOut := map[string]interface{}{}

	// Determine if this is OIDC or SAML based on fields
	if _, ok := config["redirect_uri"]; ok {
		// Handle OIDC fields
		if val, ok := config["redirect_uri"].(string); ok && val != "" {
			tfOut["redirect_uri"] = val
		}

		// Deliberately keeps an empty value, unlike the fields around it. The
		// API returns null for an app that never had URIs and "" for one whose
		// URIs were cleared, and the type assertion already separates the two:
		// null fails it. Dropping "" as well would leave state without a key
		// the practitioner's configuration still has, and the plan would never
		// go quiet.
		if val, ok := config["post_logout_redirect_uri"].(string); ok {
			tfOut["post_logout_redirect_uri"] = val
		}

		// JSON gives a bool here, while the configuration block is a map of
		// strings, so it is rendered as one. A null fails the assertion and is
		// left out, which is the distinction the pointer keeps everywhere else.
		if val, ok := config["include_amr_claims"].(bool); ok {
			tfOut["include_amr_claims"] = strconv.FormatBool(val)
		}

		if val, ok := config["login_url"].(string); ok && val != "" {
			tfOut["login_url"] = val
		}

		// Handle numeric fields, converting to string
		if val, ok := config["refresh_token_expiration_minutes"].(float64); ok && val != 0 {
			tfOut["refresh_token_expiration_minutes"] = strconv.FormatInt(int64(val), 10)
		}

		// Zero is a real value for both of these -- application type 0 is Web,
		// auth method 0 is basic -- and the API omits them at zero rather than
		// returning it. Defaulting to "0" is what an OIDC app without them
		// actually is, and skipping them kept the commonest configuration out
		// of state entirely, so the plan never settled.
		tfOut["oidc_application_type"] = "0"
		if val, ok := config["oidc_application_type"].(float64); ok {
			tfOut["oidc_application_type"] = strconv.FormatInt(int64(val), 10)
		}

		tfOut["token_endpoint_auth_method"] = "0"
		if val, ok := config["token_endpoint_auth_method"].(float64); ok {
			tfOut["token_endpoint_auth_method"] = strconv.FormatInt(int64(val), 10)
		}

		if val, ok := config["access_token_expiration_minutes"].(float64); ok && val != 0 {
			tfOut["access_token_expiration_minutes"] = strconv.FormatInt(int64(val), 10)
		}
	} else if _, ok := config["signature_algorithm"]; ok {
		// Handle SAML fields - pass through all fields from the API response
		// The OneLogin API supports many more SAML fields than just the basic ones
		for key, value := range config {
			if value != nil {
				switch key {
				case "certificate_id":
					// Handle certificate_id which comes back as float64 from JSON
					if val, ok := value.(float64); ok && val != 0 {
						tfOut[key] = strconv.FormatInt(int64(val), 10)
					}
				default:
					// Pass through all other fields, converting to string if needed for Terraform
					if strVal, ok := value.(string); ok && strVal != "" {
						tfOut[key] = strVal
					} else if numVal, ok := value.(float64); ok {
						// Always include numeric values, including zero, as they may be valid SAML configuration
						tfOut[key] = strconv.FormatInt(int64(numVal), 10)
					} else if boolVal, ok := value.(bool); ok {
						// Always include boolean values, consistent with numeric handling
						if boolVal {
							tfOut[key] = "1"
						} else {
							tfOut[key] = "0"
						}
					} else if value != nil {
						// For other types (like interface{}), convert to string
						tfOut[key] = fmt.Sprintf("%v", value)
					}
				}
			}
		}
	}

	// Ensure we always return a non-nil map
	return tfOut
}
