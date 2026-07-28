package onelogin

import "fmt"

func usersFromResponse(result interface{}) ([]interface{}, error) {
	switch v := result.(type) {
	case map[string]interface{}:
		data, ok := v["data"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid response format from API")
		}
		return data, nil
	case []interface{}:
		return v, nil
	default:
		return nil, fmt.Errorf("Invalid response format from API")
	}
}
