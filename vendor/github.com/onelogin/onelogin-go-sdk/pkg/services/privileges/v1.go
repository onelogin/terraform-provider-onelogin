package privileges

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errPrivilegesV2Context = "privileges v2 service"

// V1Service holds the information needed to interface with a repository
type V1Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V1Service {
	return &V1Service{
		Endpoint:     fmt.Sprintf("%s/api/1/privileges", host),
		Repository:   repo,
		ErrorContext: errPrivilegesV2Context,
	}
}

// Query retrieves all the privileges from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all privileges.
func (svc *V1Service) Query(query *PrivilegeQuery) ([]Privilege, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var privileges []Privilege
	for _, bytes := range resp {
		var unmarshalled []Privilege
		json.Unmarshal(bytes, &unmarshalled)
		privileges = append(privileges, unmarshalled...)
	}

	return privileges, err
}

// QueryWithAssignment retrieves all the privileges from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all privileges.
// Returns an array of errors if retrieving associated resources fails. Errors array is
// index-mapped to privileges array so you can retry the associated resources requests for
// privileges whos associated resource requests have run afoul.
func (svc *V1Service) QueryWithAssignments(query *PrivilegeQuery) ([]Privilege, []error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, []error{err}
	}

	var privileges []Privilege
	for _, bytes := range resp {
		var unmarshalled []Privilege
		json.Unmarshal(bytes, &unmarshalled)
		privileges = append(privileges, unmarshalled...)
	}
	errs := make([]error, len(privileges))
	for i := range privileges {
		e := svc.GetPrivilegeResources(&privileges[i])
		if e != nil {
			errs[i] = e
		}
	}

	return privileges, errs
}

// GetOne retrieves the privilege and assigned roles/users by id, and if successful, it returns
// a pointer to the privilege. If retrieval of roles or users fails, an error is returned also.
func (svc *V1Service) GetOne(id string) (*Privilege, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}

	var privilege Privilege

	if len(resp) < 1 {
		return nil, errors.New("invalid length of response returned")
	}

	json.Unmarshal(resp[0], &privilege)

	err = svc.GetPrivilegeResources(&privilege)
	return &privilege, err
}

// Create creates a new privilege and attaches resources to it.
// Returns an error and reverts API create on privilege if something went wrong.
func (svc *V1Service) Create(privilege *Privilege) error {
	newUsers := privilege.UserIDs
	newRoles := privilege.RoleIDs
	privilege.UserIDs = nil
	privilege.RoleIDs = nil
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    privilege,
	})
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resp, privilege); err != nil {
		return err
	}

	privilege.UserIDs = newUsers
	privilege.RoleIDs = newRoles

	if err = svc.AttachPrivilegeResources(privilege); err != nil {
		svc.Destroy(*privilege.ID)
		return err
	}

	return nil
}

// Update takes a privilege and an id and attempts to use the parameters to update it
// in the API. Modifies the privilege in place, or returns an error if one occurs
func (svc *V1Service) Update(privilege *Privilege) error {
	if privilege.ID == nil {
		return errors.New("No ID Given")
	}

	// save these off since theyre not part of the privilege API call. need to re-attach after.
	keepUsers := privilege.UserIDs
	keepRoles := privilege.RoleIDs

	resp, err := svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, *privilege.ID),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    privilege,
	})
	if err != nil {
		return err
	}

	json.Unmarshal(resp, privilege)
	privilege.UserIDs = keepUsers
	privilege.RoleIDs = keepRoles

	// attach will append new resource ids to privilege and do nothing if id is removed or unchanged.
	if err = svc.AttachPrivilegeResources(privilege); err != nil {
		fmt.Println("unable to update assigned resources, reverting privilege to last known state in remote", err)
		_, err = svc.GetOne(*privilege.ID)
		return err
	}

	privilege, err = svc.GetOne(*privilege.ID)
	if err != nil {
		return err
	}

	usersToRemove := reconcileForDiscard(privilege.UserIDs, keepUsers)
	rolesToRemove := reconcileForDiscard(privilege.RoleIDs, keepRoles)

	if err = svc.DiscardAssignment(*privilege.ID, "users", usersToRemove); err != nil {
		fmt.Println("unable to remove users, reverting privilege to last known state in remote")
		_, err = svc.GetOne(*privilege.ID)
		return err
	}

	if err = svc.DiscardAssignment(*privilege.ID, "roles", rolesToRemove); err != nil {
		fmt.Println("unable to remove roles, reverting privilege to last known state in remote")
		_, err = svc.GetOne(*privilege.ID)
		return err
	}

	return nil
}

