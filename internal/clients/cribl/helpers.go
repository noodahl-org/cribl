package cribl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleResult(resp *http.Response, err error, out interface{}) error {
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		if data, err := io.ReadAll(resp.Body); err != nil {
			return err
		} else if err := json.Unmarshal(data, out); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("status code: %v", resp.StatusCode)
	}
	return nil
}
