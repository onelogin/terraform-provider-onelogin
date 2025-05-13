package models

type Rule struct {
	ID          string   `json:"id"`          // Unique identifier for the rule
	Name        string   `json:"name"`        // Name of the rule
	Description string   `json:"description"` // Description of the rule
	Type        string   `json:"type"`        // Type of the rule (e.g., "blacklist")
	Target      string   `json:"target"`      // Target field for the rule (e.g., "location.ip")
	Source      string   `json:"source"`      // Source field for the rule (e.g., "guest-123")
	Filters     []string `json:"filters"`     // List of filters (IP addresses, etc.)
}
