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

// FlattenOIDCCredentials takes the sso object from an OIDC app's API response
// and transforms it into a map for the Terraform schema.
//
// An empty string is treated as "not reported" rather than stored as a present
// key. Terraform cannot express a null map value, so a key present with an empty
// string is indistinguishable from a real one at plan time: `sso.client_secret`
// would quietly evaluate to "" instead of failing. Omitting the key keeps the
// documented contract that indexing an uncaptured secret errors loudly.
func FlattenOIDCCredentials(ssoData map[string]interface{}) map[string]interface{} {
	tfMap := map[string]interface{}{}

	if clientID, ok := ssoData["client_id"].(string); ok && clientID != "" {
		tfMap["client_id"] = clientID
	}

	if clientSecret, ok := ssoData["client_secret"].(string); ok && clientSecret != "" {
		tfMap["client_secret"] = clientSecret
	}

	return tfMap
}

// RetainSecret merges a freshly flattened sso map with what is already in state,
// so that a client_secret captured at create survives later reads (OneLogin does
// not return it again).
//
// The merge has to hold three things true at once, and the callers rely on all
// three because d.Set replaces the whole map rather than patching keys:
//
//  1. A secret in the response always wins -- it is the freshest value.
//  2. A secret is retained only where its pairing with client_id is positively
//     confirmed. A re-issued client_id, or state holding a secret with no
//     client_id to check it against, both count as unconfirmed and the secret is
//     dropped. State is better off with the secret absent than confidently
//     wrong, because callers persist this value into other systems.
//  3. A response that omits client_id must not erase the one already in state.
func RetainSecret(prior interface{}, flattened map[string]interface{}) map[string]interface{} {
	// (1) The response carries a usable secret; nothing to retain.
	if s, ok := flattened["client_secret"].(string); ok && s != "" {
		return flattened
	}

	priorMap, ok := prior.(map[string]interface{})
	if !ok {
		return flattened
	}

	priorID, _ := priorMap["client_id"].(string)
	responseID, responseHasID := flattened["client_id"].(string)

	// (2) Retain only a secret whose pairing can be positively confirmed: the
	// response reports the same client_id state captured it against. Anything
	// else -- a re-issued client_id, or state carrying a secret with no client_id
	// to check it against, which legacy flatmap state can -- is treated as stale.
	//
	// Note this deliberately drops rather than keeps when the pairing is merely
	// unknown. OneLogin builds the sso object from one record, so a create
	// response carries client_id and client_secret together or neither, and a
	// secret with no id beside it does not arise from a healthy create. Between
	// an absent secret and a possibly-wrong one, absent is the recoverable
	// failure: it is visible, and it is fixed by recreating the app.
	if responseHasID && responseID != priorID {
		return flattened
	}

	out := make(map[string]interface{}, len(flattened)+1)
	for key, val := range flattened {
		out[key] = val
	}

	// (3) Keep a known client_id the response did not report.
	if !responseHasID && priorID != "" {
		out["client_id"] = priorID
	}

	// Retain the prior secret, applying the same empty-means-absent rule as the
	// response so an empty value cannot become permanent.
	if priorSecret, ok := priorMap["client_secret"].(string); ok && priorSecret != "" {
		out["client_secret"] = priorSecret
	}

	return out
}
