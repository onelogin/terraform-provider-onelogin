package appconfigurationschema

import (
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

func validSignatureAlgorithm(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"})
}

func getOlString(v interface{}) *string {
	if st, notNil := v.(string); notNil {
		return oltypes.String(st)
	}
	return nil
}

func getOlInt32(v interface{}) (*int32, error) {
	var (
		n   int
		err error
	)
	if st, notNil := v.(string); notNil {
		if n, err = strconv.Atoi(st); err != nil {
			return nil, err
		}
		return oltypes.Int32(int32(n)), nil
	}
	return nil, nil
}

func getOlIntAsString(v interface{}) (*string, error) {
	var err error
	if st, notNil := v.(string); notNil {
		if _, err = strconv.Atoi(st); err != nil {
			return nil, err
		}
		return oltypes.String(st), nil
	}
	return nil, nil
}

// Inflate takes a map of interfaces and uses the fields to construct
// an AppConfiguration instance.
func Inflate(s map[string]interface{}) (apps.AppConfiguration, error) {
	out := apps.AppConfiguration{}
	var err error

	// oidc fields
	out.RedirectURI = getOlString(s["redirect_uri"])
	out.PostLogoutRedirectURI = getOlString(s["post_logout_redirect_uri"])
	out.LoginURL = getOlString(s["login_url"])
	out.ProviderArn = getOlString(s["provider_arn"])
	out.IdpList = getOlString(s["idp_list"])
	out.SignatureAlgorithm = getOlString(s["signature_algorithm"])
	out.LogoutURL = getOlString(s["logout_url"])
	out.Audience = getOlString(s["audience"])
	out.ConsumerURL = getOlString(s["consumer_url"])
	out.Login = getOlString(s["login"])
	out.Recipient = getOlString(s["recipient"])
	out.Validator = getOlString(s["validator"])
	out.RelayState = getOlString(s["relaystate"])
	out.Relay = getOlString(s["relay"])

	// terraform typeMap wants all fields to be strings and we store these fields as int32
	// so we do the conversion here when assembling the resource
	if out.RefreshTokenExpirationMinutes, err = getOlInt32(s["refresh_token_expiration_minutes"]); err != nil {
		return out, err
	}
	if out.OidcApplicationType, err = getOlInt32(s["oidc_application_type"]); err != nil {
		return out, err
	}
	if out.TokenEndpointAuthMethod, err = getOlInt32(s["token_endpoint_auth_method"]); err != nil {
		return out, err
	}
	if out.AccessTokenExpirationMinutes, err = getOlInt32(s["access_token_expiration_minutes"]); err != nil {
		return out, err
	}
	if out.SAMLNotValidOnOrAafter, err = getOlIntAsString(s["saml_notonorafter"]); err != nil {
		return out, err
	}
	if out.GenerateAttributeValueTags, err = getOlIntAsString(s["generate_attribute_value_tags"]); err != nil {
		return out, err
	}
	if out.SAMLInitiaterID, err = getOlIntAsString(s["saml_initiater_id"]); err != nil {
		return out, err
	}
	if out.SAMLNotValidBefore, err = getOlIntAsString(s["saml_notbefore"]); err != nil {
		return out, err
	}
	if out.SAMLIssuerType, err = getOlIntAsString(s["saml_issuer_type"]); err != nil {
		return out, err
	}
	if out.SAMLSignElement, err = getOlIntAsString(s["saml_sign_element"]); err != nil {
		return out, err
	}
	if out.EncryptAssertion, err = getOlIntAsString(s["encrypt_assertion"]); err != nil {
		return out, err
	}
	if out.SAMLSessionNotValidOnOrAfter, err = getOlIntAsString(s["saml_sessionnotonorafter"]); err != nil {
		return out, err
	}
	if out.SAMLEncryptionMethodID, err = getOlIntAsString(s["saml_encryption_method_id"]); err != nil {
		return out, err
	}
	if out.SAMLNameIDFormatID, err = getOlIntAsString(s["saml_nameid_format_id"]); err != nil {
		return out, err
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
	if config.PostLogoutRedirectURI != nil {
		tfOut["post_logout_redirect_uri"] = *config.PostLogoutRedirectURI
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
	if config.IdpList != nil {
		tfOut["idp_list"] = *config.IdpList
	}

	if config.SignatureAlgorithm != nil {
		tfOut["signature_algorithm"] = *config.SignatureAlgorithm
	}
	return tfOut
}
