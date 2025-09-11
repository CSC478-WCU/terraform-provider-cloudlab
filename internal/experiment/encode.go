package experiment

import (
	"encoding/json"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/model"
)

func encodeSpec(spec model.ExperimentSpec) (string, error) {
	b, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
