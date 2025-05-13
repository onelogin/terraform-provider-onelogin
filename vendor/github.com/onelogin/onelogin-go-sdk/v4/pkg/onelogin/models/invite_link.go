package models

type Invite struct {
	Email         string `json:"email"`                    // The user email to generate the invite link
	PersonalEmail string `json:"personal_email,omitempty"` // Optional: the alternate email address to send the invite link
}
