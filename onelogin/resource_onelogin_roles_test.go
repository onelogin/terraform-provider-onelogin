package onelogin

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	roleschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/role"
	"github.com/stretchr/testify/assert"
)

func TestAccRole_crud(t *testing.T) {
	base := GetFixture("onelogin_role_example.tf", t)
	update := GetFixture("onelogin_role_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_roles.executive_admin", "name", "executive admin"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_roles.executive_admin", "name", "updated executive admin"),
				),
			},
		},
	})
}

// TestExtractMemberIDs tests the extractMemberIDs helper that converts sub-endpoint
// object arrays (GET /api/2/roles/{id}/apps|users|admins) into flat integer ID slices.
// Sub-endpoints return objects like {"id": 1, "name": "..."} rather than bare IDs.
func TestExtractMemberIDs(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected []int
	}{
		"nil input returns empty slice": {
			input:    nil,
			expected: []int{},
		},
		"empty array returns empty slice": {
			input:    []interface{}{},
			expected: []int{},
		},
		"objects with id fields": {
			input: []interface{}{
				map[string]interface{}{"id": float64(1), "name": "app-one"},
				map[string]interface{}{"id": float64(2), "name": "app-two"},
			},
			expected: []int{1, 2},
		},
		"skips objects missing id field": {
			input: []interface{}{
				map[string]interface{}{"name": "no-id"},
				map[string]interface{}{"id": float64(5), "name": "has-id"},
			},
			expected: []int{5},
		},
		"non-slice input returns empty slice": {
			input:    "unexpected-string",
			expected: []int{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractMemberIDs(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// memberPage is a canned sub-endpoint response: the objects on one page plus the
// After-Cursor the API would return alongside them ("" on the last page).
type memberPage struct {
	items       []interface{}
	afterCursor string
}

// stubFetcher returns a memberFetcher that serves pages in order while recording
// the query it was called with, so tests can assert on the cursor walk itself.
func stubFetcher(pages []memberPage, calls *[]roleschema.RoleQuery) memberFetcher {
	return func(_ context.Context, q *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
		*calls = append(*calls, *q)
		if len(*calls) > len(pages) {
			return nil, nil, fmt.Errorf("unexpected request %d, only %d pages defined", len(*calls), len(pages))
		}
		page := pages[len(*calls)-1]
		return page.items, &models.PaginationInfo{AfterCursor: page.afterCursor}, nil
	}
}

func memberObj(id float64) interface{} {
	return map[string]interface{}{"id": id, "name": fmt.Sprintf("member-%v", id)}
}

// TestFetchAllMemberIDsWalksEveryPage is the regression test for the bug this
// pagination exists to prevent: role sub-endpoints are paginated, so stopping at
// the first page writes a truncated membership list into state and produces a
// permanent diff for any role larger than one page.
func TestFetchAllMemberIDsWalksEveryPage(t *testing.T) {
	var calls []roleschema.RoleQuery
	fetch := stubFetcher([]memberPage{
		{items: []interface{}{memberObj(1), memberObj(2)}, afterCursor: "cursor-2"},
		{items: []interface{}{memberObj(3), memberObj(4)}, afterCursor: "cursor-3"},
		{items: []interface{}{memberObj(5)}, afterCursor: ""},
	}, &calls)

	ids, err := fetchAllMemberIDs(context.Background(), fetch)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, ids, "every page's members should be accumulated")
	if !assert.Len(t, calls, 3, "should keep requesting until After-Cursor is empty") {
		return
	}

	// First request asks for a page size; subsequent ones carry the cursor with
	// limit and page cleared, which the V2 API requires.
	assert.Equal(t, roleschema.RoleQuery{Limit: rolePageLimit}, calls[0])
	assert.Equal(t, roleschema.RoleQuery{Cursor: "cursor-2"}, calls[1])
	assert.Equal(t, roleschema.RoleQuery{Cursor: "cursor-3"}, calls[2])
}

func TestFetchAllMemberIDsSinglePage(t *testing.T) {
	var calls []roleschema.RoleQuery
	fetch := stubFetcher([]memberPage{
		{items: []interface{}{memberObj(7)}, afterCursor: ""},
	}, &calls)

	ids, err := fetchAllMemberIDs(context.Background(), fetch)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []int{7}, ids)
	assert.Len(t, calls, 1, "an empty After-Cursor should end the walk immediately")
}

func TestFetchAllMemberIDsEmptyMembership(t *testing.T) {
	var calls []roleschema.RoleQuery
	fetch := stubFetcher([]memberPage{
		{items: []interface{}{}, afterCursor: ""},
	}, &calls)

	ids, err := fetchAllMemberIDs(context.Background(), fetch)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []int{}, ids, "a role with no members reads as empty, not nil")
}

// TestFetchAllMemberIDsPropagatesError checks that a mid-walk failure surfaces
// rather than being silently truncated into a partial membership list.
func TestFetchAllMemberIDsPropagatesError(t *testing.T) {
	callCount := 0
	fetch := func(_ context.Context, _ *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
		callCount++
		if callCount == 1 {
			return []interface{}{memberObj(1)}, &models.PaginationInfo{AfterCursor: "cursor-2"}, nil
		}
		return nil, nil, fmt.Errorf("request failed with status: 502")
	}

	ids, err := fetchAllMemberIDs(context.Background(), fetch)
	if !assert.Error(t, err) {
		return
	}
	assert.Nil(t, ids, "a failed walk must not return a partial list")
	assert.Contains(t, err.Error(), "502")
}

// TestFetchAllMemberIDsStalledCursor guards the loop against a server that keeps
// returning the same cursor, which would hang a plan instead of failing it.
func TestFetchAllMemberIDsStalledCursor(t *testing.T) {
	fetch := func(_ context.Context, _ *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
		return []interface{}{memberObj(1)}, &models.PaginationInfo{AfterCursor: "stuck"}, nil
	}

	_, err := fetchAllMemberIDs(context.Background(), fetch)
	if !assert.Error(t, err) {
		return
	}
	assert.Contains(t, err.Error(), "pagination stalled")
}
