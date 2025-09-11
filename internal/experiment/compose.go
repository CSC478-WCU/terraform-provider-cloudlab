package experiment

import "encoding/json"

func composeParams(project, name, specJSON string) map[string]any {
	bindingsMap := map[string]string{profileParamSpecJSON: specJSON}
	b, _ := json.Marshal(bindingsMap)

	return map[string]any{
		"proj":     project,
		"profile":  profileName,
		"name":     name,
		"bindings": string(b), // JSON object as STRING (required by portal API)
	}
}
