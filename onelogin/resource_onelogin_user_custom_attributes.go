package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// validCustomAttributePosition rejects a negative position at plan time.
//
// Positions start at 1 and 0 stands in for "no position", so the two payload
// builders read a negative differently: GetOk counts -1 as set and would send it
// on create, while the update builder only treats a value above 0 as a position
// and would clear it instead. Neither is what the configuration asked for.
func validCustomAttributePosition(val interface{}, key string) (warns []string, errs []error) {
	if position, ok := val.(int); ok && position < 0 {
		errs = append(errs, fmt.Errorf("%s must be 0 or greater, got %d", key, position))
	}
	return warns, errs
}

// checkCustomAttributeDefinitionName rejects a definition that has no name. The
// API requires one on both create and update, and the schema cannot ask for it:
// a name is only needed for a definition, not for a value set on one user, and
// Required applies to both shapes alike.
//
// This is checked at apply rather than in a CustomizeDiff because the two shapes
// are told apart by user_id, which is commonly a reference to a user created in
// the same run. That value is still unknown while planning, so a plan-time check
// would read a user-value resource as a definition and demand a name it does not
// need.
func checkCustomAttributeDefinitionName(d *schema.ResourceData) error {
	if d.Get("name").(string) == "" {
		return fmt.Errorf("name is required for a custom attribute definition (shortname %q)", d.Get("shortname").(string))
	}
	return nil
}

// userCustomAttributeDefinitionCreateInput builds the "user_field" body handed to
// the create endpoint.
//
// "position" is left out unless it is actually set. OneLogin defaults a new
// definition's position to null, so omitting the key asks for exactly that.
// GetOk reports false for an unset int as well as for a literal 0, which is the
// behaviour wanted here: positions start at 1, so 0 only ever means "no
// position".
func userCustomAttributeDefinitionCreateInput(d *schema.ResourceData) map[string]interface{} {
	input := map[string]interface{}{
		"name":      d.Get("name").(string),
		"shortname": d.Get("shortname").(string),
	}

	if position, ok := d.GetOk("position"); ok {
		input["position"] = position.(int)
	}

	return input
}

// userCustomAttributeDefinitionUpdateInput builds the body handed to the update
// endpoint. Unlike create, the API takes this body unwrapped rather than nested
// under "user_field".
//
// It differs from the create input in two ways:
//
//   - "name" and "shortname" are always sent. Both are required by the API, so a
//     body carrying only a changed position would be rejected for the missing
//     fields.
//   - "position" is always sent too, as an explicit null once it is no longer
//     set. Dropping the key instead would leave the old position in place, so a
//     definition could never be put back to having none.
//
// Sending it unconditionally is what the plan already describes: a position left
// out of the configuration plans as 0, so an update that clears it is the change
// the user was shown, not a surprise.
func userCustomAttributeDefinitionUpdateInput(d *schema.ResourceData) map[string]interface{} {
	input := map[string]interface{}{
		"name":      d.Get("name").(string),
		"shortname": d.Get("shortname").(string),
		"position":  nil,
	}

	if position := d.Get("position").(int); position > 0 {
		input["position"] = position
	}

	return input
}

// UserCustomAttributes returns a resource with the CRUD methods and Terraform Schema defined
func UserCustomAttributes() *schema.Resource {
	return &schema.Resource{
		Create:   userCustomAttributesCreate,
		Read:     userCustomAttributesRead,
		Update:   userCustomAttributesUpdate,
		Delete:   userCustomAttributesDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				// Optional rather than Required because this resource has two
				// shapes: a definition, which needs a name, and a value set on
				// one user, which does not. Requiring it forced a throwaway name
				// onto every user-value resource, which the documented usage and
				// the acceptance tests never supplied. Definitions are checked
				// for it below, at apply.
				Description: "Name of the custom attribute. Required when managing a definition, ignored when setting a value on a user",
			},
			"shortname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Short name identifier for the custom attribute",
			},
			"position": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"user_id", "value"},
				ValidateFunc:  validCustomAttributePosition,
				Description:   "Ordering of the custom attribute definition. Positions start at 1; leaving this unset, or setting it to 0, leaves the definition without a position",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "User ID to set this custom attribute for (for user-specific custom attributes)",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value of the custom attribute (for user-specific custom attributes)",
			},
		},
	}
}

func userCustomAttributesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*onelogin.OneloginSDK)

	// Check if we're creating a definition or setting a value
	_, hasUserId := d.GetOk("user_id")
	_, hasValue := d.GetOk("value")

	// If we have both user_id and value, we're setting a value for a user
	if hasUserId && hasValue {
		userId := d.Get("user_id")
		value := d.Get("value")

		// Set the custom attribute value for this user
		userIdInt := userId.(int)
		userIdInt32 := int32(userIdInt)

		// Get the user first to get its current state
		user, err := client.GetUserByID(userIdInt, nil)
		if err != nil {
			log.Printf("[ERROR] Error getting user %d: %v", userIdInt, err)
			return err
		}

		// Initialize or update custom attributes
		userMap, ok := user.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse user response")
		}

		// Get or initialize custom attributes
		customAttrs, ok := userMap["custom_attributes"].(map[string]interface{})
		if !ok {
			customAttrs = make(map[string]interface{})
		}

		// Set the custom attribute value
		shortname := d.Get("shortname").(string)
		customAttrs[shortname] = value

		// Create a user object with just the custom attributes
		userUpdate := models.User{
			ID:               userIdInt32,
			CustomAttributes: customAttrs,
		}

		// Update the user
		_, err = client.UpdateUser(userIdInt, userUpdate)
		if err != nil {
			log.Printf("[ERROR] Error setting custom attribute for user %d: %v", userIdInt, err)
			return err
		}

		// For user-specific custom attributes, use {user_id}_{shortname} as the ID
		d.SetId(fmt.Sprintf("%d_%s", userIdInt, shortname))
		return userCustomAttributesRead(d, m)
	} else {
		// Otherwise, we're creating a new custom attribute definition,
		// wrapped in a user_field object as required by API
		if err := checkCustomAttributeDefinitionName(d); err != nil {
			return err
		}

		payload := map[string]interface{}{
			"user_field": userCustomAttributeDefinitionCreateInput(d),
		}

		// Create custom attribute
		result, err := client.CreateCustomAttributes(payload)
		if err != nil {
			log.Printf("[ERROR] Error creating custom attribute: %v", err)
			return err
		}

		// Extract ID from result
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse custom attribute creation response")
		}

		id, ok := resultMap["id"].(float64)
		if !ok {
			return fmt.Errorf("failed to extract custom attribute ID from response")
		}

		attributeID := int(id)
		// For attribute definitions, prefix the ID with "attr_" to distinguish from user attribute values
		d.SetId(fmt.Sprintf("attr_%d", attributeID))

		return userCustomAttributesRead(d, m)
	}
}

func userCustomAttributesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*onelogin.OneloginSDK)

	// Special case for a new custom attribute definition resource
	if d.Id() == "" {
		shortname := d.Get("shortname").(string)
		if shortname != "" {
			d.SetId(shortname)
		}
		return nil
	}

	// Check if this is an attribute definition ID (prefixed with "attr_")
	if len(d.Id()) > 5 && d.Id()[:5] == "attr_" {
		attrIdStr := d.Id()[5:] // Remove the prefix
		attrId, err := strconv.Atoi(attrIdStr)
		if err != nil {
			return fmt.Errorf("invalid attribute ID format: %v", err)
		}

		// Get all custom attributes
		attributes, err := client.GetCustomAttributes()
		if err != nil {
			return fmt.Errorf("error retrieving custom attributes: %v", err)
		}

		attrList, ok := attributes.([]interface{})
		if !ok {
			return fmt.Errorf("invalid custom attributes response format")
		}

		// Find the attribute with matching ID
		for _, attr := range attrList {
			attrMap, ok := attr.(map[string]interface{})
			if !ok {
				continue
			}

			id, ok := attrMap["id"].(float64)
			if !ok {
				continue
			}

			if int(id) == attrId {
				d.Set("name", attrMap["name"])
				d.Set("shortname", attrMap["shortname"])
				// A definition without a position reports null, which decodes to
				// no float64. 0 is what the schema reports for an unset int, so
				// the two agree and an unmanaged position is not seen as drift.
				if position, ok := attrMap["position"].(float64); ok {
					d.Set("position", int(position))
				} else {
					d.Set("position", 0)
				}
				return nil
			}
		}

		// If we get here, the attribute wasn't found
		d.SetId("")
		return nil
	}

	// Check if this is a user-specific custom attribute (format: "userId_shortname")
	parts := splitUserCustomAttributeId(d.Id())
	if len(parts) == 2 {
		userId, err := strconv.Atoi(parts[0])
		if err != nil {
			// If the first part isn't a number, this might be a shortname for an attribute definition
			d.Set("shortname", d.Id())
			return nil
		}

		shortname := parts[1]
		// Read the user to get their custom attributes
		user, err := client.GetUserByID(userId, nil)
		if err != nil {
			return fmt.Errorf("error reading user %d: %v", userId, err)
		}

		userMap, ok := user.(map[string]interface{})
		if !ok {
			d.SetId("")
			return nil
		}

		customAttrs, ok := userMap["custom_attributes"].(map[string]interface{})
		if !ok || customAttrs == nil {
			// No custom attributes found
			d.SetId("")
			return nil
		}

		value, ok := customAttrs[shortname]
		if !ok {
			// Custom attribute not found for this user
			d.SetId("")
			return nil
		}

		d.Set("user_id", userId)
		d.Set("shortname", shortname)
		d.Set("value", value)

		return nil
	}

	// For shortname-based IDs, just set the shortname
	shortname := d.Id()
	d.Set("shortname", shortname)

	return nil
}

func userCustomAttributesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*onelogin.OneloginSDK)

	// Check if this is an attribute definition ID (prefixed with "attr_")
	if len(d.Id()) > 5 && d.Id()[:5] == "attr_" {
		attrIdStr := d.Id()[5:] // Remove the prefix
		attrId, err := strconv.Atoi(attrIdStr)
		if err != nil {
			return fmt.Errorf("invalid attribute ID format: %v", err)
		}

		if err := checkCustomAttributeDefinitionName(d); err != nil {
			return err
		}

		// Update the custom attribute
		_, err = client.UpdateCustomAttributes(attrId, userCustomAttributeDefinitionUpdateInput(d))
		if err != nil {
			log.Printf("[ERROR] Error updating custom attribute %d: %v", attrId, err)
			return err
		}

		return userCustomAttributesRead(d, m)
	}

	// Check if this is a user-specific custom attribute
	parts := splitUserCustomAttributeId(d.Id())
	if len(parts) == 2 {
		userId, err := strconv.Atoi(parts[0])
		if err != nil {
			// This might be a shortname-based ID, just update the state
			shortname := d.Get("shortname").(string)
			d.SetId(shortname)
			return userCustomAttributesRead(d, m)
		}

		shortname := parts[1]
		userIdInt32 := int32(userId)

		// Get the user to update custom attributes
		user, err := client.GetUserByID(userId, nil)
		if err != nil {
			return fmt.Errorf("error reading user %d: %v", userId, err)
		}

		userMap, ok := user.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse user response")
		}

		// Get or initialize custom attributes
		customAttrs, ok := userMap["custom_attributes"].(map[string]interface{})
		if !ok {
			customAttrs = make(map[string]interface{})
		}

		// Update the attribute
		customAttrs[shortname] = d.Get("value")

		// Create a user object with just the custom attributes
		userUpdate := models.User{
			ID:               userIdInt32,
			CustomAttributes: customAttrs,
		}

		// Update the user
		_, err = client.UpdateUser(userId, userUpdate)
		if err != nil {
			log.Printf("[ERROR] Error updating custom attribute for user %d: %v", userId, err)
			return err
		}

		return userCustomAttributesRead(d, m)
	}

	// Otherwise, just update the shortname value in the state
	shortname := d.Get("shortname").(string)
	d.SetId(shortname)

	return userCustomAttributesRead(d, m)
}

func userCustomAttributesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*onelogin.OneloginSDK)

	// Check if this is an attribute definition ID (prefixed with "attr_")
	if len(d.Id()) > 5 && d.Id()[:5] == "attr_" {
		attrIdStr := d.Id()[5:] // Remove the prefix
		attrId, err := strconv.Atoi(attrIdStr)
		if err != nil {
			return fmt.Errorf("invalid attribute ID format: %v", err)
		}

		// Delete the custom attribute
		_, err = client.DeleteCustomAttributes(attrId)
		if err != nil {
			log.Printf("[ERROR] Error deleting custom attribute %d: %v", attrId, err)
			return err
		}

		d.SetId("")
		return nil
	}

	// Check if this is a user-specific custom attribute
	parts := splitUserCustomAttributeId(d.Id())
	if len(parts) == 2 {
		userId, err := strconv.Atoi(parts[0])
		if err != nil {
			// This might be a shortname-based ID, just clear the ID
			d.SetId("")
			return nil
		}

		shortname := parts[1]
		userIdInt32 := int32(userId)

		// Get the user to update custom attributes
		user, err := client.GetUserByID(userId, nil)
		if err != nil {
			return fmt.Errorf("error reading user %d: %v", userId, err)
		}

		userMap, ok := user.(map[string]interface{})
		if !ok {
			d.SetId("")
			return nil
		}

		// Get custom attributes
		customAttrs, ok := userMap["custom_attributes"].(map[string]interface{})
		if ok && customAttrs != nil {
			// Remove the attribute by setting it to nil
			customAttrs[shortname] = nil

			// Create a user object with just the custom attributes
			userUpdate := models.User{
				ID:               userIdInt32,
				CustomAttributes: customAttrs,
			}

			// Update the user
			_, err = client.UpdateUser(userId, userUpdate)
			if err != nil {
				log.Printf("[ERROR] Error clearing custom attribute for user %d: %v", userId, err)
				return err
			}
		}

		d.SetId("")
		return nil
	}

	// For any other IDs, just remove the ID
	d.SetId("")
	return nil
}

// Helper function to split a user custom attribute ID in the format "userId_shortname"
func splitUserCustomAttributeId(id string) []string {
	var result []string
	var currentPart string

	underscore := false
	for i, c := range id {
		if c == '_' && !underscore {
			result = append(result, currentPart)
			currentPart = ""
			underscore = true
		} else {
			currentPart += string(c)
		}

		// If this is the last character, add the current part
		if i == len(id)-1 {
			result = append(result, currentPart)
		}
	}

	return result
}
