package helper

import "encoding/json"

func PrettyJson(value interface{}) ([]byte, error) {
	return json.MarshalIndent(value, "", "    ")

}