// Destroy deletes the privilege with the given id, and if successful, it returns nil
func (svc *V1Service) Destroy(id string) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}

func (svc *V1Service) DiscardAssignment(pID, resourceType string, assignments []int) error {
	var discardError error
	if len(assignments) > 0 {
		c := make(chan bool, len(assignments))
		var wg sync.WaitGroup
		for _, n := range assignments {
			wg.Add(1)
			go func(id int, c chan bool, wg *sync.WaitGroup) {
				defer wg.Done()
				_, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
					URL:        fmt.Sprintf("%s/%s/%s/%d", svc.Endpoint, pID, resourceType, id),
					Headers:    map[string]string{"Content-Type": "application/json"},
					AuthMethod: "bearer",
				})
				c <- err != nil
			}(n, c, &wg)
			wg.Wait()
			close(c)
		}
		for r := range c {
			if r {
				discardError = errors.New("unable to remove resources")
			}
		}
	}
	return discardError
}

type AttachedResponse struct {
	errResource string
	out         []int
}

// GetPrivilegeResources takes a privilege and requests the RoleIDs and UserIDs associated
// with the privilege and attaches an array of Role and User Ids to the privilege object.
// Returns an error if any request fails. Will only add Ids from successful calls
// It is possible to be in a 'half-finished' state and caller is responsible for retry or deciding next steps
func (svc *V1Service) GetPrivilegeResources(p *Privilege) error {
	u := make(chan AttachedResponse)
	r := make(chan AttachedResponse)

	go svc.getResourcesByType("users", *p.ID, u)
	go svc.getResourcesByType("roles", *p.ID, r)

	uResp := <-u
	rResp := <-r

	if uResp.errResource == "" {
		p.UserIDs = uResp.out
	}
	if rResp.errResource == "" {
		p.RoleIDs = rResp.out
	}
	return collectErrors(uResp.errResource, rResp.errResource, "read")
}

func (svc *V1Service) getResourcesByType(resourceType, privilegeID string, c chan AttachedResponse) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%s/%s", svc.Endpoint, privilegeID, resourceType),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		c <- AttachedResponse{errResource: resourceType}
	}
	ar := map[string][]int{}
	if len(resp) > 0 {
		json.Unmarshal(resp[0], &ar)
	}

	c <- AttachedResponse{out: ar[resourceType]}
	close(c)
}

type AttachResponse struct {
	errResource string
	out         bool
}

// AttachPrivilegeResources takes a privilege with RoleIDs and UserIDs specified and attempts to create
// the association between the privilege and the requested resources.
// Returns an error if the save fails.
func (svc *V1Service) AttachPrivilegeResources(p *Privilege) error {
	u := make(chan AttachResponse)
	r := make(chan AttachResponse)

	go svc.attachResourcesByType("roles", *p.ID, p.RoleIDs, u)
	go svc.attachResourcesByType("users", *p.ID, p.UserIDs, r)

	uResp := <-u
	rResp := <-r

	return collectErrors(uResp.errResource, rResp.errResource, "assign")
}

func (svc *V1Service) attachResourcesByType(resourceType, privilegeID string, resourceIDs []int, c chan AttachResponse) {
	// assign the resources if there are resources to be assigned
	if len(resourceIDs) > 0 {
		resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
			URL:        fmt.Sprintf("%s/%s/%s", svc.Endpoint, privilegeID, resourceType),
			Headers:    map[string]string{"Content-Type": "application/json"},
			AuthMethod: "bearer",
			Payload:    map[string][]int{resourceType: resourceIDs},
		})
		if err != nil {
			c <- AttachResponse{errResource: resourceType}
		}

		ar := map[string]bool{}
		json.Unmarshal(resp, &ar)

		if ar["success"] {
			c <- AttachResponse{out: ar["success"]}
		} else {
			c <- AttachResponse{errResource: resourceType}
		}
	} else {
		c <- AttachResponse{out: true} // skip assignment flow - nothing to do here
	}
	close(c)
}

func collectErrors(usersError, rolesError, action string) error {
	erroredResources := []string{}
	if usersError != "" {
		erroredResources = append(erroredResources, usersError)
	}
	if rolesError != "" {
		erroredResources = append(erroredResources, rolesError)
	}
	if len(erroredResources) == 0 {
		return nil
	}
	return fmt.Errorf("unable to %s %v", action, erroredResources)
}

func reconcileForDiscard(all, keep []int) []int {
	k := map[int]bool{}
	discard := []int{}
	for _, n := range keep {
		k[n] = true
	}
	for _, m := range all {
		if !k[m] {
			discard = append(discard, m)
		}
	}
	return discard
}
