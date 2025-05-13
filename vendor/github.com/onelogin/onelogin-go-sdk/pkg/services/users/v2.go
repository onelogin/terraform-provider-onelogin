package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errUsersV2Context = "users v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V2Service {
	return &V2Service{
		Endpoint:     fmt.Sprintf("%s/api/2/users", host),
		Repository:   repo,
		ErrorContext: errUsersV2Context,
	}
}

// Query retrieves all the users from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all users
func (svc *V2Service) Query(query *UserQuery) ([]User, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}
	var users []User
	for _, bytes := range resp {
		var unmarshalled []User
		json.Unmarshal(bytes, &unmarshalled)
		if len(users) == 0 {
			users = unmarshalled
		} else {
			users = append(users, unmarshalled...)
		}
	}
	return users, nil
}

// GetOne retrieves the user by id and returns it
func (svc *V2Service) GetOne(id int32) (*User, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}
	var user User

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	json.Unmarshal(resp[0], &user)
	return &user, nil
}

// GetApps retrieves the list of apps for a given user by id, it returns
// an array of apps for the user.
func (svc *V2Service) GetApps(id int32) ([]UserApp, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/apps", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}
	var apps []UserApp

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}
	log.Println(string(resp[0]))
	err = json.Unmarshal(resp[0], &apps)
	return apps, err
}

// Create takes a user without an id and attempts to use the parameters to create it
// in the API. Modifies the user in place, or returns an error if one occurs
func (svc *V2Service) Create(user *User) error {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    user,
	})
	if err != nil {
		return err
	}
	json.Unmarshal(resp, user)
	return nil
}

// Update takes a user and an id and attempts to use the parameters to update it
// in the API. Modifies the user in place, or returns an error if one occurs
func (svc *V2Service) Update(user *User) error {
	if user.ID == nil {
		return errors.New("No ID Given")
	}
	resp, err := svc.UpdateRaw(*user.ID, user)
	if err != nil {
		return err
	}
	json.Unmarshal(resp, user)
	return nil
}

// UpdateRaw takes a user and an id and attempts to use the parameters to update it
// in the API. Returns the raw response bytes or an error.
func (svc *V2Service) UpdateRaw(id int32, user interface{}) ([]byte, error) {
	return svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    user,
	})
}

// Destroy deletes the user with the given id, and if successful, it returns nil
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
