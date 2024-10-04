package write

import (
	"encoding/json"
	"github.com/ankodd/todo-server/pkg/models/http/response"
)

func Write(response *response.Response) error {
	response.Writer.WriteHeader(response.Status)

	err := json.NewEncoder(response.Writer).Encode(response)
	if err != nil {
		return err
	}

	return nil
}
