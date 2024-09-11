package handlers

import "encoding/json"

func newInternalError(cause error, message string) ([]byte, error) {
	internalError := map[string]string{
		"cause":   cause.Error(),
		"message": message,
	}

	b, err := json.Marshal(internalError)
	return b, err
}
