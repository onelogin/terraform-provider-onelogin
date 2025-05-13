package olhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/onelogin/onelogin-go-sdk/internal/customerrors"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/utils"
)

const resourceRequestuestContext = "ol http service"

var errInvalidRequestInput = errors.New("Invalid input for request creation")

type OLHTTPService struct {
	ErrorContext     string
	Config           services.HTTPServiceConfig
	ClientCredential ClientCredential
}

// New uses the cfg to generate the new auth service, and returns
// the created auth service for version 2.
func New(cfg services.HTTPServiceConfig) *OLHTTPService {
	return &OLHTTPService{
		Config:           cfg,
		ErrorContext:     resourceRequestuestContext,
		ClientCredential: ClientCredential{},
	}
}

// Read exectues the HTTP GET method with the given url and query parameters passed in the payload.
// It will retrieve all available resources from the remote, up to the specified "limit" if given, that meet the query criteria.
// This also implies that if a url is for one resource i.e. /resources/:id then only that resource will be returned

// This assumes pagination is implemented using a After-Cursor response header and the Read method
// will use this until the remote stops responding with the After-Cursor header, indicating we have
// run out of pages. This is not to be confused with the Cursor header which can be set by the
// caller in the request to Read which will offset the remote query starting at that page.
func (svc OLHTTPService) Read(r interface{}) ([][]byte, error) {
	resourceRequest := r.(OLHTTPRequest)
	req, reqErr := http.NewRequest(http.MethodGet, resourceRequest.URL, nil)
	if reqErr != nil {
		return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, reqErr)
	}

	if resourceRequest.Payload != nil {
		if err := attachQueryParameters(req, resourceRequest.Payload); err != nil {
			return nil, err
		}
	}

	var (
		allData [][]byte
		next    string
	)

	if os.Getenv("OL_LOG_LEVEL") == "debug" {
		log.Printf("[ONELOGIN HTTP DEBUG] Making Read Request to %s with parameters: %s \n", resourceRequest.URL, resourceRequest.Payload)
	}
	resp, data, err := svc.executeHTTP(req, resourceRequest)
	if err != nil {
		return [][]byte{}, err
	}
	for {
		allData = append(allData, data)
		next = resp.Header.Get("After-Cursor")
		if next == "" {
			break
		}
		totalPages, _ := strconv.Atoi(resp.Header.Get("Total-Pages"))
		currentPage, _ := strconv.Atoi(resp.Header.Get("Current-Page"))
		if currentPage > totalPages {
			break
		}
		params := req.URL.Query()
		// respect page size if given
		if params.Has("limit") {
			break
		}
		params.Set("cursor", next)
		req.URL.RawQuery = params.Encode()
		resp, data, err = svc.executeHTTP(req, resourceRequest)

		if err != nil {
			return [][]byte{}, err
		}
	}
	return allData, err
}

// Create creates a new resource in the remote location over HTTP
func (svc OLHTTPService) Create(r interface{}) ([]byte, error) {
	resourceRequest := r.(OLHTTPRequest)
	var (
		req    *http.Request
		reqErr error
	)
	if resourceRequest.Payload != nil {
		bodyToSend, marshErr := json.Marshal(resourceRequest.Payload)
		if marshErr != nil {
			return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, marshErr)
		}
		if os.Getenv("OL_LOG_LEVEL") == "debug" {
			log.Printf("[ONELOGIN HTTP DEBUG] Making Create Request to %s with payload: %s \n", resourceRequest.URL, bodyToSend)
		}
		req, reqErr = http.NewRequest(http.MethodPost, resourceRequest.URL, bytes.NewBuffer(bodyToSend))
	} else {
		if os.Getenv("OL_LOG_LEVEL") == "debug" {
			log.Printf("[ONELOGIN HTTP DEBUG] Making Create Request to %s with no payload \n", resourceRequest.URL)
		}
		req, reqErr = http.NewRequest(http.MethodPost, resourceRequest.URL, bytes.NewBuffer([]byte("")))
	}
	if reqErr != nil {
		return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, reqErr)
	}
	_, data, err := svc.executeHTTP(req, resourceRequest)
	if err != nil {
		return []byte{}, err
	}
	return data, err
}

// Update updates a resource in its remote location over HTTP
func (svc OLHTTPService) Update(r interface{}) ([]byte, error) {
	resourceRequest := r.(OLHTTPRequest)
	var (
		req    *http.Request
		reqErr error
	)
	if resourceRequest.Payload != nil {
		bodyToSend, marshErr := json.Marshal(resourceRequest.Payload)
		if marshErr != nil {
			return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, marshErr)
		}
		if os.Getenv("OL_LOG_LEVEL") == "debug" {
			log.Printf("[ONELOGIN HTTP DEBUG] Making Update Request to %s with payload: %s \n", resourceRequest.URL, bodyToSend)
		}
		req, reqErr = http.NewRequest(http.MethodPut, resourceRequest.URL, bytes.NewBuffer(bodyToSend))
	} else {
		if os.Getenv("OL_LOG_LEVEL") == "debug" {
			log.Printf("[ONELOGIN HTTP DEBUG] Making Create Request to %s with no payload: \n", resourceRequest.URL)
		}
		req, reqErr = http.NewRequest(http.MethodPut, resourceRequest.URL, bytes.NewBuffer([]byte("")))
	}
	if reqErr != nil {
		return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, reqErr)
	}
	_, data, err := svc.executeHTTP(req, resourceRequest)
	if err != nil {
		return []byte{}, err
	}
	return data, err
}

