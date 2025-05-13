package appssoschema

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// FlattenOIDC takes a SSOOpenId instance and creates a map
func FlattenOIDC(sso models.SSOOpenId) map[string]interface{} {
	return map[string]interface{}{
		"client_id": sso.ClientID,
	}
}

// FlattenSAMLCert takes a SSOSAML instance and uses the Certificate node to create the map
func FlattenSAMLCert(sso models.SSOSAML) map[string]interface{} {
	return map[string]interface{}{
		"name":  sso.Certificate.Name,
		"value": sso.Certificate.Value,
	}
}

// FlattenSAML takes a SSOSAML instance and creates a map
func FlattenSAML(sso models.SSOSAML) map[string]interface{} {
	return map[string]interface{}{
		"metadata_url": sso.MetadataURL,
		"acs_url":      sso.AcsURL,
		"sls_url":      sso.SlsURL,
		"issuer":       sso.Issuer,
	}
}

// FlattenSSO takes an SSO interface and creates a map based on its actual type
func FlattenSSO(sso interface{}) map[string]interface{} {
	// Check if it's a SAML SSO
	if samlSSO, ok := sso.(models.SSOSAML); ok {
		return FlattenSAML(samlSSO)
	}

	// Check if it's an OpenID SSO
	if oidcSSO, ok := sso.(models.SSOOpenId); ok {
		return FlattenOIDC(oidcSSO)
	}

	// Return empty map if sso has unknown type
	return map[string]interface{}{}
}

// FlattenCert takes an SSO interface and creates a certificate map if it's a SAML app
func FlattenCert(sso interface{}) map[string]interface{} {
	// Check if it's a SAML SSO
	if samlSSO, ok := sso.(models.SSOSAML); ok {
		return FlattenSAMLCert(samlSSO)
	}

	// Return empty map if sso has unknown type or is not SAML
	return map[string]interface{}{}
}

// Flatten takes an interface{} that is likely a map[string]interface{} from the API response and
// transforms it into a map for the Terraform schema
func Flatten(ssoData map[string]interface{}) map[string]interface{} {
	tfMap := map[string]interface{}{}

	// Set known fields if they exist
	if metadataURL, ok := ssoData["metadata_url"].(string); ok {
		tfMap["metadata_url"] = metadataURL
	}

	if acsURL, ok := ssoData["acs_url"].(string); ok {
		tfMap["acs_url"] = acsURL
	}

	if slsURL, ok := ssoData["sls_url"].(string); ok {
		tfMap["sls_url"] = slsURL
	}

	if issuer, ok := ssoData["issuer"].(string); ok {
		tfMap["issuer"] = issuer
	}

	if clientID, ok := ssoData["client_id"].(string); ok {
		tfMap["client_id"] = clientID
	}

	// Return the flattened map
	return tfMap
}
