package onelogin

import (
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

const (
	BrandingPath string = "api/2/branding/brands"
)

func (sdk *OneloginSDK) ListAccountBrands(query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetBrandByID(brandID int, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateBrand(brand mod.Brand) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, brand)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateBrand(brandID int, brand mod.Brand) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, brand)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteBrand(brandID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetBrandApps(brandID int, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "apps")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// Branding Service: Templates
func (sdk *OneloginSDK) ListBrandTemplates(brandID int, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) CreateBrandTemplate(brandID int, template *mod.BrandTemplate) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates")
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Post(&p, template)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

// func (sdk *OneloginSDK) CreateBrandTemplate1(brandID int, template mod.BrandTemplate) (interface{}, error) {
// 	// Prepare the API path for the request
// 	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Marshal the template into JSON using the custom MarshalJSON function
// 	templateJSON, err := json.Marshal(template)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling template: %v", err)
// 	}

// 	// Send the POST request with the marshaled JSON data
// 	resp, err := sdk.Client.Post(&p, templateJSON)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check and return the response
// 	return utl.CheckHTTPResponse(resp)
// }

func (sdk *OneloginSDK) GetBrandTemplates(brandID int, templateID int, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateBrandTemplate(brandID int, templateID int, template *mod.BrandTemplate) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, template)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) DeleteBrandTemplate(brandID int, templateID int) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateID)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Delete(&p)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetBrandTemplateByType(brandID int, templateType string, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateType)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetBrandTemplateByTypeAndLocale(brandID int, templateType, locale string, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateType, locale)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) UpdateBrandTemplateByTypeAndLocale(brandID int, templateType, locale string, requestBody interface{}) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, brandID, "templates", templateType, locale)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Put(&p, requestBody)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetMasterBrandTemplateByType(templateType string, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, "master", "templates", templateType)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}

func (sdk *OneloginSDK) GetMasterBrandTemplateByTypeAndLocale(templateType, locale string, query mod.Queryable) (interface{}, error) {
	p, err := utl.BuildAPIPath(BrandingPath, "master", "templates", templateType, locale)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.Client.Get(&p, query)
	if err != nil {
		return nil, err
	}
	return utl.CheckHTTPResponse(resp)
}
