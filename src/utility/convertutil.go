package utility

import (
	"bytes"
	"encoding/json"
)

func ConvertToJSON(indata interface{}) (string, error) {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(indata); err != nil {
		return "", nil
	}

	return b.String(), nil
}
