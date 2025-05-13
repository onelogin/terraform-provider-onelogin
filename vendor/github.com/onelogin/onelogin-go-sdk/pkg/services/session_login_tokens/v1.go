package sessionlogintokens

import (
	"encoding/json"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errSessionLoginTokenV1Context = "Session Login Tokens v2 service"

type V1Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new apps service v2.
func New(repo services.Repository, host string) *V1Service {
	return &V1Service{
		Endpoint:     fmt.Sprintf("%s/api/1/login/auth", host),
		Repository:   repo,
		ErrorContext: errSessionLoginTokenV1Context,
	}
}

// Create takes a SessionLoginToken request that represents an end-user's credentials
// and returns a Session Token that represents an authenticated session
func (svc *V1Service) Create(request *SessionLoginTokenRequest) (*SessionLoginToken, error) {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    request,
	})
	if err != nil {
		return nil, err
	}
	var newSessionToken SessionLoginToken
	json.Unmarshal(resp, &newSessionToken)
	return &newSessionToken, nil
}
