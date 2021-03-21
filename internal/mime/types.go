package mime

import (
	_ "embed"
	"encoding/json"
)

//go:embed types.json
var mimeTypes []byte

func Types() ([]string, error) {
	var s []string

	err := json.Unmarshal(mimeTypes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