// Destroy executes a HTTP destroy and removes the resource from its location in a remote
func (svc OLHTTPService) Destroy(r interface{}) ([]byte, error) {
	resourceRequest := r.(OLHTTPRequest)
	var (
		req    *http.Request
		reqErr error
	)
	if resourceRequest.Payload != nil {
		bodyToSend, marshErr := json.Marshal(resourceRequest.Payload)
		if marshErr != nil {
			return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, marshErr)
		}
		req, reqErr = http.NewRequest(http.MethodDelete, resourceRequest.URL, bytes.NewBuffer(bodyToSend))
	} else {
		req, reqErr = http.NewRequest(http.MethodDelete, resourceRequest.URL, nil)
	}
	if reqErr != nil {
		return nil, customerrors.OneloginErrorWrapper(resourceRequestuestContext, reqErr)
	}
	if os.Getenv("OL_LOG_LEVEL") == "debug" {
		log.Printf("[ONELOGIN HTTP DEBUG] Making Delete Request to %s \n", resourceRequest.URL)
	}
	_, data, err := svc.executeHTTP(req, resourceRequest)
	if err != nil {
		return []byte{}, err
	}
	return data, err
}

// creates a http query string from the request payload
func attachQueryParameters(req *http.Request, payload interface{}) error {
	params := req.URL.Query()
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var m map[string]string
	if err = json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		if v != "" {
			params.Add(utils.ToSnakeCase(k), v)
		}
	}
	req.URL.RawQuery = params.Encode()
	return nil
}

// attaches http request headers supplied by caller and auth headers depending
// on the request's auth type (e.g. bearer or basic)
func (svc *OLHTTPService) attachHeaders(req *http.Request, resourceRequest OLHTTPRequest) error {
	// set headers
	for key, val := range resourceRequest.Headers {
		req.Header.Set(key, val)
	}
	switch strings.ToLower(resourceRequest.AuthMethod) {
	case "bearer":
		if (svc.ClientCredential == ClientCredential{}) {
			if err := SetBearerToken(svc); err != nil {
				return err
			}
		}
		if svc.ClientCredential.AccessToken != nil {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *svc.ClientCredential.AccessToken))
		}
	case "basic":
		req.SetBasicAuth(svc.Config.ClientID, svc.Config.ClientSecret)
	default:
		// no auth headers
	}
	return nil
}

// executes the http request, initiates retry on expired bearer tokens and returns the
// response's byte array resource representation
func (svc *OLHTTPService) executeHTTP(req *http.Request, resourceRequest OLHTTPRequest) (*http.Response, []byte, error) {
	if err := svc.attachHeaders(req, resourceRequest); err != nil {
		return nil, nil, err
	}
	resp, err := svc.Config.Client.Do(req)
	if err != nil {
		log.Println("Executing Request To", req.URL, "With", resourceRequest.Payload)
		log.Println("HTTP Transport Error", err)
		return nil, nil, customerrors.ReqErrorWrapper(resp, svc.ErrorContext, err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, customerrors.ReqErrorWrapper(resp, svc.ErrorContext, err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusUnauthorized, resp.StatusCode == http.StatusForbidden:
		if resourceRequest.AuthMethod == "bearer" {
			if err := SetBearerToken(svc); err != nil {
				return nil, nil, err
			}
			return svc.executeHTTP(req, resourceRequest)
		}
		return nil, nil, customerrors.OneloginErrorWrapper(svc.ErrorContext, errors.New("unauthorized"))
	case resp.StatusCode >= 400 && resp.StatusCode <= 499:
		return nil, nil, customerrors.OneloginErrorWrapper(svc.ErrorContext, errors.New(string(responseData)))
	case resp.StatusCode >= 500 && resp.StatusCode <= 599:
		return nil, nil, customerrors.OneloginErrorWrapper(svc.ErrorContext, errors.New("unable to connect"))
	default:
		return resp, responseData, nil
	}
}

// requests a fresh access token
func (svc *OLHTTPService) mintBearerToken() (ClientCredential, error) {
	resp, err := svc.Create(OLHTTPRequest{
		URL:        fmt.Sprintf("%s/auth/oauth2/v2/token", svc.Config.BaseURL),
		Headers:    map[string]string{"Content-Type": "application/json"},
		Payload:    AuthBody{GrantType: "client_credentials"},
		AuthMethod: "basic",
	})
	if err != nil {
		return ClientCredential{}, err
	}

	var output ClientCredential
	if err = json.Unmarshal(resp, &output); err != nil {
		return ClientCredential{}, err
	}
	return output, nil
}

// force overwrite the service's memoized access token
func SetBearerToken(svc *OLHTTPService) error {
	cred, err := svc.mintBearerToken()
	svc.ClientCredential = cred
	if err != nil {
		return err
	}
	return nil
}
