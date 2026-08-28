// Package brand holds the Terraform schema for onelogin_brands and the
// translation between it and the /api/2/branding/brands representation.
package brand

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// field is one brand attribute that maps straight onto an API field of the same
// name. The attributes that need more than a rename -- name, master, and the
// two images -- are handled separately in Schema, RequestBody and Flatten.
type field struct {
	Name string
	Type schema.ValueType
	// Localised marks a value the API stores as a JSON object of locale to
	// text, carried on the wire as a string. See localisedDescription.
	Localised   bool
	Description string
}

// fields lists every plain brand attribute the SDK's models.Brand can send.
//
// The API returns more than this -- custom_links, use_custom_smtp_setting, and
// on the master brand domain_name, allowed_redirect_urls,
// hide_login_forgot_password and the two password_self_service_* fields. None
// of them is on models.Brand, so none can be written through the SDK's
// CreateBrand/UpdateBrand, and a schema attribute that cannot be sent would be
// a lie. They are left out until the model carries them.
var fields = []field{
	{Name: "enabled", Type: schema.TypeBool, Description: "Whether the brand is enabled. Defaults to false on create."},
	{Name: "custom_support_enabled", Type: schema.TypeBool, Description: "Whether the login page offers a support request link."},

	{Name: "custom_color", Type: schema.TypeString, Description: "Primary brand colour, as a hex value such as `#000000`."},
	{Name: "custom_accent_color", Type: schema.TypeString, Description: "Secondary brand colour, as a hex value such as `#016B91`."},
	{Name: "custom_masking_color", Type: schema.TypeString, Description: "Colour of the masking layer drawn over the background image, as a hex value."},
	{Name: "custom_masking_opacity", Type: schema.TypeInt, Description: "Opacity of custom_masking_color. Defaults to 0."},

	{Name: "enable_custom_label_for_login_screen", Type: schema.TypeBool, Description: "Whether the Username/Email field on the login screen uses a custom label."},
	{Name: "custom_label_text_for_login_screen", Type: schema.TypeString, Localised: true, Description: "Label for the Username/Email field on the login screen."},

	{Name: "login_instruction_title", Type: schema.TypeString, Localised: true, Description: "Link text that opens the login instruction screen."},
	{Name: "login_instruction", Type: schema.TypeString, Localised: true, Description: "Text of the login instruction screen, styled in Markdown."},

	{Name: "hide_onelogin_footer", Type: schema.TypeBool, Description: "Whether to hide the OneLogin footer at the bottom of the login page."},
	{Name: "mfa_enrollment_message", Type: schema.TypeString, Description: "Text replacing the default message on the first screen of MFA registration."},
}

// localisedDescription explains the JSON-object-in-a-string shape shared by the
// localised attributes, appended to each of their descriptions so it appears
// wherever the argument does.
const localisedDescription = " Carried as a JSON object of locale to text, encoded as a string, e.g. `jsonencode({ en = \"Username\" })`. Whitespace and key order are not significant."

// Schema returns the onelogin_brands resource schema.
func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// Required because the API says so: a create without it is
		// 422 {"name":"Value is required."}, confirmed against the API.
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the brand.",
		},
		// Read-only. Every account has exactly one master brand, created with
		// the account; the API will not make an existing brand the master and
		// a create always returns master false.
		"master": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether this is the account's master brand. Read-only.",
		},

		// The two images are write-only, and deliberately not Computed.
		//
		// They go out as base64 and come back as an object -- file_size,
		// content_type, updated_at and a urls block of resized variants --
		// so there is nothing to write back into a string attribute. Marking
		// them Computed would leave them permanently "known after apply", and
		// flattening the object would be a type error. Both were confirmed
		// against the API.
		//
		// The cost is that a logo replaced outside Terraform is not detected
		// as drift: state keeps the base64 that was last applied. Changing the
		// configured value still replaces the image.
		"logo": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Base64-encoded PNG for the login page logo, under 1MB. Write-only: the API returns image URLs rather than the data, so a logo changed outside Terraform is not seen as drift.",
		},
		"background": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Base64-encoded JPG or PNG for the login page background, under 5MB. Write-only, for the same reason as logo.",
		},
	}

	// Every plain field is Optional and Computed. Optional because none has to
	// be set, and Computed because the API fills in a default for several of
	// them on create -- enabled false, custom_color #000000, custom_accent_color
	// #016B91, custom_masking_color #000000, custom_masking_opacity 0 -- and
	// without it a brand configured with two attributes shows a diff for the
	// rest on the next plan.
	//
	// The cost is the usual one: removing an attribute from the configuration
	// no longer resets it. State keeps the last value the API returned and
	// nothing is sent. Setting it back explicitly is the way to undo one.
	for _, f := range fields {
		description := f.Description
		attribute := &schema.Schema{
			Type:     f.Type,
			Optional: true,
			Computed: true,
		}
		if f.Localised {
			description += localisedDescription
			attribute.DiffSuppressFunc = suppressEquivalentJSON
		}
		attribute.Description = description
		s[f.Name] = attribute
	}

	return s
}

