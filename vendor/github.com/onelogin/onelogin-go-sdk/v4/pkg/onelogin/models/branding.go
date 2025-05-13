package models

import "encoding/json"

type Brand struct {
	Name                            *string `json:"name"`
	Master                          *bool   `json:"master,omitempty"`
	Enabled                         *bool   `json:"enabled,omitempty"`                              // Indicates if the brand is enabled or not (default=false)
	CustomSupportEnabled            *bool   `json:"custom_support_enabled,omitempty"`               // Indicates if custom support is enabled
	CustomColor                     *string `json:"custom_color,omitempty"`                         // Primary brand color (Hex value)
	CustomAccentColor               *string `json:"custom_accent_color,omitempty"`                  // Secondary brand color (Hex value)
	CustomMaskingColor              *string `json:"custom_masking_color,omitempty"`                 // Color for the masking layer above the background image
	CustomMaskingOpacity            *int32  `json:"custom_masking_opacity,omitempty"`               // Opacity for the custom_masking_color
	EnableCustomLabelForLoginScreen *bool   `json:"enable_custom_label_for_login_screen,omitempty"` // Custom Username/Email field label for login screen
	CustomLabelTextForLoginScreen   *string `json:"custom_label_text_for_login_screen,omitempty"`   // Custom label text for Username/Email field
	LoginInstructionTitle           *string `json:"login_instruction_title,omitempty"`              // Link text for login instruction screen
	LoginInstruction                *string `json:"login_instruction,omitempty"`                    // Text for login instruction screen (Markdown)
	HideOneLoginFooter              *bool   `json:"hide_onelogin_footer,omitempty"`                 // Indicates if the OneLogin footer will appear at the bottom of the login page
	MFAEnrollmentMessage            *string `json:"mfa_enrollment_message,omitempty"`               // Custom text for MFA Registration initial screen
	Background                      *string `json:"background,omitempty"`                           // Base64 encoded image data for background (JPG/PNG, <5MB)
	Logo                            *string `json:"logo,omitempty"`                                 // Base64 encoded image data for logo (PNG, <1MB)
	ID                              int32   `json:"id,omitempty"`
}

type BrandTemplate struct {
	Type     string `json:"type"`     // Template type (e.g., email_forgot_password)
	Locale   string `json:"locale"`   // Locale (e.g., "es" for Spanish, "en" for english)
	Template string `json:"template"` // Template content (Email/SMS content)
}

// type BrandTemplate1 struct {
// 	Type     string          `json:"type,omitempty"`
// 	Locale   string          `json:"locale,omitempty"`
// 	Template json.RawMessage `json:"template,omitempty"` // Use RawMessage to avoid extra escaping
// }

type TemplateContent struct {
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Plain   string `json:"plain"`
}

func (t *BrandTemplate) MarshalJSON() ([]byte, error) {
	// Define type to break recursive calls to MarshalJSON.
	// The type X has all of fields for BrandTemplate, but
	// non of the methods.
	type X BrandTemplate

	// Marshal a type that shadows BrandTemplate.Template
	// with a raw JSON field of the same name.
	return json.Marshal(
		struct {
			*X
			Template json.RawMessage `json:"template"`
		}{
			(*X)(t),
			json.RawMessage(t.Template),
		})
}

// Custom MarshalJSON to handle the template field
// func (bt *BrandTemplate) MarshalJSON() ([]byte, error) {
// 	type Alias BrandTemplate
// 	return json.Marshal(&struct {
// 		*Alias
// 		Template string `json:"template"`
// 	}{
// 		Alias:    (*Alias)(bt),
// 		Template: bt.Template,
// 	})
// }

// // Helper function to generate JSON string from TemplateContent
// func GenerateTemplateJSON(subject, html, plain string) (string) {
// 	// Create the TemplateContent struct with dynamic values
// 	templateContent := TemplateContent{
// 		Subject: subject,
// 		HTML:    html,
// 		Plain:   plain,
// 	}

// 	// Marshal it into a JSON string
// 	encodedTemplate:= json.Marshal(templateContent)

// 	// Return the raw JSON string
// 	return string(encodedTemplate), nil
// }
