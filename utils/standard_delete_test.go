package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deletableResource(t *testing.T, id string) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}, map[string]interface{}{})
	d.SetId(id)
	return d
}

// TestStandardDeleteFuncNotFound covers a delete that finds the resource
// already gone. That is the state the delete was asking for, and reporting it
// as a failure leaves the resource in state for the next run to fail on again
// -- a destroy that cannot finish.
func TestStandardDeleteFuncNotFound(t *testing.T) {
	t.Run("treats a 404 as done", func(t *testing.T) {
		d := deletableResource(t, "12345")

		diags := StandardDeleteFunc(context.Background(), d, func(string) (interface{}, error) {
			return nil, errors.New("request failed with status: 404")
		}, "Thing")

		if diags.HasError() {
			t.Fatalf("expected a 404 to be treated as deleted, got %v", diags)
		}
		if d.Id() != "" {
			t.Fatalf("expected the id to be cleared, got %q", d.Id())
		}
	})

	t.Run("still reports anything else", func(t *testing.T) {
		d := deletableResource(t, "12345")

		diags := StandardDeleteFunc(context.Background(), d, func(string) (interface{}, error) {
			return nil, errors.New("request failed with status: 500")
		}, "Thing")

		if !diags.HasError() {
			t.Fatal("expected a 500 to be reported")
		}
		if d.Id() == "" {
			t.Fatal("expected the id to survive a failed delete: the resource is still there")
		}
	})

	t.Run("clears the id on success", func(t *testing.T) {
		d := deletableResource(t, "12345")

		if diags := StandardDeleteFunc(context.Background(), d, func(string) (interface{}, error) {
			return nil, nil
		}, "Thing"); diags.HasError() {
			t.Fatalf("expected a clean delete, got %v", diags)
		}
		if d.Id() != "" {
			t.Fatalf("expected the id to be cleared, got %q", d.Id())
		}
	})
}
