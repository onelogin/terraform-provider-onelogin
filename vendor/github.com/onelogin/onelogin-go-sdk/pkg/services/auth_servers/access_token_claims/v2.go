package accesstokenclaims

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errAppsV2Context = "access token claims v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V2Service {
	return &V2Service{
		Endpoint:     fmt.Sprintf("%s/api/2/api_authorizations", host),
		Repository:   repo,
		ErrorContext: errAppsV2Context,
	}
}

// Query retrieves all the access token claims from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all access token claims.
func (svc *V2Service) Query(query *AccessTokenClaimsQuery) ([]AccessTokenClaim, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s/claims", svc.Endpoint, query.AuthServerID),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}

	var accessTokenClaims []AccessTokenClaim
	for _, bytes := range resp {
		var unmarshalled []AccessTokenClaim
		json.Unmarshal(bytes, &unmarshalled)
		accessTokenClaims = append(accessTokenClaims, unmarshalled...)
	}
	return accessTokenClaims, nil
}

// Create creates a new access token claim in place and returns an error if something went wrong
func (svc *V2Service) Create(accessTokenClaim *AccessTokenClaim) error {
	if accessTokenClaim.AuthServerID == nil {
		return errors.New("AuthServerID required on the payload")
	}
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/claims", svc.Endpoint, *accessTokenClaim.AuthServerID),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    accessTokenClaim,
	})
	if err != nil {
		return err
	}
	respObj := map[string]int32{}
	json.Unmarshal(resp, &respObj)
	accessTokenClaim.ID = oltypes.Int32(respObj["id"])
	return nil
}

// Update updates an existing access token claim in place or returns an error if something went wrong
func (svc *V2Service) Update(accessTokenClaim *AccessTokenClaim) error {
	if accessTokenClaim.ID == nil || accessTokenClaim.AuthServerID == nil {
		return errors.New("Both ID and AuthServerID are required on the payload")
	}
	_, err := svc.UpdateRaw(*accessTokenClaim.AuthServerID, *accessTokenClaim.ID, accessTokenClaim)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRaw updates an existing access token claim and returns the raw response or an error if something went wrong
func (svc *V2Service) UpdateRaw(authServerId int32, claimId, accessTokenClaim interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/claims/%d", svc.Endpoint, authServerId, claimId),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    accessTokenClaim,
	})
}

// Destroy takes the access token claim id and access token claim id and removes the access token claim from the API.
// Returns an error if something went wrong.
func (svc *V2Service) Destroy(accessTokenClaimId int32, id int32) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/claims/%d", svc.Endpoint, accessTokenClaimId, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}
