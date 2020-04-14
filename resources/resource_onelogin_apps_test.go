package resources

import(
  "testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/onelogin/onelogin-go-sdk/pkg/client"
  "github.com/onelogin/onelogin-go-sdk/pkg/models"
  "github.com/onelogin/onelogin-terraform-provider/resources/app"
  "github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
  "github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
)

func TestAccApp_crud(t *testing.T){
  resource.Test(t, resource,TestCase{
    Steps: []resource.TestStep{{

      }}
  })
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
  testAccProvider = Provider().(*schema.Provider)
  testAccProviders = map[string]terraform.ResourceProvider{
    "example": testAccProvider,
  }
}
