package legalvalues

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errLegalValuesContext = "legal values service"

// V2Service holds the information needed to interface with a repository
type LegalValuesService struct {
	BaseURL, ErrorContext string
	Repository            services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *LegalValuesService {
	return &LegalValuesService{
		BaseURL:      fmt.Sprintf("%s/api/2", host),
		Repository:   repo,
		ErrorContext: errLegalValuesContext,
	}
}

func (svc *LegalValuesService) Query(address string, outShape interface{}) error {
	respBytes, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		AuthMethod: "bearer",
		URL:        fmt.Sprintf("%s/%s", svc.BaseURL, address),
		Headers:    map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		return err
	}

	for _, respByte := range respBytes {
		if err := json.Unmarshal(respByte, outShape); err != nil {
			return errors.New("Unable to unpack response")
		}
	}

	return nil
}
