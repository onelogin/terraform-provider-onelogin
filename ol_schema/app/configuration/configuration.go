package appconfigurationschema

import (
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

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
		if n, err = strconv.Atoi(st); err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
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
		outOidc := models.ConfigurationOpenId{}

		// Set OIDC fields
		outOidc.RedirectURI = getString(s["redirect_uri"])
		outOidc.LoginURL = getString(s["login_url"])

		// Convert string to int for these fields
		if outOidc.RefreshTokenExpirationMinutes, err = getInt(s["refresh_token_expiration_minutes"]); err != nil {
			return nil, err
		}
		if outOidc.OidcApplicationType, err = getInt(s["oidc_application_type"]); err != nil {
			return nil, err
		}
		if outOidc.TokenEndpointAuthMethod, err = getInt(s["token_endpoint_auth_method"]); err != nil {
			return nil, err
		}
		if outOidc.AccessTokenExpirationMinutes, err = getInt(s["access_token_expiration_minutes"]); err != nil {
			return nil, err
		}

		return outOidc, nil
	} else if configType == "saml" {
		outSaml := models.ConfigurationSAML{}

		// Set SAML fields
		outSaml.SignatureAlgorithm = getString(s["signature_algorithm"])

		// Handle provider_arn which can be string or interface{}
		if s["provider_arn"] != nil {
			outSaml.ProviderArn = s["provider_arn"]
		}

		// Convert string to int for certificate_id
		certId, err := getInt(s["certificate_id"])
		if err != nil {
			return nil, err
		}
		outSaml.CertificateID = certId

		return outSaml, nil
	}

	// Return an empty map if we can't determine the type
	return map[string]interface{}{}, nil
}

// FlattenOIDC takes an instance of ConfigurationOpenId and returns a map of interface{}
func FlattenOIDC(config models.ConfigurationOpenId) map[string]interface{} {
	tfOut := map[string]interface{}{}

	// Add non-empty fields
	if config.RedirectURI != "" {
		tfOut["redirect_uri"] = config.RedirectURI
	}

	if config.LoginURL != "" {
		tfOut["login_url"] = config.LoginURL
	}

	// Terraform typeMap wants all strings so we convert int to string here
	if config.RefreshTokenExpirationMinutes != 0 {
		tfOut["refresh_token_expiration_minutes"] = strconv.FormatInt(int64(config.RefreshTokenExpirationMinutes), 10)
	}

	if config.OidcApplicationType != 0 {
		tfOut["oidc_application_type"] = strconv.FormatInt(int64(config.OidcApplicationType), 10)
	}

	if config.TokenEndpointAuthMethod != 0 {
		tfOut["token_endpoint_auth_method"] = strconv.FormatInt(int64(config.TokenEndpointAuthMethod), 10)
	}

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
	tfOut := map[string]interface{}{}

	// Determine if this is OIDC or SAML based on fields
	if _, ok := config["redirect_uri"]; ok {
		// Handle OIDC fields
		if val, ok := config["redirect_uri"].(string); ok && val != "" {
			tfOut["redirect_uri"] = val
		}

		if val, ok := config["login_url"].(string); ok && val != "" {
			tfOut["login_url"] = val
		}

		// Handle numeric fields, converting to string
		if val, ok := config["refresh_token_expiration_minutes"].(float64); ok && val != 0 {
			tfOut["refresh_token_expiration_minutes"] = strconv.FormatInt(int64(val), 10)
		}

		if val, ok := config["oidc_application_type"].(float64); ok && val != 0 {
			tfOut["oidc_application_type"] = strconv.FormatInt(int64(val), 10)
		}

		if val, ok := config["token_endpoint_auth_method"].(float64); ok && val != 0 {
			tfOut["token_endpoint_auth_method"] = strconv.FormatInt(int64(val), 10)
		}

		if val, ok := config["access_token_expiration_minutes"].(float64); ok && val != 0 {
			tfOut["access_token_expiration_minutes"] = strconv.FormatInt(int64(val), 10)
		}
	} else if _, ok := config["signature_algorithm"]; ok {
		// Handle SAML fields
		if val, ok := config["signature_algorithm"].(string); ok && val != "" {
			tfOut["signature_algorithm"] = val
		}

		if val, ok := config["provider_arn"]; ok && val != nil {
			tfOut["provider_arn"] = val
		}

		if val, ok := config["certificate_id"].(float64); ok && val != 0 {
			tfOut["certificate_id"] = int(val)
		}
	}

	return tfOut
}
