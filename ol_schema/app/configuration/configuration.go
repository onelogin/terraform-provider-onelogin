package appconfigurationschema

import (
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/terraform-provider-onelogin/utils"
	"strconv"
)

func validSignatureAlgorithm(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"})
}

// Inflate takes a map of interfaces and uses the fields to construct
// an AppConfiguration instance.
func Inflate(s map[string]interface{}) (apps.AppConfiguration, error) {
	out := apps.AppConfiguration{}
	var (
		st     string
		n      int
		notNil bool
		err    error
	)

	// oidc fields
	if st, notNil = s["redirect_uri"].(string); notNil {
		out.RedirectURI = oltypes.String(st)
	}
	if st, notNil = s["login_url"].(string); notNil {
		out.LoginURL = oltypes.String(st)
	}
	// terraform typeMap wants all fields to be strings and we store these fields as in32
	// so we do the conversion here when assembling the resource
	if st, notNil = s["refresh_token_expiration_minutes"].(string); notNil {
		if n, err = strconv.Atoi(st); err != nil {
			return out, err
		}
		out.RefreshTokenExpirationMinutes = oltypes.Int32(int32(n))
	}

	if st, notNil = s["oidc_application_type"].(string); notNil {
		if n, err = strconv.Atoi(st); err != nil {
			return out, err
		}
		out.OidcApplicationType = oltypes.Int32(int32(n))
	}

	if st, notNil = s["token_endpoint_auth_method"].(string); notNil {
		if n, err = strconv.Atoi(st); err != nil {
			return out, err
		}
		out.TokenEndpointAuthMethod = oltypes.Int32(int32(n))
	}

	if st, notNil = s["access_token_expiration_minutes"].(string); notNil {
		if n, err = strconv.Atoi(st); err != nil {
			return out, err
		}
		out.AccessTokenExpirationMinutes = oltypes.Int32(int32(n))
	}

	// saml fields
	if st, notNil = s["provider_arn"].(string); notNil {
		out.ProviderArn = oltypes.String(st)
	}
	if st, notNil = s["signature_algorithm"].(string); notNil {
		out.SignatureAlgorithm = oltypes.String(st)
	}
	return out, nil
}

// FlattenOIDC takes an instance of AppConfiguration and return an array of
// maps. Fields differ depending on if the app is a SAML or OIDC app.
func FlattenOIDC(config apps.AppConfiguration) map[string]interface{} {
	tfOut := map[string]interface{}{}
	if config.RedirectURI != nil {
		tfOut["redirect_uri"] = *config.RedirectURI
	}
	if config.LoginURL != nil {
		tfOut["login_url"] = *config.LoginURL
	}
	// Terraform typeMap wants all strings so we convert int32 to string here
	if config.RefreshTokenExpirationMinutes != nil {
		tfOut["refresh_token_expiration_minutes"] = strconv.FormatInt(int64(*config.RefreshTokenExpirationMinutes), 10)
	}
	if config.OidcApplicationType != nil {
		tfOut["oidc_application_type"] = strconv.FormatInt(int64(*config.OidcApplicationType), 10)
	}
	if config.TokenEndpointAuthMethod != nil {
		tfOut["token_endpoint_auth_method"] = strconv.FormatInt(int64(*config.TokenEndpointAuthMethod), 10)
	}
	if config.AccessTokenExpirationMinutes != nil {
		tfOut["access_token_expiration_minutes"] = strconv.FormatInt(int64(*config.AccessTokenExpirationMinutes), 10)
	}
	return tfOut
}

// FlattenSAML takes an instance of AppConfiguration and return an array of
// maps. Fields differ depending on if the app is a SAML or OIDC app.
func FlattenSAML(config apps.AppConfiguration) map[string]interface{} {
	tfOut := map[string]interface{}{}
	if config.ProviderArn != nil {
		tfOut["provider_arn"] = *config.ProviderArn
	}
	if config.SignatureAlgorithm != nil {
		tfOut["signature_algorithm"] = *config.SignatureAlgorithm
	}
	return tfOut
}
