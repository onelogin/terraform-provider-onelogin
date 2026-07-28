package onelogin

import "testing"

func TestUsersFromResponseWrappedData(t *testing.T) {
	result := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{"id": 1.0},
		},
	}

	data, err := usersFromResponse(result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("expected 1 user, got %d", len(data))
	}
}

func TestUsersFromResponseBareArray(t *testing.T) {
	result := []interface{}{
		map[string]interface{}{"id": 1.0},
		map[string]interface{}{"id": 2.0},
	}

	data, err := usersFromResponse(result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 users, got %d", len(data))
	}
}

func TestUsersFromResponseInvalidShape(t *testing.T) {
	result := map[string]interface{}{
		"unexpected": "value",
	}

	_, err := usersFromResponse(result)
	if err == nil {
		t.Fatal("expected error for invalid response shape")
	}
}