// suppressEquivalentJSON hides a diff between two JSON documents that differ
// only in formatting or key order.
//
// The localised attributes are JSON objects the API hands back re-encoded, so
// the string that comes out is rarely byte-identical to the one that went in
// even when nothing changed. Comparing the decoded values instead keeps a plan
// empty when the brand really is unchanged.
//
// Anything that does not parse is compared as a plain string, so a malformed
// value still shows a diff rather than being quietly treated as equal.
func suppressEquivalentJSON(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}

	var oldValue, newValue interface{}
	if json.Unmarshal([]byte(old), &oldValue) != nil {
		return false
	}
	if json.Unmarshal([]byte(new), &newValue) != nil {
		return false
	}
	return reflect.DeepEqual(oldValue, newValue)
}

// ConfiguredKeys returns the names of the attributes the practitioner actually
// wrote, so that an attribute left out can be told from one written as false, 0
// or "". d.Get cannot make that distinction -- it returns the zero value for
// both -- and the difference decides whether the attribute is sent.
//
// raw comes from d.GetRawConfig(). It is null outside a plan or an apply, and
// nil is returned in that case to mean "no configuration information", which is
// not the same answer as "nothing was configured"; see requestValue.
func ConfiguredKeys(raw cty.Value) map[string]bool {
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsObjectType() {
		return nil
	}

	keys := map[string]bool{}
	for name, value := range raw.AsValueMap() {
		if value.IsNull() {
			continue
		}
		keys[name] = true
	}
	return keys
}

// Getter is the part of *schema.ResourceData that RequestBody reads, so that
// the body can be built from anything holding the same values.
type Getter interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

// RequestBody returns the models.Brand for a create or an update, carrying only
// the attributes the configuration sets.
//
// Sending only configured attributes matters because every plain attribute is
// Optional and Computed, so one the practitioner left out still holds whatever
// the API last returned; sending that back would have the provider assert a
// value nobody asked for.
//
// name is always set. It is Required, so it is always configured, and it is the
// one field on models.Brand without omitempty -- a nil Name marshals to
// "name": null, which the API answers with 422 "Value must be a string." even
// though omitting the key entirely is a 200. Both were confirmed against the
// API. Always setting it keeps this resource clear of that.
//
// configured comes from ConfiguredKeys and may be nil; see requestValue.
func RequestBody(d Getter, configured map[string]bool) models.Brand {
	name, _ := d.Get("name").(string)
	body := models.Brand{Name: &name}

	for _, f := range fields {
		value, ok := requestValue(d, configured, f.Name)
		if !ok {
			continue
		}
		assign(&body, f.Name, value)
	}

	// The images are write-only, so they are sent whenever configured and
	// never read back.
	if value, ok := requestValue(d, configured, "logo"); ok {
		if text, ok := value.(string); ok && text != "" {
			body.Logo = &text
		}
	}
	if value, ok := requestValue(d, configured, "background"); ok {
		if text, ok := value.(string); ok && text != "" {
			body.Background = &text
		}
	}

	return body
}

