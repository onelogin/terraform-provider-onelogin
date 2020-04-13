package configuration

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppConfiguration returns a key/value map of the various fields that make up
// the AppConfiguration field for a OneLogin App.
func SAMLConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"provider_arn": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"signature_algorithm": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validSignatureAlgo,
		},
	}
}

func InflateSAMLConfiguration(s *map[string]interface{}) *models.AppConfiguration {
	out := models.AppConfiguration{}
	var st string
	var notNil bool
	if st, notNil = (*s)["provider_arn"].(string); notNil {
		out.ProviderArn = oltypes.String(st)
	}
	if st, notNil = (*s)["signature_algorithm"].(string); notNil {
		out.SignatureAlgorithm = oltypes.String(st)
	}
	return &out
}

func validSignatureAlgo(val interface{}, key string) (warns []string, errs []error) {
	validOpts := []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"}
	v := val.(string)
	isValid := false
	for _, o := range validOpts {
		isValid = v == o
		if isValid {
			break
		}
	}
	if !isValid {
		errs = append(errs, fmt.Errorf("signature_algorithm must be one of %v, got: %s", validOpts, v))
	}
	return
}
