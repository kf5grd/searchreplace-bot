package logr

import "encoding/json"

// toJson marshals an object into a json string
func toJson(b interface{}) string {
	s, _ := json.Marshal(b)
	return string(s)
}

// toJson marshals an object into a json string with indenting
func toJsonPretty(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "  ")
	return string(s)
}
