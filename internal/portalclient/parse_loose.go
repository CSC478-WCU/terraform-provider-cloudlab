package portalclient

import (
	"encoding/json"
	"fmt"

	portal "github.com/csc478-wcu/portalctl/portal"
)

// ParseStatusJSONLoose extracts and decodes the first top-level JSON object
// from s, tolerating non-JSON text before/after. We use this for reads only.
func ParseStatusJSONLoose(s string) (*portal.StatusPayload, error) {
	// Try strict first.
	var fast portal.StatusPayload
	if err := json.Unmarshal([]byte(s), &fast); err == nil {
		return &fast, nil
	}

	// Scan for a JSON object and decode the first one that works.
	start, depth := -1, 0
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr {
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth > 0 {
				depth--
				if depth == 0 && start >= 0 {
					var out portal.StatusPayload
					if err := json.Unmarshal([]byte(s[start:i+1]), &out); err == nil {
						return &out, nil
					}
					start = -1 // keep scanning (there might be another object)
				}
			}
		}
	}
	return nil, fmt.Errorf("no decodable JSON object found in status output (len=%d)", len(s))
}
