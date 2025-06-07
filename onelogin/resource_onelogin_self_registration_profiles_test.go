package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSelfRegistrationProfile_crud(t *testing.T) {
	base := GetFixture("onelogin_self_registration_profile_example.tf", t)
	update := GetFixture("onelogin_self_registration_profile_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "name", "Test Profile"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "url", "test-profile"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "enabled", "true"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "moderated", "false"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "helptext", "Welcome to our community. Please follow the guidelines."),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "thankyou_message", "Thank you for joining!"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "email_verification_type", "Email MagicLink"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "name", "Updated Test Profile"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "url", "test-profile"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "moderated", "true"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "helptext", "Updated welcome message."),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "thankyou_message", "Updated thank you message!"),
					resource.TestCheckResourceAttr("onelogin_self_registration_profiles.basic_test", "email_verification_type", "Email OTP"),
				),
			},
		},
	})
}