// requestValue returns an attribute's value and whether it should be sent.
//
// When configured is nil the raw configuration was not available, and the
// fallback is d.GetOk, which treats a zero value as unset. That can drop a
// deliberate false or 0, but the alternative is sending values the practitioner
// never wrote. Terraform populates the raw configuration for create and update,
// so this only matters somewhere that does not, such as a test harness.
func requestValue(d Getter, configured map[string]bool, name string) (interface{}, bool) {
	if configured == nil {
		return d.GetOk(name)
	}
	if !configured[name] {
		return nil, false
	}
	return d.Get(name), true
}

// assign writes one attribute onto the brand. The switch is exhaustive over
// fields; a name not listed here is dropped rather than sent, which is why
// TestRequestBodyCoversEveryField exists.
func assign(body *models.Brand, name string, value interface{}) {
	switch name {
	case "enabled":
		body.Enabled = boolPtr(value)
	case "custom_support_enabled":
		body.CustomSupportEnabled = boolPtr(value)
	case "custom_color":
		body.CustomColor = stringPtr(value)
	case "custom_accent_color":
		body.CustomAccentColor = stringPtr(value)
	case "custom_masking_color":
		body.CustomMaskingColor = stringPtr(value)
	case "custom_masking_opacity":
		body.CustomMaskingOpacity = int32Ptr(value)
	case "enable_custom_label_for_login_screen":
		body.EnableCustomLabelForLoginScreen = boolPtr(value)
	case "custom_label_text_for_login_screen":
		body.CustomLabelTextForLoginScreen = stringPtr(value)
	case "login_instruction_title":
		body.LoginInstructionTitle = stringPtr(value)
	case "login_instruction":
		body.LoginInstruction = stringPtr(value)
	case "hide_onelogin_footer":
		body.HideOneLoginFooter = boolPtr(value)
	case "mfa_enrollment_message":
		body.MFAEnrollmentMessage = stringPtr(value)
	}
}

func boolPtr(value interface{}) *bool {
	flag, ok := value.(bool)
	if !ok {
		return nil
	}
	return &flag
}

func stringPtr(value interface{}) *string {
	text, ok := value.(string)
	if !ok {
		return nil
	}
	return &text
}

func int32Ptr(value interface{}) *int32 {
	number, ok := toInt(value)
	if !ok {
		return nil
	}
	converted := int32(number)
	return &converted
}

// Flatten writes a brand from the API into ResourceData.
//
// Only the keys the response carries are written. The API leaves out attributes
// a brand has never been given -- a new brand comes back without
// hide_onelogin_footer or mfa_enrollment_message at all -- and writing a zero
// value in their place would put a value in state for something the brand does
// not have.
//
// logo and background are never written back: see the note in Schema.
func Flatten(d *schema.ResourceData, brand map[string]interface{}) error {
	for _, name := range []string{"name", "master"} {
		value, ok := brand[name]
		if !ok || value == nil {
			continue
		}
		if err := d.Set(name, value); err != nil {
			return err
		}
	}

	for _, f := range fields {
		value, ok := brand[f.Name]
		if !ok || value == nil {
			continue
		}
		if f.Type == schema.TypeInt {
			number, ok := toInt(value)
			if !ok {
				continue
			}
			value = number
		}
		if err := d.Set(f.Name, value); err != nil {
			return err
		}
	}

	return nil
}

// ID reads the ID out of a brand response as the string Terraform keeps.
func ID(brand map[string]interface{}) (string, bool) {
	number, ok := toInt(brand["id"])
	if ok {
		return strconv.Itoa(number), true
	}
	if text, ok := brand["id"].(string); ok && text != "" {
		return text, true
	}
	return "", false
}

// toInt converts a number from a decoded JSON response. Numbers arrive as
// float64 from encoding/json, but a caller that built the map itself may have
// used int, so both are accepted.
func toInt(value interface{}) (int, bool) {
	switch number := value.(type) {
	case float64:
		return int(number), true
	case int:
		return number, true
	case int32:
		return int(number), true
	case int64:
		return int(number), true
	default:
		return 0, false
	}
}
