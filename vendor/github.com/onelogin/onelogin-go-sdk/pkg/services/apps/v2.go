package apps

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/internal/customerrors"
	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errAppsV2Context = "apps v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V2Service {
	return &V2Service{
		Endpoint:     fmt.Sprintf("%s/api/2/apps", host),
		Repository:   repo,
		ErrorContext: errAppsV2Context,
	}
}

// Query retrieves all the apps from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all apps.
func (svc *V2Service) Query(query *AppsQuery) ([]App, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, bytes := range resp {
		var unmarshalled []App
		json.Unmarshal(bytes, &unmarshalled)
		if len(apps) == 0 {
			apps = unmarshalled
		} else {
			apps = append(apps, unmarshalled...)
		}
	}

	return apps, nil
}

// GetOne retrieves the app by id, and if successful, it returns
// a pointer to the app.
func (svc *V2Service) GetOne(id int32) (*App, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}

	var app App
	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	json.Unmarshal(resp[0], &app)

	return &app, nil
}

// GetUsers retrieves the list of users for a given app by id, it returns
// an array of users for the app.
func (svc *V2Service) GetUsers(id int32) ([]AppUser, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/users", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}

	var users []AppUser
	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	err = json.Unmarshal(resp[0], &users)

	return users, err
}

// Create creates a new app, and if successful, it returns a pointer to the app.
func (svc *V2Service) Create(app *App) error {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    app,
	})
	if err != nil {
		return err
	}
	json.Unmarshal(resp, app)
	return nil
}

// Update updates an existing app in place or returns an error if something went wrong

// Update is unique in that the API does not fully support Parameters as first-class
// resources and are managed by nesting them in the App. This means that a partial
// update state could exist if, for example, a parameter failed to delete or be updated
// while other parameter changes succeeded. In order to ensure the client is given an
// accurate representation of what has been persisted to the API, we call out to the GetOne
// to simply return what is currently in the API, rather than updating in place. This is a
// temporary holdover until parameters is dealt with in a consistent fashion to other nested resources like app rules
func (svc *V2Service) Update(app *App) (*App, error) {
	if app.ID == nil {
		return nil, errors.New("No ID Given")
	}
	requestedParametersState := make(map[string]AppParameters, len(app.Parameters))
	for k, p := range app.Parameters {
		requestedParametersState[k] = p
	}
	resp, err := svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, *app.ID),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    app,
	})
	if err != nil {
		return &App{}, err
	}
	json.Unmarshal(resp, app)

	pruneParamErr := svc.pruneParameters(requestedParametersState, app)

	if pruneParamErr != nil {
		var recoverErr error
		app, recoverErr = svc.GetOne(*app.ID)
		if recoverErr != nil {
			return nil, err
		}
		return app, pruneParamErr
	}
	// re-read the app so we return one with all the parameters changes made via each individual parameters call
	return svc.GetOne(*app.ID)
}

// Destroy deletes the app for the id, and if successful, it returns nil
func (svc *V2Service) Destroy(id int32) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}

// Given a list of requested parameters, go to the API, and pluck (delete) all the parameters that are not on the
// request list. At this point the app holds all existing parameters in the API.
// Rules not on the request list are assumed to be removed by the caller.
func (svc *V2Service) pruneParameters(requestedParams map[string]AppParameters, app *App) error {
	var delErrors []error
	keepMap := make(map[int32]bool, len(requestedParams))
	for _, param := range requestedParams {
		var id int32
		if param.ID == nil {
			// If we weren't given an id for our parameter, get it from the app object
			id = *app.Parameters[*param.ParamKeyName].ID
		} else {
			id = *param.ID
		}
		keepMap[id] = true
	}
	// no need to call down app parameters. parameters returned as part of app update
	for _, delCandidate := range app.Parameters {
		if !keepMap[*delCandidate.ID] {
			if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
				URL:        fmt.Sprintf("%s/%d/parameters/%d", svc.Endpoint, *app.ID, *delCandidate.ID),
				Headers:    map[string]string{"Content-Type": "application/json"},
				AuthMethod: "bearer",
			}); err != nil {
				delErrors = append(delErrors, err)
			}
		}
	}
	return customerrors.StackErrors(delErrors)
}
