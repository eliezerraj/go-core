package utils

import (
	"encoding/json"
	"fmt"
)

// FormatHeadersAsJSON converts HTTP headers map to a formatted JSON string
func FormatHeadersAsJSON(headers map[string][]string) []byte {
	jsonBytes, err := json.Marshal(headers)
	if err != nil {
		return []byte(fmt.Sprintf("%v", headers))
	}
	return jsonBytes
}
