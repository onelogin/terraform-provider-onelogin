package appssoschema

import (
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
)

// FlattenOIDC takes a AppSSO instance and creates a map
func FlattenOIDC(sso apps.AppSso) map[string]interface{} {
	return map[string]interface{}{
		"client_id":     sso.ClientID,
		"client_secret": sso.ClientSecret,
	}
}

// FlattenSAMLCert takes a AppSSO instance and uses the Certificate node to create the map
func FlattenSAMLCert(sso apps.AppSso) map[string]interface{} {
	return map[string]interface{}{
		"name":  sso.Certificate.Name,
		"value": sso.Certificate.Value,
	}
}

// FlattenSAML takes a AppSSO instance and creates a map
func FlattenSAML(sso apps.AppSso) map[string]interface{} {
	return map[string]interface{}{
		"metadata_url": sso.MetadataURL,
		"acs_url":      sso.AcsURL,
		"sls_url":      sso.SlsURL,
		"issuer":       sso.Issuer,
	}
}
