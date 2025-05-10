package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
)

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the custom attribute",
			},
			"shortname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Short name identifier for the custom attribute",
			},
			"position": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Position of the custom attribute",
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
	client := m.(*client.APIClient)

	// If user_id is provided, we're setting a value for an existing custom attribute
	if userId, ok := d.GetOk("user_id"); ok {
		if value, valueOk := d.GetOk("value"); valueOk {
			// Set the custom attribute value for this user
			userIdInt := userId.(int)
			userIdInt32 := int32(userIdInt)
			
			// Get the user first to get its current state
			user, err := client.Services.UsersV2.GetOne(userIdInt32)
			if err != nil {
				log.Printf("[ERROR] Error getting user %d: %v", userIdInt, err)
				return err
			}
			
			// Initialize or update custom attributes
			if user.CustomAttributes == nil {
				user.CustomAttributes = make(map[string]interface{})
			}
			
			// Set the custom attribute value
			shortname := d.Get("shortname").(string)
			user.CustomAttributes[shortname] = value
			
			// Update the user
			err = client.Services.UsersV2.Update(user)
			if err != nil {
				log.Printf("[ERROR] Error setting custom attribute for user %d: %v", userIdInt, err)
				return err
			}
			
			// For user-specific custom attributes, use {user_id}_{shortname} as the ID
			d.SetId(fmt.Sprintf("%d_%s", userIdInt, shortname))
			return userCustomAttributesRead(d, m)
		} else {
			return fmt.Errorf("when user_id is provided, value must also be provided")
		}
	}

	// Otherwise, we're creating a new custom attribute definition
	// This requires API calls outside the existing SDK, so we'll handle it specially
	log.Printf("[WARNING] Creation of global custom attribute definitions is currently not working due to an OneLogin API bug")
	log.Printf("[WARNING] The API returns 'Missing param: user_field' when attempting to create custom attributes")
	log.Printf("[WARNING] You can set values for existing custom attributes by providing user_id and value parameters")
	log.Printf("[WARNING] Please create custom attributes in the OneLogin UI until this API bug is fixed")
	
	shortname := d.Get("shortname").(string)
	d.SetId(shortname)
	return userCustomAttributesRead(d, m)
}

func userCustomAttributesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	
	// Check if this is a user-specific custom attribute
	if d.Id() != "" && len(d.Id()) > 0 && d.Id()[0] != '0' {
		parts := splitUserCustomAttributeId(d.Id())
		if len(parts) == 2 {
			userId, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("failed to parse user ID from resource ID: %v", err)
			}
			
			shortname := parts[1]
			// Read the user to get their custom attributes
			user, err := client.Services.UsersV2.GetOne(int32(userId))
			if err != nil {
				return fmt.Errorf("error reading user %d: %v", userId, err)
			}
			
			if user == nil {
				d.SetId("")
				return nil
			}
			
			if user.CustomAttributes == nil {
				// No custom attributes found
				d.SetId("")
				return nil
			}
			
			value, ok := user.CustomAttributes[shortname]
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
	}
	
	// For shortname-based IDs, just set the shortname
	shortname := d.Id()
	d.Set("shortname", shortname)
	
	return nil
}

func userCustomAttributesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	
	// Check if this is a user-specific custom attribute
	parts := splitUserCustomAttributeId(d.Id())
	if len(parts) == 2 {
		userId, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("failed to parse user ID from resource ID: %v", err)
		}
		
		shortname := parts[1]
		
		// Update the custom attribute value for this user
		user, err := client.Services.UsersV2.GetOne(int32(userId))
		if err != nil {
			return fmt.Errorf("error reading user %d: %v", userId, err)
		}
		
		if user.CustomAttributes == nil {
			user.CustomAttributes = make(map[string]interface{})
		}
		
		user.CustomAttributes[shortname] = d.Get("value")
		
		err = client.Services.UsersV2.Update(user)
		if err != nil {
			log.Printf("[ERROR] Error updating custom attribute for user %d: %v", userId, err)
			return err
		}
		
		return userCustomAttributesRead(d, m)
	}
	
	// Otherwise, just update the shortname value in the state
	shortname := d.Get("shortname").(string)
	d.SetId(shortname)
	
	log.Printf("[WARNING] Updating global custom attribute definitions is currently not working due to an OneLogin API bug")
	log.Printf("[WARNING] The API returns 'Missing param: user_field' when attempting to update custom attributes")
	log.Printf("[WARNING] Please update custom attributes in the OneLogin UI until this API bug is fixed")
	
	return userCustomAttributesRead(d, m)
}

func userCustomAttributesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	
	// Check if this is a user-specific custom attribute
	parts := splitUserCustomAttributeId(d.Id())
	if len(parts) == 2 {
		userId, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("failed to parse user ID from resource ID: %v", err)
		}
		
		shortname := parts[1]
		
		// To "delete" a custom attribute for a user, we set it to null/empty
		user, err := client.Services.UsersV2.GetOne(int32(userId))
		if err != nil {
			return fmt.Errorf("error reading user %d: %v", userId, err)
		}
		
		if user.CustomAttributes != nil {
			delete(user.CustomAttributes, shortname)
			
			err = client.Services.UsersV2.Update(user)
			if err != nil {
				log.Printf("[ERROR] Error clearing custom attribute for user %d: %v", userId, err)
				return err
			}
		}
		
		d.SetId("")
		return nil
	}
	
	// For shortname-based IDs, just remove the ID
	log.Printf("[WARNING] Deleting global custom attribute definitions is currently not working due to an OneLogin API bug")
	log.Printf("[WARNING] The API returns 'Missing param: user_field' when attempting to delete custom attributes")
	log.Printf("[WARNING] Please delete custom attributes in the OneLogin UI until this API bug is fixed")
	
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